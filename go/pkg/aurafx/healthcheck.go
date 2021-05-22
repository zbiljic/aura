package aurafx

import (
	"net/http"

	"github.com/heptiolabs/healthcheck"
	"go.uber.org/fx"
)

var healthcheckfx = fx.Provide(
	ProvideHealthCheckHandler,
	NewHealthCheckAdmin,
)

func ProvideHealthCheckHandler() healthcheck.Handler {
	return healthcheck.NewHandler()
}

func NewHealthCheckAdmin(health healthcheck.Handler) AdminHandlerResult {
	return AdminHandlerResult{
		AdminHandlers: map[string]http.Handler{
			// Expose a liveness check
			"/health/live": http.HandlerFunc(health.LiveEndpoint),
			// Expose a readiness check
			"/health/ready": http.HandlerFunc(health.ReadyEndpoint),
		},
	}
}
