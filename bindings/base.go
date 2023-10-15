package bindings

import "github.com/KarlGW/azfunc/data"

// Base represents a base Function App trigger. With custom handlers
// all bindings that are not HTTP output bindings share the same data structure.
type Base struct {
	name string
	data.Raw
}

// Name returns the name of the binding.
func (b Base) Name() string {
	return b.name
}

// Write data to the binding.
func (b *Base) Write(d []byte) (int, error) {
	b.Raw = data.Raw(d)
	return len(b.Raw), nil
}

// NewBase creates a new base output binding.
func NewBase(name string, options ...Option) *Base {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	return &Base{
		name: name,
		Raw:  opts.Data,
	}
}
