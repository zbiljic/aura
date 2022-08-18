package otelaura

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Tracer encapsulates tracing abilities.
type Tracer struct {
	Config *Config

	log            *zap.SugaredLogger
	tracerProvider trace.TracerProvider
	shutdownFn     func(context.Context) error
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
	// common resources
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.TelemetrySDKLanguageGo,
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(t.Config.ServiceName),
		),
	)
	if err != nil {
		return fmt.Errorf("could not set resources: %w", err)
	}

	spanProcessorOptionFn := func(sync bool, exporter sdktrace.SpanExporter) sdktrace.TracerProviderOption {
		if sync {
			return sdktrace.WithSyncer(exporter)
		}
		return sdktrace.WithBatcher(exporter)
	}

	switch strings.ToLower(t.Config.Provider) {
	case "jaeger":
		var (
			options []jaeger.AgentEndpointOption
		)

		if t.Config.Jaeger.AgentHost != "" {
			options = append(options, jaeger.WithAgentHost(t.Config.Jaeger.AgentHost))
		}

		if t.Config.Jaeger.AgentPort != "" {
			options = append(options, jaeger.WithAgentPort(t.Config.Jaeger.AgentPort))
		}

		options = append(options, jaeger.WithLogger(zap.NewStdLog(t.log.Desugar())))

		// Create the Jaeger exporter
		exporter, err := jaeger.New(
			jaeger.WithAgentEndpoint(options...),
		)
		if err != nil {
			return err
		}

		t.tracerProvider = sdktrace.NewTracerProvider(
			spanProcessorOptionFn(t.Config.Sync, exporter),
			sdktrace.WithResource(resources),
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		)

		t.shutdownFn = exporter.Shutdown

		otel.SetTracerProvider(t.tracerProvider)

		t.log.Infof("Jaeger tracer configured")

	case "zipkin":
		// Create the Zipkin exporter
		exporter, err := zipkin.New(
			t.Config.Zipkin.ServerURL,
			zipkin.WithLogger(zap.NewStdLog(t.log.Desugar())),
		)
		if err != nil {
			return err
		}

		t.tracerProvider = sdktrace.NewTracerProvider(
			spanProcessorOptionFn(t.Config.Sync, exporter),
			sdktrace.WithResource(resources),
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		)

		t.shutdownFn = exporter.Shutdown

		otel.SetTracerProvider(t.tracerProvider)

		t.log.Infof("Zipkin tracer configured")

	case "":
		t.log.Infof("no tracer configured - skipping tracing setup")
	default:
		return fmt.Errorf("unknown tracer: %s", t.Config.Provider)
	}

	// setup propagators
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)

	return nil
}

// IsLoaded returns true if the tracer has been loaded.
func (t *Tracer) IsLoaded() bool {
	if t == nil || t.tracerProvider == nil {
		return false
	}
	return true
}

// TracerProvider returns the wrapped tracer.
func (t *Tracer) TracerProvider() trace.TracerProvider {
	return t.tracerProvider
}

// Shutdown stops the Tracer.
func (t *Tracer) Shutdown(ctx context.Context) {
	if t.shutdownFn != nil {
		err := t.shutdownFn(ctx)
		if err != nil {
			t.log.Errorf("unable to shutdown exporter: %w", err)
		}
	}
}
