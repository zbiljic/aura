package main

import (
	"go.uber.org/fx"

	"github.com/zbiljic/aura/go/pkg/aurafx"
	grpc_fx "github.com/zbiljic/aura/go/pkg/grpc/fx"
)

var configfx = fx.Provide(
	provideRootConfig,
	provideGRPCConfig,
	provideGatewayConfig,
)

func provideRootConfig(config *config) *aurafx.Config {
	return &config.Config
}

func provideGRPCConfig(config *config) *grpc_fx.GRPCConfig {
	return config.GRPC
}

func provideGatewayConfig(config *config) *grpc_fx.GatewayConfig {
	return config.Gateway
}
