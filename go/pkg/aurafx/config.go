package aurafx

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"

	"github.com/zbiljic/aura/go/pkg/logger"
	otelaura "github.com/zbiljic/aura/go/pkg/otel"
)

var configfx = fx.Provide(
	ProvideLoggerConfig,
	ProvideTracingConfig,
	ProvideDebugConfig,
	ProvideAdminConfig,
)

type Config struct {
	AppName string           `json:"app_name" validate:"required"`
	Logger  *logger.Config   `json:"logger" validate:"dive"`
	Tracing *otelaura.Config `json:"tracing" validate:"dive"`
	Debug   *DebugConfig     `json:"debug" validate:"dive"`
	Admin   *AdminConfig     `json:"admin" validate:"dive"`
}

type namer interface {
	withName(string)
}

func (config *Config) withName(name string) {
	config.AppName = name
}

func ProvideLoggerConfig(config *Config) *logger.Config {
	return config.Logger
}

func ProvideTracingConfig(config *Config) *otelaura.Config {
	tracingConfig := config.Tracing
	if tracingConfig.ServiceName == "" {
		tracingConfig.ServiceName = config.AppName
	}
	return tracingConfig
}

func ProvideDebugConfig(config *Config) *DebugConfig {
	return config.Debug
}

func ProvideAdminConfig(config *Config) *AdminConfig {
	return config.Admin
}

func loadEnvironment(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Load(filename)
	} else {
		err = godotenv.Load()
		// handle if .env file does not exists, this is OK
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

// LoadConfig loads configuration.
func LoadConfig(filename, appName string, spec interface{}) error {
	if err := loadEnvironment(filename); err != nil {
		return err
	}

	if config, ok := spec.(namer); ok {
		config.withName(appName)
	}

	if err := envconfig.Process(appName, spec); err != nil {
		return err
	}

	if err := ValidateConfig(spec); err != nil {
		return err
	}

	return nil
}

func ValidateConfig(spec interface{}) error {
	validate := validator.New()

	if err := validate.Struct(spec); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors[0]
	}

	return nil
}
