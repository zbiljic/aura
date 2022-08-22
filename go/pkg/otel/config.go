package otelaura

import "time"

type Config struct {
	// Specifies the service name to use on the tracer.
	ServiceName string `json:"service_name" split_words:"true"`

	// Set this to the tracing backend you wish to use.
	// If omitted or empty, tracing will be disabled.
	Provider string `json:"provider" default:"jaeger"`

	// Use synchronous span exporter processor.
	//
	// This is not recommended for production use.
	Sync bool `json:"sync" default:"false"`

	OTLP   *OTLPConfig   `json:"otlp" validate:"dive"`
	Jaeger *JaegerConfig `json:"jaeger" validate:"dive"`
	Zipkin *ZipkinConfig `json:"zipkin" validate:"dive"`
}

// OTLPConfig encapsulates OTLP exporter configuration.
type OTLPConfig struct {
	Endpoint    string            `json:"endpoint"`
	Insecure    bool              `json:"insecure" default:"true"`
	Headers     map[string]string `json:"headers"`
	Compression string            `json:"compression"`
	Timeout     time.Duration     `json:"timeout"`
	Protocol    string            `json:"protocol" default:"grpc"`
}

// JaegerConfig encapsulates jaeger's configuration.
type JaegerConfig struct {
	AgentHost string `json:"agent_host" split_words:"true"`
	AgentPort string `json:"agent_port" split_words:"true"`
}

// ZipkinConfig encapsulates zipkin's configuration.
type ZipkinConfig struct {
	ServerURL string `json:"server_url" split_words:"true"`
}
