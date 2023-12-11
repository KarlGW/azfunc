package bindings

import (
	"encoding/json"
)

// Output represents an outgoing response to the Functuon Host.
type Output struct {
	Outputs     map[string]Bindable
	Logs        []string
	ReturnValue any
	http        *HTTP
}

// MarshalJSON implements custom marshaling to produce
// the required JSON structure as expected by the
// function host.
func (o Output) MarshalJSON() ([]byte, error) {
	type Alias Output

	temp := struct {
		Outputs map[string]any `json:"Outputs"`
		Alias
	}{
		Outputs: make(map[string]any),
		Alias:   Alias(o),
	}

	for key, binding := range o.Outputs {
		if b, ok := binding.(*httpBinding); ok {
			temp.Outputs[key] = b
		} else {
			temp.Outputs[key] = binding.Data()
		}
	}

	return json.Marshal(temp)
}

// JSON returns the JSON encoding of Output.
func (o Output) JSON() []byte {
	if o.http != nil && !o.http.IsZero() {
		o.Outputs[o.http.Name()] = o.http.toHTTPBinding()
	}
	b, _ := json.Marshal(o)
	return b
}

// AddBindings one or more bindings to Output.
func (o *Output) AddBindings(bindings ...Bindable) {
	if o.Outputs == nil {
		o.Outputs = make(map[string]Bindable, len(bindings))
	}

	for _, binding := range bindings {
		if b, ok := binding.(*HTTP); ok {
			o.http = b
		} else {
			o.Outputs[binding.Name()] = binding
		}
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

// Binding returns the binding with the provided name, if no binding
// with that name exists, return a new base binding with the
// provided name.
func (o Output) Binding(name string) Bindable {
	binding, ok := o.Outputs[name]
	if !ok {
		o.Outputs[name] = NewBase(name)
		return o.Outputs[name]
	}
	return binding
}

// HTTP returns the HTTP binding of output if any is set.
// If not set it will create, set and return it.
func (o *Output) HTTP() *HTTP {
	if o.http == nil || o.http.IsZero() {
		o.http = NewHTTP()
		return o.http
	}
	return o.http
}

// OutputOptions contains options for creating a new
// Output.
type OutputOptions struct {
	Bindings    []Bindable
	Logs        []string
	ReturnValue any
	http        *HTTP
}

// Output option is a function that sets OutputOptions.
type OutputOption func(o *OutputOptions)

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
		http:        opts.http,
	}
	output.AddBindings(opts.Bindings...)

	return output
}

// WithBindings add one or more bindings to OutputOptions
func WithBindings(bindings ...Bindable) OutputOption {
	return func(o *OutputOptions) {
		for _, binding := range bindings {
			if b, ok := binding.(*HTTP); ok {
				o.http = b
			} else {
				o.Bindings = append(o.Bindings, binding)
			}
		}
	}
}
