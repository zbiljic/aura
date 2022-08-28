package aurafx

import (
	"go.uber.org/fx"

	"github.com/zbiljic/aura/go/pkg/logger"
	"github.com/zbiljic/aura/go/pkg/tracing"
)

var configfx = fx.Provide(
	ProvideLoggerConfig,
	ProvideTracingConfig,
	ProvideDebugConfig,
	ProvideAdminConfig,
)

type Config struct {
	AppName string         `json:"app_name" validate:"required" alias:"appName"`
	Logger  logger.Config  `json:"logger" validate:"dive"`
	Tracing tracing.Config `json:"tracing" validate:"dive"`
	Debug   DebugConfig    `json:"debug" validate:"dive"`
	Admin   AdminConfig    `json:"admin" validate:"dive"`
}

func ProvideLoggerConfig(config Config) logger.Config {
	return config.Logger
}

func ProvideTracingConfig(config Config) tracing.Config {
	tracingConfig := config.Tracing
	if tracingConfig.ServiceName == "" {
		tracingConfig.ServiceName = config.AppName
	}
	return tracingConfig
}

func ProvideDebugConfig(config Config) DebugConfig {
	return config.Debug
}

func ProvideAdminConfig(config Config) AdminConfig {
	return config.Admin
}
