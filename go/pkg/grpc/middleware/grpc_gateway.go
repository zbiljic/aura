package grpc_middleware

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/metadata"
)

func TracingMetadataAnnotator(ctx context.Context, _ *http.Request) metadata.MD {
	carrier := make(metadata.MD)
	otelgrpc.Inject(ctx, &carrier)
	return carrier
}
