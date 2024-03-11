package bindings

import "github.com/KarlGW/azfunc/data"

// Bindable is the interface that wraps around methods Data, Name and Write.
type Bindable interface {
	// Data returns the data of the binding.
	Data() data.Raw
	// Name returns the name of the binding.
	Name() string
	// Write to the binding.
	Write([]byte) (int, error)
}
