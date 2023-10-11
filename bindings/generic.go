package bindings

import "github.com/KarlGW/azfunc/data"

// Generic represents a generic Function App trigger. With custom handlers
// all bindings that are not HTTP output bindings share the same data structure.
type Generic struct {
	name string
	data.Raw
}

// Name returns the name of the binding.
func (b Generic) Name() string {
	return b.name
}

// Write data to the binding.
func (b *Generic) Write(d []byte) (int, error) {
	b.Raw = data.Raw(d)
	return len(b.Raw), nil
}

// NewGeneric creates a new generic output binding.
func NewGeneric(name string, options ...Option) *Generic {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	return &Generic{
		name: name,
		Raw:  opts.Data,
	}
}
