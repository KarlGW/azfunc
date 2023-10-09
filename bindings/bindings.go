package bindings

import (
	"encoding/json"
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

// Binding aliases.

// QueueBinding represents a Function App Queue Binding and contains
// the outgoing queue message data.
type QueueBinding = Generic

// NewQueueBinding creates a new Queue output binding.
var NewQueueBinding = NewGeneric
