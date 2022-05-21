package otelaura

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

	Jaeger *JaegerConfig `json:"jaeger" validate:"dive"`
	Zipkin *ZipkinConfig `json:"zipkin" validate:"dive"`
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
