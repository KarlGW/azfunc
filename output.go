package azfunc

import (
	"encoding/json"

	"github.com/KarlGW/azfunc/data"
	"github.com/KarlGW/azfunc/output"
)

// outputable is the interface that wraps around methods Data, Name and Write.
type outputable interface {
	// Data returns the data of the binding.
	Data() data.Raw
	// Name returns the name of the binding.
	Name() string
	// Write to the binding.
	Write([]byte) (int, error)
}

// outputs represents an outgoing response to the Function Host.
// Make private? This is not used by exposed functionality.
type outputs struct {
	outputs     map[string]outputable
	returnValue any
	http        *output.HTTP
	logs        []string
}

// MarshalJSON implements custom marshaling to produce
// the required JSON structure as expected by the
// function host.
func (o outputs) MarshalJSON() ([]byte, error) {
	type Alias outputs

	temp := struct {
		Outputs     map[string]any `json:"Outputs"`
		ReturnValue any            `json:"ReturnValue"`
		Logs        []string       `json:"Logs"`
		Alias
	}{
		Outputs: make(map[string]any),
		Alias:   Alias(o),
	}

	for key, binding := range o.outputs {
		if b, ok := binding.(*output.HTTP); ok {
			temp.Outputs[key] = b
		} else {
			temp.Outputs[key] = binding.Data()
		}
	}

	return json.Marshal(temp)
}

// json returns the JSON encoding of Output.
func (o outputs) json() []byte {
	if o.http != nil {
		o.outputs[o.http.Name()] = o.http
	}
	b, _ := json.Marshal(o)
	return b
}

// addOutputs one or more output bindings to functionOutput.
func (o *outputs) Add(outputs ...outputable) {
	if o.outputs == nil {
		o.outputs = make(map[string]outputable, len(outputs))
	}

	for _, binding := range outputs {
		if b, ok := binding.(*output.HTTP); ok {
			o.http = b
		} else {
			o.outputs[binding.Name()] = binding
		}
	}
}

// SetReturnValue sets ReturnValue of Output.
func (o *outputs) SetReturnValue(v any) {
	o.returnValue = v
}

// Get returns the output binding with the provided name, if no output binding
// with that name exists, return a new generic output binding with the
// provided name.
func (o outputs) Get(name string) outputable {
	binding, ok := o.outputs[name]
	if !ok {
		o.outputs[name] = output.NewGeneric(name)
		return o.outputs[name]
	}
	return binding
}

// Binding returns the output binding with the provided name, if no output binding
// with that name exists, return a new generic output binding with the
// provided name.
func (o outputs) Binding(name string) outputable {
	return o.Get(name)
}

// Output returns the output binding with the provided name, if no output binding
// with that name exists, return a new generic output binding with the
// provided name.
func (o outputs) Output(name string) outputable {
	return o.Get(name)
}

// HTTP returns the HTTP output binding if any is set.
// If not set it will create, set and return it.
func (o *outputs) HTTP() *output.HTTP {
	if o.http == nil {
		o.http = output.NewHTTP()
		return o.http
	}
	return o.http
}

// outputsOptions contains options for creating a new
// Output.
type outputsOptions struct {
	returnValue any
	http        *output.HTTP
	outputs     []outputable
	Logs        []string
}

// outputsOption is a function that sets OutputOptions.
type outputsOption func(o *outputsOptions)

// newOutputs creates a new outputs containing output bindings to be used for creating
// the response back to the Function host.
// make private?
func newOutputs(options ...outputsOption) *outputs {
	opts := outputsOptions{}
	for _, option := range options {
		option(&opts)
	}

	var logs []string
	if len(opts.Logs) > 0 {
		logs = make([]string, len(opts.Logs))
		copy(logs, opts.Logs)
	}

	outputs := &outputs{
		logs:        logs,
		returnValue: opts.returnValue,
		http:        opts.http,
	}
	outputs.Add(opts.outputs...)

	return outputs
}

// withOutputs add one or more output bindings to OutputOptions
func withOutputs(outputs ...outputable) outputsOption {
	return func(o *outputsOptions) {
		for _, binding := range outputs {
			if b, ok := binding.(*output.HTTP); ok {
				o.http = b
			} else {
				o.outputs = append(o.outputs, binding)
			}
		}
	}
}
