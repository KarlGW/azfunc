package azfunc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Bindable is the interface that wraps around method Name.
type Bindable interface {
	Name() string
}

// Output represents an outgoing response to the Functuon Host.
type Output struct {
	Outputs     map[string]Bindable
	Logs        []string
	ReturnValue any
}

// JSON returns the JSON encoding of Output.
func (o Output) JSON() []byte {
	b, _ := json.Marshal(o)
	return b
}

// AddBindings one or more bindings to Output.
func (o *Output) AddBindings(bindings ...Bindable) {
	if o.Outputs == nil {
		o.Outputs = make(map[string]Bindable, len(bindings))
	}

	for _, binding := range bindings {
		o.Outputs[binding.Name()] = binding
	}
}

// Log adds a message to the Logs of Output.
func (o *Output) Log(msg string) {
	if o.Logs == nil {
		o.Logs = make([]string, 0)
	}
	o.Logs = append(o.Logs, msg)
}

// SetReturnValue sets ReturnValue of Output.
func (o *Output) SetReturnValue(v any) {
	o.ReturnValue = v
}

// OutputOptions contains options for creating a new
// Output.
type OutputOptions struct {
	Bindings    []Bindable
	Logs        []string
	ReturnValue any
}

// Output option is a function that sets OutputOptions.
type OutputOption func(o *OutputOptions)

// WithBindings add one or more bindings to OutputOptions
func WithBindings(bindings ...Bindable) OutputOption {
	return func(o *OutputOptions) {
		o.Bindings = bindings
	}
}

// NewOutput creates a new Output containing binding to be used for creating
// the response back to the Function host.
func NewOutput(options ...OutputOption) Output {
	opts := OutputOptions{}
	for _, option := range options {
		option(&opts)
	}

	var logs []string
	if len(opts.Logs) > 0 {
		logs = make([]string, len(opts.Logs))
		copy(logs, opts.Logs)
	}

	output := Output{
		Logs:        logs,
		ReturnValue: opts.ReturnValue,
	}
	output.AddBindings(opts.Bindings...)

	return output
}

// HTTPBinding represents an HTTP output binding.
type HTTPBinding struct {
	StatusCode string            `json:"statusCode"`
	Body       Payload           `json:"body"`
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
		Body:       body,
		Headers:    hdr,
	}
}

// GenericBinding represents a generic Function App trigger. With custom handlers
// all bindings that are not HTTP output bindings share the same data structure.
type GenericBinding struct {
	name string
	Payload
}

// Name returns the name of the binding.
func (b GenericBinding) Name() string {
	return b.name
}

// NewGenericBinding creates a new generic output binding.
func NewGenericBinding(name string, data []byte) GenericBinding {
	return GenericBinding{
		name:    name,
		Payload: data,
	}
}

// Binding aliases.

// QueueBinding represents a Function App Queue Binding and contains
// the outgoing queue message data.
type QueueBinding = GenericBinding

// NewQueueBinding creates a new Queue output binding.
var NewQueueBinding = NewGenericBinding
