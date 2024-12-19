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
type outputs struct {
	outputs     map[string]outputable
	log         invocationLogger
	http        *output.HTTP
	returnValue any
}

// MarshalJSON implements custom marshaling to produce
// the required JSON structure as expected by the
// function host.
func (o outputs) MarshalJSON() ([]byte, error) {
	temp := struct {
		Outputs     map[string]any `json:"Outputs"`
		ReturnValue any            `json:"ReturnValue"`
		Logs        []string       `json:"Logs"`
	}{
		Outputs:     make(map[string]any),
		ReturnValue: o.returnValue,
		Logs:        o.log.Entries(),
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

// outputsOptions contains options for creating a new
// Output.
type outputsOptions struct {
	http    *output.HTTP
	outputs []outputable
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

	outputs := &outputs{
		http: opts.http,
		log:  newInvocationLogger(),
	}
	outputs.Add(opts.outputs...)

	return outputs
}

// json returns the JSON encoding of Output.
func (o outputs) json() []byte {
	if o.http != nil {
		o.outputs[o.http.Name()] = o.http
	}
	b, _ := json.Marshal(o)
	return b
}

// Add one or more output bindings to functionOutput.
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

// Log returns the invocation logger. The invocation logger writes
// to the outputs logs field that is used by the function host
// to handle logging. The log entries are not written until
// the function has run to completion.
func (o outputs) Log() InvocationLogger {
	return o.log
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
