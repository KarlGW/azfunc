package bindings

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KarlGW/azfunc/data"
)

// HTTP represents an HTTP output binding.
type HTTP struct {
	header     http.Header
	name       string
	body       data.Raw
	statusCode int
}

// MarshalJSON implements custom marshaling to create the
// required JSON structure as expected by the function host.
func (b HTTP) MarshalJSON() ([]byte, error) {
	headers := make(map[string]string)
	for k, v := range b.header {
		headers[k] = v[0]
	}

	return json.Marshal(struct {
		Headers    map[string]string `json:"headers"`
		StatusCode string            `json:"statusCode"`
		Body       data.Raw          `json:"body"`
	}{
		Headers:    headers,
		StatusCode: strconv.Itoa(b.statusCode),
		Body:       b.body,
	})
}

// Data returns the data of the binding.
func (b HTTP) Data() data.Raw {
	return b.body
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
	b.body = data.Raw(d)
	return len(b.body), nil
}

// WriteHeader sets the response header with the provided
// status code.
func (b *HTTP) WriteHeader(statusCode int) {
	b.statusCode = statusCode
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

	b.statusCode = statusCode
	b.body = data.Raw(body)
	b.header = opts.Header
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
		statusCode: opts.StatusCode,
		body:       opts.Body,
		header:     opts.Header,
	}
}
