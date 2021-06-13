package corsx

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/zbiljic/aura/go/pkg/httpx"
)

func CORSHandler(config *Config) httpx.HandlerResult {
	result := httpx.HandlerResult{}

	if options, enabled := config.CORSOptions(); enabled {
		result.HTTPMiddleware = []func(h http.Handler) http.Handler{
			cors.New(options).Handler,
		}
	}

	return result
}
