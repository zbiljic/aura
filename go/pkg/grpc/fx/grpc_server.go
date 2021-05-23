package grpc_fx

import (
	"context"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Enabled bool   `json:"enabled" default:"true"`
	Addr    string `json:"addr" default:":7070"`
}

type GRPCServerInputResult struct {
	fx.Out

	UnaryInterceptors  []grpc.UnaryServerInterceptor  `group:"grpc_unary_server_interceptor,flatten"`
	StreamInterceptors []grpc.StreamServerInterceptor `group:"grpc_stream_server_interceptor,flatten"`
	Services           []RegisterFn                   `group:"service,flatten"`
}

type GRPCServerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	OnErrorCh chan error

	Log        *zap.SugaredLogger
	GRPCConfig *GRPCConfig

	UnaryInterceptors  []grpc.UnaryServerInterceptor  `group:"grpc_unary_server_interceptor"`
	StreamInterceptors []grpc.StreamServerInterceptor `group:"grpc_stream_server_interceptor"`
	Services           []RegisterFn                   `group:"service"`
}

func NewGRPCServer(p GRPCServerParams) error {
	addr, err := net.ResolveTCPAddr("tcp", p.GRPCConfig.Addr)
	if err != nil {
		return fmt.Errorf("could not resolve TCP address: %w", err)
	}

	hostPort := addr.String()

	serverLogger := p.Log.Desugar()

	// Make sure that log statements internal to gRPC library are logged.
	grpc_zap.ReplaceGrpcLogger(serverLogger.WithOptions(zap.IncreaseLevel(zap.ErrorLevel)))

	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(serverLogger),
		grpc_validator.UnaryServerInterceptor(),
	}

	// register additional unary server interceptors
	if len(p.UnaryInterceptors) > 0 {
		p.Log.Debugf("registering %d unary server interceptors", len(p.UnaryInterceptors))
		unaryServerInterceptors = append(unaryServerInterceptors, p.UnaryInterceptors...)
	}

	streamServerInterceptors := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(),
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(serverLogger),
		grpc_validator.StreamServerInterceptor(),
	}

	// register additional stream server interceptors
	if len(p.StreamInterceptors) > 0 {
		p.Log.Debugf("registering %d stream server interceptors", len(p.StreamInterceptors))
		streamServerInterceptors = append(streamServerInterceptors, p.StreamInterceptors...)
	}

	serverOptions := []grpc.ServerOption{}
	serverOptions = append(serverOptions, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(unaryServerInterceptors...)),
	)
	serverOptions = append(serverOptions, grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(streamServerInterceptors...)),
	)

	gRPCServer := grpc.NewServer(serverOptions...)

	// register services
	if len(p.Services) > 0 {
		p.Log.Debugf("registering %d services", len(p.Services))

		for _, s := range p.Services {
			name, srvFn, _ := s()

			p.Log.Debugf("register service: %s", name)

			srvFn(gRPCServer)
		}
	}

	// register metrics
	grpc_prometheus.DefaultServerMetrics.InitializeMetrics(gRPCServer)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				p.Log.Infof("starting gRPC server at: %s", hostPort)

				netListener, err := net.Listen("tcp", hostPort)
				if err != nil {
					p.Log.Errorf("failed to listen: %w", err)
					p.OnErrorCh <- err
					return
				}

				if err := gRPCServer.Serve(netListener); err != nil {
					p.Log.Errorf("failed to serve: %w", err)
					p.OnErrorCh <- err
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Log.Info("stopping gRPC server")

			gRPCServer.GracefulStop()

			p.Log.Debug("stopped gRPC server")
			return nil
		},
	})

	return nil
}
