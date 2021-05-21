package tracing

import (
	"fmt"
	"io"
	"strings"

	"github.com/opentracing/opentracing-go"
	zipkinOT "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinHttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jaegerConf "github.com/uber/jaeger-client-go/config"
	jaegerZap "github.com/uber/jaeger-client-go/log/zap"
	jaegerZipkin "github.com/uber/jaeger-client-go/zipkin"
	jaegerProm "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
)

// Tracer encapsulates tracing abilities.
type Tracer struct {
	Config *Config

	log    *zap.SugaredLogger
	tracer opentracing.Tracer
	closer io.Closer
}

func New(log *zap.SugaredLogger, c *Config) (*Tracer, error) {
	t := &Tracer{Config: c, log: log}

	if err := t.setup(); err != nil {
		return nil, err
	}

	return t, nil
}

// setup sets up the tracer.
func (t *Tracer) setup() error {
	switch strings.ToLower(t.Config.Provider) {
	case "jaeger":
		jc, err := jaegerConf.FromEnv()

		if err != nil {
			return err
		}

		if t.Config.Jaeger.LocalAgentHostPort != "" {
			jc.Reporter.LocalAgentHostPort = t.Config.Jaeger.LocalAgentHostPort
		} else {
			jc.Reporter.LocalAgentHostPort = fmt.Sprintf("127.0.0.1:%d", jaeger.DefaultUDPSpanServerPort)
		}

		if t.Config.Jaeger.SamplingType != "" {
			jc.Sampler.Type = t.Config.Jaeger.SamplingType
		} else {
			jc.Sampler.Type = jaeger.SamplerTypeConst
		}

		if t.Config.Jaeger.SamplingValue != 0 {
			jc.Sampler.Param = t.Config.Jaeger.SamplingValue
		} else {
			jc.Sampler.Param = 1
		}

		if t.Config.Jaeger.SamplingServerURL != "" {
			jc.Sampler.SamplingServerURL = t.Config.Jaeger.SamplingServerURL
		} else {
			jc.Sampler.SamplingServerURL = jaeger.DefaultSamplingServerURL
		}

		var configs []jaegerConf.Option

		if t.Config.Jaeger.MaxTagValueLength > 0 &&
			t.Config.Jaeger.MaxTagValueLength != jaeger.DefaultMaxTagValueLength {
			configs = append(configs, jaegerConf.MaxTagValueLength(t.Config.Jaeger.MaxTagValueLength))
		}

		// This works in other jaeger clients, but is not part of jaeger-client-go
		if t.Config.Jaeger.Propagation == "b3" {
			zipkinPropagator := jaegerZipkin.NewZipkinB3HTTPHeaderPropagator()
			configs = append(
				configs,
				jaegerConf.Injector(opentracing.HTTPHeaders, zipkinPropagator),
				jaegerConf.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
			)
		}

		logAdapt := jaegerZap.NewLogger(t.log.Desugar())
		configs = append(configs, jaegerConf.Logger(logAdapt))

		factory := jaegerProm.New() // By default uses prometheus.DefaultRegisterer
		configs = append(configs, jaegerConf.Metrics(factory))

		closer, err := jc.InitGlobalTracer(
			t.Config.ServiceName,
			configs...,
		)

		if err != nil {
			return err
		}

		t.closer = closer
		t.tracer = opentracing.GlobalTracer()
		t.log.Infof("Jaeger tracer configured")

	case "zipkin":
		if t.Config.Zipkin.ServerURL == "" {
			return errors.Errorf("Zipkin's server url is required")
		}

		reporter := zipkinHttp.NewReporter(t.Config.Zipkin.ServerURL)

		endpoint, err := zipkin.NewEndpoint(t.Config.ServiceName, "")

		if err != nil {
			return err
		}

		nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))

		if err != nil {
			return err
		}

		opentracing.SetGlobalTracer(zipkinOT.Wrap(nativeTracer))

		t.closer = reporter
		t.tracer = opentracing.GlobalTracer()
		t.log.Infof("Zipkin tracer configured")

	case "":
		t.log.Infof("No tracer configured - skipping tracing setup")
	default:
		return errors.Errorf("unknown tracer: %s", t.Config.Provider)
	}

	return nil
}

// IsLoaded returns true if the tracer has been loaded.
func (t *Tracer) IsLoaded() bool {
	if t == nil || t.tracer == nil {
		return false
	}
	return true
}

// Tracer returns the wrapped tracer.
func (t *Tracer) Tracer() opentracing.Tracer {
	return t.tracer
}

// Close closes the tracer.
func (t *Tracer) Close() {
	if t.closer != nil {
		err := t.closer.Close()
		if err != nil {
			t.log.Errorf("Unable to close tracer: %w", err)
		}
	}
}
