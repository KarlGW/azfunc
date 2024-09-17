package output

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/KarlGW/azfunc/data"
)

// HTTP represents an HTTP output binding.
type HTTP struct {
	header     http.Header
	name       string
	body       data.Raw
	statusCode int
}

// HTTPOptions contains options for an HTTP output binding.
type HTTPOptions struct {
	// Header sets the body of the binding.
	Header http.Header
	// Name sets the name of the binding.
	Name string
	// Body sets the body of the binding.
	Body data.Raw
	// StatusCode sets the status code of the binding.
	StatusCode int
}

// HTTPOption is a function that sets options on an HTTP output binding.
type HTTPOption func(o *HTTPOptions)

// MarshalJSON implements custom marshaling to create the
// required JSON structure as expected by the function host.
func (o HTTP) MarshalJSON() ([]byte, error) {
	headers := make(map[string]string, len(o.header))
	for k, v := range o.header {
		headers[k] = v[0]
	}

	return json.Marshal(struct {
		Headers    map[string]string `json:"headers"`
		StatusCode string            `json:"statusCode"`
		Body       data.Raw          `json:"body"`
	}{
		Headers:    headers,
		StatusCode: strconv.Itoa(o.statusCode),
		Body:       o.body,
	})
}

// Data returns the data of the binding.
func (o HTTP) Data() data.Raw {
	return o.body
}

// Name returns the name of the binding. In case of an HTTP binding
// it is always "res".
func (o HTTP) Name() string {
	if len(o.name) > 0 {
		return o.name
	}
	return "res"
}

// Write the provided data to the body of the HTTP bindings.
func (o *HTTP) Write(d []byte) (int, error) {
	o.body = data.Raw(d)
	return len(o.body), nil
}

// WriteHeader sets the response header with the provided
// status code.
func (o *HTTP) WriteHeader(statusCode int) {
	o.statusCode = statusCode
}

// Header returns the header(s) of the HTTP binding.
func (o *HTTP) Header() http.Header {
	if o.header == nil {
		o.header = http.Header{}
		return o.header
	}
	return o.header
}

// WriteResponse writes the provided status code, body and options to
// the HTTP binding. Supports option WithHeader.
func (o *HTTP) WriteResponse(statusCode int, body []byte, options ...HTTPOption) {
	opts := HTTPOptions{}
	for _, option := range options {
		option(&opts)
	}

	o.statusCode = statusCode
	o.body = data.Raw(body)
	o.header = opts.Header
}

// NewHTTP creates a new HTTP output binding.
func NewHTTP(options ...HTTPOption) *HTTP {
	opts := HTTPOptions{
		Header: http.Header{},
	}
	for _, option := range options {
		option(&opts)
	}
	if len(opts.Name) == 0 {
		opts.Name = "res"
	}
	if opts.StatusCode < http.StatusContinue || opts.StatusCode > http.StatusNetworkAuthenticationRequired {
		opts.StatusCode = http.StatusOK
	}

	return &HTTP{
		name:       opts.Name,
		statusCode: opts.StatusCode,
		body:       opts.Body,
		header:     opts.Header,
	}
}

// WithHeader adds the provided header to a HTTP binding.
func WithHeader(header http.Header) HTTPOption {
	return func(o *HTTPOptions) {
		if o.Header == nil {
			o.Header = http.Header{}
		}
		for k, v := range header {
			o.Header.Add(k, strings.Join(v, ", "))
		}
	}
}
