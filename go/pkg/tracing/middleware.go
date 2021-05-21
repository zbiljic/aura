package tracing

import (
	"fmt"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
)

func NewTracedHttpHandler(tracer opentracing.Tracer, h http.Handler) http.Handler {
	return nethttp.Middleware(
		tracer,
		h,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return fmt.Sprintf("HTTP-gRPC %s %s", r.Method, r.URL.String())
		}),
	)
}
