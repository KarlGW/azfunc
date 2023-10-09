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

// NewGeneric creates a new generic output binding.
func NewGeneric(name string, data []byte) Generic {
	return Generic{
		name: name,
		Raw:  data,
	}
}
