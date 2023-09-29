package azfunc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// binding is the interface that wraps around method Name.
type binding interface {
	Name() string
}

// Output represents an outgoing response to the Functuon Host.
type Output struct {
	Outputs     map[string]binding
	Logs        []string
	ReturnValue any
}

// JSON returns the JSON encoding of Output.
func (r Output) JSON() []byte {
	b, _ := json.MarshalIndent(r, "", "  ")
	return b
}

// HTTPBinding represents an HTTP output binding.
type HTTPBinding struct {
	StatusCode string            `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

// Name returns the name of the binding. In case of an HTTP binding
// it is always "res".
func (b HTTPBinding) Name() string {
	return "res"
}

// NewHTTPBinding creates a new HTTP output binding.
func NewHTTPBinding(statusCode int, body []byte, header ...http.Header) HTTPBinding {
	hdr := make(map[string]string, len(header))
	for _, h := range header {
		for k, v := range h {
			hdr[k] = strings.Join(v, ", ")
		}
	}
	return HTTPBinding{
		StatusCode: strconv.Itoa(statusCode),
		Body:       string(body),
		Headers:    hdr,
	}
}

// Data is the data contained in a generic binding.
type Data []byte

// GenericBinding represents a generic Function App trigger. With custom handlers
// all bindings that are not HTTP output bindings share the same data structure.
type GenericBinding struct {
	name string
	Data
}

// Name returns the name of the binding.
func (b GenericBinding) Name() string {
	return b.name
}

// NewGenericBinding creates a new generic output binding.
func NewGenericBinding(name string, data Data) GenericBinding {
	return GenericBinding{
		name: name,
		Data: data,
	}
}

// NewOutput creates a new Output containing binding to be used for creating
// the response back to the Function host.
func NewOutput(bindings ...binding) (Output, error) {
	outputs := make(map[string]binding, len(bindings))
	for _, binding := range bindings {
		outputs[binding.Name()] = binding
	}

	return Output{
		Outputs: outputs,
	}, nil
}

// QueueBinding represents a Function App Queue Binding and contains
// the outgoing queue message data.
type QueueBinding = GenericBinding

// NewQueueBinding creates a new Queue output binding.
var NewQueueBinding = NewGenericBinding
