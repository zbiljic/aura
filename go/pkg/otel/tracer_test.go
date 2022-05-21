package otelaura_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	otelaura "github.com/zbiljic/aura/go/pkg/otel"
)

type zipkinSpanRequest struct {
	Id            string
	TraceId       string
	Timestamp     uint64
	Name          string
	LocalEndpoint struct {
		ServiceName string
	}
	Tags map[string]string
}

func TestZipkinTracer(t *testing.T) {
	done := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(done)

		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		var spans []zipkinSpanRequest
		err = json.Unmarshal(body, &spans)

		assert.NoError(t, err)

		assert.NotEmpty(t, spans[0].Id)
		assert.NotEmpty(t, spans[0].TraceId)
		// NOTE: span name and service name are lowercased as per swagger definition
		assert.Equal(t, "testoperation", spans[0].Name)
		assert.Equal(t, "test", spans[0].LocalEndpoint.ServiceName)
		assert.NotNil(t, spans[0].Tags["testTag"])
		assert.Equal(t, "true", spans[0].Tags["testTag"])
	}))
	defer ts.Close()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	_, err = otelaura.New(log.Sugar(), &otelaura.Config{
		ServiceName: "Test",
		Provider:    "zipkin",
		Sync:        true,
		Zipkin: &otelaura.ZipkinConfig{
			ServerURL: ts.URL,
		},
	})
	assert.NoError(t, err)

	tr := otel.GetTracerProvider().Tracer("test")
	_, span := tr.Start(context.Background(), "testOperation", trace.WithSpanKind(trace.SpanKindServer))
	span.SetAttributes(
		attribute.Bool("testTag", true),
	)
	span.End()

	select {
	case <-done:
	case <-time.After(time.Millisecond * 1500):
		t.Fatalf("Test server did not receive spans")
	}
}
