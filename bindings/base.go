package bindings

import "github.com/KarlGW/azfunc/data"

// Base represents a base output binding. With custom handlers
// all bindings that are not HTTP output bindings share the
// same data structure.
type Base struct {
	name string
	data data.Raw
}

// BaseOptions contains options for a Base output binding.
type BaseOptions struct {
	// Name sets the name of the binding.
	Name string
	// Data sets the data of the binding.
	Data data.Raw
}

// BaseOption is a function that sets options on a Base output binding
type BaseOption func(o *BaseOptions)

// Data returns the data of the binding.
func (b Base) Data() data.Raw {
	return b.data
}

// Name returns the name of the binding.
func (b Base) Name() string {
	return b.name
}

// Write data to the binding.
func (b *Base) Write(d []byte) (int, error) {
	b.data = data.Raw(d)
	return len(b.data), nil
}

// NewBase creates a new base output binding.
func NewBase(name string, options ...BaseOption) *Base {
	opts := BaseOptions{}
	for _, option := range options {
		option(&opts)
	}
	return &Base{
		name: name,
		data: opts.Data,
	}
}
