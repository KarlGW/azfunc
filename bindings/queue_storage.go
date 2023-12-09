package bindings

import "github.com/KarlGW/azfunc/data"

// QueueStorage represents a queue storage output binding.
type QueueStorage struct {
	name string
	data.Raw
}

// Name returns the name of the binding.
func (b QueueStorage) Name() string {
	return b.name
}

// Write data to the binding.
func (b *QueueStorage) Write(d []byte) (int, error) {
	b.Raw = data.Raw(d)
	return len(b.Raw), nil
}

// NewQueueStorage creates a new queue storage output binding.
func NewQueueStorage(name string, options ...Option) *QueueStorage {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	return &QueueStorage{
		name: name,
		Raw:  opts.Data,
	}
}
