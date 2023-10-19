package bindings

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/KarlGW/azfunc/data"
)

// HTTP represents an HTTP output binding.
type HTTP struct {
	StatusCode int
	Body       data.Raw
	header     http.Header
	name       string
}

// Name returns the name of the binding. In case of an HTTP binding
// it is always "res".
func (b HTTP) Name() string {
	if len(b.name) > 0 {
		return b.name
	}
	return "res"
}

// Write the provided data to the body of the HTTP bindings.
func (b *HTTP) Write(d []byte) (int, error) {
	b.Body = data.Raw(d)
	return len(b.Body), nil
}

// WriteHeader sets the response header with the provided
// status code.
func (b *HTTP) WriteHeader(statusCode int) {
	b.StatusCode = statusCode
}

// Header returns the header(s) of the HTTP binding.
func (b *HTTP) Header() http.Header {
	if b.header == nil {
		b.header = http.Header{}
		return b.header
	}
	return b.header
}

// WriteResponse writes the provided status code, body and options to
// the HTTP binding. Supports option WithHeader.
func (b *HTTP) WriteResponse(statusCode int, body []byte, options ...Option) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	b.StatusCode = statusCode
	b.Body = data.Raw(body)
	b.header = opts.Header
}

// IsZero checks if the HTTP binding is unset.
func (b HTTP) IsZero() bool {
	return b.StatusCode == 0 && b.Body == nil && b.header == nil
}

// toHTTPBinding returns the httpBinding representation
// of HTTP.
func (b HTTP) toHTTPBinding() *httpBinding {
	headers := make(map[string]string)
	for k, v := range b.header {
		headers[k] = strings.Join(v, ", ")
	}
	return &httpBinding{
		StatusCode: strconv.Itoa(b.StatusCode),
		Body:       b.Body,
		Headers:    headers,
	}
}

// NewHTTP creates a new HTTP output binding.
func NewHTTP(options ...Option) *HTTP {
	opts := Options{
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
		StatusCode: opts.StatusCode,
		Body:       opts.Body,
		header:     opts.Header,
	}
}

// httpBinding is the Output binding representation of an
// HTTP binding.
type httpBinding struct {
	StatusCode string            `json:"statusCode"`
	Body       data.Raw          `json:"body"`
	Headers    map[string]string `json:"headers"`
	name       string
}

// Name returns the name of the binding. In case of an HTTP binding
// it is always "res".
func (b httpBinding) Name() string {
	if len(b.name) > 0 {
		return b.name
	}
	return "res"
}

// Write the provided data to the body of the HTTP bindings.
func (b *httpBinding) Write(d []byte) (int, error) {
	b.Body = data.Raw(d)
	return len(b.Body), nil
}
