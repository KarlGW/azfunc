package bindings

import "github.com/KarlGW/azfunc/data"

// Queue represents a Queue Storage output binding.
type Queue struct {
	name string
	data data.Raw
}

// QueueOptions contains options for a Queue Storage output binding.
type QueueOptions struct {
	// Name sets the name of the binding.
	Name string
	// Data sets the data of the binding.
	Data data.Raw
}

// QueueOption is a function that sets options on a Queue Storage output binding.
type QueueOption func(o *QueueOptions)

// Data returns the data of the binding.
func (b Queue) Data() data.Raw {
	return b.data
}

// Name returns the name of the binding.
func (b Queue) Name() string {
	return b.name
}

// Write data to the binding.
func (b *Queue) Write(d []byte) (int, error) {
	b.data = data.Raw(d)
	return len(b.data), nil
}

// NewQueue creates a new queue storage output binding.
func NewQueue(name string, options ...QueueOption) *Queue {
	opts := QueueOptions{}
	for _, option := range options {
		option(&opts)
	}
	return &Queue{
		name: name,
		data: opts.Data,
	}
}
