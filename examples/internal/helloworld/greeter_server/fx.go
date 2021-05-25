package main

import (
	"go.uber.org/fx"

	"github.com/zbiljic/aura/go/pkg/aurafx"
	grpc_fx "github.com/zbiljic/aura/go/pkg/grpc/fx"
)

var fxModule = fx.Options(
	configfx,
	aurafx.Module,
	fx.Invoke(grpc_fx.NewGRPCServer),
	fx.Invoke(grpc_fx.NewGateway),
	serverfx,
)
