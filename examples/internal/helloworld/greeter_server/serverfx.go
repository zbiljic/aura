package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_fx "github.com/zbiljic/aura/go/pkg/grpc/fx"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/zbiljic/aura/examples/internal/proto/helloworld"
)

var serverfx = fx.Provide(newServer)

type serverParams struct {
	fx.In

	Log *zap.SugaredLogger
}

func newServer(p serverParams) grpc_fx.GRPCServerInputResult {
	server := &server{log: p.Log}

	return grpc_fx.GRPCServerInputResult{
		Services: []grpc_fx.RegisterFn{
			func() (string, grpc_fx.ServerRegisterFn, grpc_fx.HandlerRegisterFn) {
				return "helloworld",
					func(s *grpc.Server) {
						pb.RegisterGreeterServer(s, server)
					},
					func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
						return pb.RegisterGreeterHandler(ctx, mux, conn)
					}
			},
		},
	}
}
