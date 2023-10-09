package bindings

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/KarlGW/azfunc/data"
)

// HTTP represents an HTTP output binding.
type HTTP struct {
	StatusCode string            `json:"statusCode"`
	Body       data.Raw          `json:"body"`
	Headers    map[string]string `json:"headers"`
}

// Name returns the name of the binding. In case of an HTTP binding
// it is always "res".
func (b HTTP) Name() string {
	return "res"
}

// NewHTTP creates a new HTTP output binding.
func NewHTTP(statusCode int, body []byte, header ...http.Header) HTTP {
	hdr := make(map[string]string, len(header))
	for _, h := range header {
		for k, v := range h {
			hdr[k] = strings.Join(v, ", ")
		}
	}
	return HTTP{
		StatusCode: strconv.Itoa(statusCode),
		Body:       body,
		Headers:    hdr,
	}
}
