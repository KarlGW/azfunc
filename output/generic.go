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
func (b Generic) Data() data.Raw {
	return b.data
}

// Name returns the name of the binding.
func (b Generic) Name() string {
	return b.name
}

// Write data to the binding.
func (b *Generic) Write(d []byte) (int, error) {
	b.data = data.Raw(d)
	return len(b.data), nil
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
