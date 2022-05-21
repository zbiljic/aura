package otelaura

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func NewTracedHttpHandler(h http.Handler, operation string, provider trace.TracerProvider) http.Handler {
	return otelhttp.NewHandler(
		h,
		operation,
		otelhttp.WithTracerProvider(provider),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("%s %s %s", operation, r.Method, r.URL.String())
		}),
		otelhttp.WithPropagators(propagation.TraceContext{}),
	)
}
