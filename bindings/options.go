package bindings

import (
	"net/http"

	"github.com/KarlGW/azfunc/data"
)

// Options for bindings. Not all options are viable for all
// bindings.
type Options struct {
	// StatusCode sets the status code on an HTTP binding.
	StatusCode int
	// Body sets the body of an HTTP binding.
	Body data.Raw
	// Header sets the body of an HTTP binding.
	Header http.Header
	// Data sets the data of a base binding.
	Data data.Raw
	// Name sets the name of a binding.
	Name string
}

// Option is a function that sets Options.
type Option func(o *Options)
