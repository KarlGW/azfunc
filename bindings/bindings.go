package bindings

// Bindable is the interface that wraps around method Name.
type Bindable interface {
	// Name returns the name of the binding.
	Name() string
	// Write to the binding.
	Write([]byte) (int, error)
}

// Binding aliases.

// QueueBinding represents a Function App Queue Binding and contains
// the outgoing queue message data.
type QueueBinding = Generic

// NewQueueBinding creates a new Queue output binding.
var NewQueueBinding = NewGeneric
