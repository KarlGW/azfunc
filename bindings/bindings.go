package bindings

// Bindable is the interface that wraps around methods Name and Write.
type Bindable interface {
	// Name returns the name of the binding.
	Name() string
	// Write to the binding.
	Write([]byte) (int, error)
}

// Binding aliases.

// Queue represents a Function App Queue Binding and contains
// the outgoing queue message data.
type Queue = Base

// NewQueue creates a new Queue output binding.
var NewQueue = NewBase
