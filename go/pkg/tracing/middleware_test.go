package tracing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zbiljic/aura/go/pkg/tracing"
)

var mockedTracer *mocktracer.MockTracer

func init() {
	mockedTracer = mocktracer.New()
	opentracing.SetGlobalTracer(mockedTracer)
}

func TestTracedHttpHandler(t *testing.T) {
	expectedTagsSuccess := map[string]interface{}{
		string(ext.SpanKind):       ext.SpanKindEnum("server"),
		string(ext.Component):      "net/http",
		string(ext.HTTPUrl):        "/",
		string(ext.HTTPMethod):     "GET",
		string(ext.HTTPStatusCode): uint16(200),
	}

	expectedTagsError := map[string]interface{}{
		string(ext.SpanKind):       ext.SpanKindEnum("server"),
		string(ext.Component):      "net/http",
		string(ext.HTTPUrl):        "/",
		string(ext.HTTPMethod):     "GET",
		string(ext.HTTPStatusCode): uint16(400),
	}

	testCases := []struct {
		httpStatus      int
		testDescription string
		expectedTags    map[string]interface{}
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
			defer mockedTracer.Reset()

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.httpStatus)
			})

			ts := httptest.NewServer(tracing.NewTracedHttpHandler(mockedTracer, mux))
			defer ts.Close()

			_, err := http.Get(ts.URL)
			require.NoError(t, err)

			spans := mockedTracer.FinishedSpans()
			assert.Len(t, spans, 1)
			span := spans[0]

			assert.Equal(t, test.expectedTags, span.Tags())
		})
	}
}
