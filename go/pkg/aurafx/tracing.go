package aurafx

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"

	otelaura "github.com/zbiljic/aura/go/pkg/otel"
)

var tracingfx = fx.Options(
	fx.Provide(ProvideTracer),
	fx.Provide(ProvideTracerProvider),
	fx.Invoke(NewTracerCloser),
)

func ProvideTracer(
	log *zap.SugaredLogger,
	tracingConfig otelaura.Config,
) (*otelaura.Tracer, error) {
	tracer, err := otelaura.New(log, tracingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %v", err)
	}

	return tracer, nil
}

func ProvideTracerProvider(tracer *otelaura.Tracer) trace.TracerProvider {
	return tracer.TracerProvider()
}

type TracingParams struct {
	fx.In

	Lifecycle fx.Lifecycle

	Log    *zap.SugaredLogger
	Tracer *otelaura.Tracer
}

func NewTracerCloser(p TracingParams) error {
	p.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			p.Log.Infof("stopping tracing")

			p.Tracer.Shutdown(ctx)

			return nil
		},
	})

	return nil
}
