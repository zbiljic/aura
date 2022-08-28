package main

import (
	"github.com/zbiljic/aura/go/pkg/aurafx"
	grpc_fx "github.com/zbiljic/aura/go/pkg/grpc/fx"
)

type config struct {
	aurafx.Config
	GRPC    grpc_fx.GRPCConfig    `json:"grpc" validate:"dive"`
	Gateway grpc_fx.GatewayConfig `json:"gateway" validate:"dive"`
}
