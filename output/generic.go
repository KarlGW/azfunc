package output

import "github.com/KarlGW/azfunc/data"

// Generic represents a generic output binding. With custom handlers
// all bindings that are not HTTP output bindings share the
// same data structure.
type Generic struct {
	name string
	data data.Raw
}

// GenericOptions contains options for a Generic output binding.
type GenericOptions struct {
	// Name sets the name of the binding.
	Name string
	// Data sets the data of the binding.
	Data data.Raw
}

// GenericOption is a function that sets options on a Generic output binding
type GenericOption func(o *GenericOptions)

// Data returns the data of the binding.
func (o Generic) Data() data.Raw {
	return o.data
}

// Name returns the name of the binding.
func (o Generic) Name() string {
	return o.name
}

// Write data to the binding.
func (o *Generic) Write(d []byte) (int, error) {
	o.data = data.Raw(d)
	return len(o.data), nil
}

// NewGeneric creates a new generic output binding.
func NewGeneric(name string, options ...GenericOption) *Generic {
	opts := GenericOptions{}
	for _, option := range options {
		option(&opts)
	}
	return &Generic{
		name: name,
		data: opts.Data,
	}
}
