package grpc_middleware

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
)

func TracingMetadataAnnotator(ctx context.Context, _ *http.Request) metadata.MD {
	span := opentracing.SpanFromContext(ctx)
	return MetadataFromSpan(span)
}

func MetadataFromSpan(span opentracing.Span) metadata.MD {
	ctx := span.Context()
	carrier := make(map[string]string)

	//nolint:errcheck
	span.Tracer().Inject(
		ctx,
		opentracing.TextMap,
		opentracing.TextMapCarrier(carrier),
	)

	return metadata.New(carrier)
}
