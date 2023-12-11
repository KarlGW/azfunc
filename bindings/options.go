package bindings

import (
	"net/http"
	"strings"

	"github.com/KarlGW/azfunc/data"
)

// Options for bindings. Not all options are viable for all
// bindings.
type Options struct {
	// Header sets the body of an HTTP binding.
	Header http.Header
	// Name sets the name of a binding.
	Name string
	// Body sets the body of an HTTP binding.
	Body data.Raw
	// Data sets the data of a base binding.
	Data data.Raw
	// StatusCode sets the status code on an HTTP binding.
	StatusCode int
}

// Option is a function that sets Options.
type Option func(o *Options)

// WithHeader adds the provided header to a HTTP binding.
func WithHeader(header http.Header) Option {
	return func(o *Options) {
		if o.Header == nil {
			o.Header = http.Header{}
		}
		for k, v := range header {
			o.Header.Add(k, strings.Join(v, ", "))
		}
	}
}
