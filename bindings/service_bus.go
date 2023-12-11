package bindings

import "github.com/KarlGW/azfunc/data"

// ServiceBus represents a service bus output binding.
type ServiceBus struct {
	name string
	data data.Raw
}

// Data returns the data of the binding.
func (b ServiceBus) Data() data.Raw {
	return b.data
}

// Name returns the name of the binding.
func (b ServiceBus) Name() string {
	return b.name
}

// Write data to the binding.
func (b *ServiceBus) Write(d []byte) (int, error) {
	b.data = data.Raw(d)
	return len(b.data), nil
}

// NewServiceBus creates a new service bus output binding.
func NewServiceBus(name string, options ...Option) *ServiceBus {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	return &ServiceBus{
		name: name,
		data: opts.Data,
	}
}
