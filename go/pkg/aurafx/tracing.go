package aurafx

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/zbiljic/aura/go/pkg/tracing"
)

var tracingfx = fx.Options(
	fx.Provide(ProvideTracer),
	fx.Invoke(NewTracerCloser),
)

func ProvideTracer(
	log *zap.SugaredLogger,
	tracingConfig *tracing.Config,
) (opentracing.Tracer, error) {
	tracer, err := tracing.New(log, tracingConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %v", err)
	}

	return tracer.Tracer(), nil
}

type TracingParams struct {
	fx.In

	Lifecycle fx.Lifecycle

	Log    *zap.SugaredLogger
	Tracer opentracing.Tracer
}

func NewTracerCloser(p TracingParams) error {
	p.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			p.Log.Infof("closing tracing")

			if tracer, ok := p.Tracer.(io.Closer); ok {
				tracer.Close()
			}

			return nil
		},
	})

	return nil
}
