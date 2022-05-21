package otelaura_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	otelaura "github.com/zbiljic/aura/go/pkg/otel"
)

func TestTracedHttpHandler(t *testing.T) {
	expectedTagsSuccess := map[attribute.Key]attribute.Value{
		semconv.HTTPServerNameKey: attribute.StringValue("test"),
		semconv.HTTPTargetKey:     attribute.StringValue("/"),
		semconv.HTTPMethodKey:     attribute.StringValue("GET"),
		semconv.HTTPStatusCodeKey: attribute.IntValue(200),
	}

	expectedTagsError := map[attribute.Key]attribute.Value{
		semconv.HTTPServerNameKey: attribute.StringValue("test"),
		semconv.HTTPTargetKey:     attribute.StringValue("/"),
		semconv.HTTPMethodKey:     attribute.StringValue("GET"),
		semconv.HTTPStatusCodeKey: attribute.IntValue(400),
	}

	testCases := []struct {
		httpStatus      int
		testDescription string
		expectedTags    map[attribute.Key]attribute.Value
	}{
		{
			testDescription: "success http response",
			httpStatus:      http.StatusOK,
			expectedTags:    expectedTagsSuccess,
		},
		{
			testDescription: "error http response",
			httpStatus:      http.StatusBadRequest,
			expectedTags:    expectedTagsError,
		},
	}

	for _, test := range testCases {
		t.Run(test.testDescription, func(t *testing.T) {
			spanRecorder := tracetest.NewSpanRecorder()
			provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.httpStatus)
			})

			ts := httptest.NewServer(otelaura.NewTracedHttpHandler(mux, "test", provider))
			defer ts.Close()

			_, err := http.Get(ts.URL)
			require.NoError(t, err)

			spans := spanRecorder.Ended()
			assert.Len(t, spans, 1)

			if !spans[0].SpanContext().IsValid() {
				t.Fatalf("invalid span created: %#v", spans[0].SpanContext())
			}

			spanAttributes := map[attribute.Key]attribute.Value{}
			for _, kv := range spans[0].Attributes() {
				// NOTE: only check defined keys
				if _, ok := test.expectedTags[kv.Key]; ok {
					spanAttributes[kv.Key] = kv.Value
				}
			}

			assert.Equal(t, test.expectedTags, spanAttributes)
		})
	}
}
