package grpc_fx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	grpc_middleware "github.com/zbiljic/aura/go/pkg/grpc/middleware"
	"github.com/zbiljic/aura/go/pkg/tracing"
)

type GatewayConfig struct {
	Enabled      bool          `json:"enabled" default:"true"`
	Addr         string        `json:"addr" default:":8080"`
	ReadTimeout  time.Duration `json:"read_timeout" default:"30s" split_words:"true"`
	WriteTimeout time.Duration `json:"write_timeout" default:"30s" split_words:"true"`
	IdleTimeout  time.Duration `json:"idle_timeout" default:"120s" split_words:"true"`
}

type GatewayParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	OnErrorCh chan error

	Log           *zap.SugaredLogger
	Tracer        opentracing.Tracer
	GRPCConfig    *GRPCConfig
	GatewayConfig *GatewayConfig

	Services []RegisterFn `group:"service"`
}

func NewGateway(p GatewayParams) error {
	addr, err := net.ResolveTCPAddr("tcp", p.GatewayConfig.Addr)
	if err != nil {
		return fmt.Errorf("could not resolve TCP address: %w", err)
	}

	hostPort := addr.String()

	ctx := context.Background()

	conn, err := grpc.DialContext(
		ctx,
		p.GRPCConfig.Addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("could not dial gRPC server: %w", err)
	}

	gatewayMux := runtime.NewServeMux(
		runtime.WithMetadata(grpc_middleware.TracingMetadataAnnotator),
	)

	// register handlers
	if len(p.Services) > 0 {
		p.Log.Debugf("registering %d handlers", len(p.Services))

		for _, s := range p.Services {
			name, _, handlerFn := s()

			p.Log.Debugf("register handler: %s", name)

			if err := handlerFn(ctx, gatewayMux, conn); err != nil {
				return fmt.Errorf("failed to register handler for '%s': %w", name, err)
			}
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", gatewayMux)

	server := &http.Server{
		Addr:         hostPort,
		Handler:      tracing.NewTracedHttpHandler(p.Tracer, mux),
		ReadTimeout:  p.GatewayConfig.ReadTimeout,
		WriteTimeout: p.GatewayConfig.WriteTimeout,
		IdleTimeout:  p.GatewayConfig.IdleTimeout,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				p.Log.Infof("starting gRPC gateway at: %s", hostPort)

				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					p.Log.Errorf("unable to start server: %w", err)
					p.OnErrorCh <- err
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Log.Info("stopping gRPC gateway")

			if err := conn.Close(); err != nil {
				return err
			}

			return server.Shutdown(ctx)
		},
	})

	return nil
}