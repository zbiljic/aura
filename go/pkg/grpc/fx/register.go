package grpc_fx

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type RegisterFn func() (string, ServerRegisterFn, HandlerRegisterFn)

type ServerRegisterFn func(*grpc.Server)

type HandlerRegisterFn func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
