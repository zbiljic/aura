package corsx

import "github.com/rs/cors"

type Config struct {
	// Enable CORS
	// If set to true, CORS will be enabled and preflight-requests (OPTION) will be answered.
	Enabled bool `json:"enabled"`

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string `json:"allowed_origins" split_words:"true"`
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string `json:"allowed_methods" split_words:"true"`
	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string `json:"allowed_headers" split_words:"true"`
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string `json:"exposed_headers" split_words:"true"`
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int `json:"max_age" split_words:"true"`
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool `json:"allow_credentials" split_words:"true"`
	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool `json:"options_passthrough" split_words:"true"`
	// Debugging flag adds additional output to debug server side CORS issues
	Debug bool `json:"debug"`
}

func (c *Config) CORSOptions() (cors.Options, bool) {
	return cors.Options{
		AllowedOrigins:     c.AllowedOrigins,
		AllowedMethods:     c.AllowedMethods,
		AllowedHeaders:     c.AllowedHeaders,
		ExposedHeaders:     c.ExposedHeaders,
		MaxAge:             c.MaxAge,
		AllowCredentials:   c.AllowCredentials,
		OptionsPassthrough: c.OptionsPassthrough,
		Debug:              c.Debug,
	}, c.Enabled
}
