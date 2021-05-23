package tracing

type Config struct {
	// Specifies the service name to use on the tracer.
	ServiceName string `json:"service_name" split_words:"true"`

	// Set this to the tracing backend you wish to use.
	// If omitted or empty, tracing will be disabled.
	Provider string `json:"provider" default:"jaeger"`

	Jaeger *JaegerConfig `json:"jaeger" validate:"dive"`
	Zipkin *ZipkinConfig `json:"zipkin" validate:"dive"`
}

// JaegerConfig encapsulates jaeger's configuration.
type JaegerConfig struct {
	LocalAgentHostPort string  `json:"local_agent_address" split_words:"true"`
	SamplingType       string  `json:"sampling_type" split_words:"true"`
	SamplingValue      float64 `json:"sampling_value" split_words:"true"`
	SamplingServerURL  string  `json:"sampling_server_url" split_words:"true"`
	Propagation        string  `json:"propagation" default:"jaeger"`
	MaxTagValueLength  int     `json:"max_tag_value_length" split_words:"true"`
}

// ZipkinConfig encapsulates zipkin's configuration.
type ZipkinConfig struct {
	ServerURL string `json:"server_url" default:"http://localhost:9411" split_words:"true"`
}
