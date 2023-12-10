package bindings

// Bindable is the interface that wraps around methods Name and Write.
type Bindable interface {
	// Name returns the name of the binding.
	Name() string
	// Write to the binding.
	Write([]byte) (int, error)
}
