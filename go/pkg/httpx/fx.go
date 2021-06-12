package httpx

import (
	"net/http"

	"go.uber.org/fx"
)

type HandlerResult struct {
	fx.Out

	HTTPMiddleware []func(http.Handler) http.Handler `group:"http_middleware,flatten"`
}
