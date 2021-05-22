package aurafx

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

var metricsfx = fx.Provide(
	ProvidePrometheusRegistry,
	NewMetricsAdmin,
)

func ProvidePrometheusRegistry() (prometheus.Registerer, prometheus.Gatherer) {
	return prometheus.DefaultRegisterer, prometheus.DefaultGatherer
}

func NewMetricsAdmin(gatherer prometheus.Gatherer) AdminHandlerResult {
	return AdminHandlerResult{
		AdminHandlers: map[string]http.Handler{
			// Expose prometheus metrics
			"/metrics": promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}),
		},
	}
}
