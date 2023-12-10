package bindings

import "github.com/KarlGW/azfunc/data"

// ServiceBus represents a service bus output binding.
type ServiceBus struct {
	name string
	data.Raw
}

// Name returns the name of the binding.
func (b ServiceBus) Name() string {
	return b.name
}

// Write data to the binding.
func (b *ServiceBus) Write(d []byte) (int, error) {
	b.Raw = data.Raw(d)
	return len(b.Raw), nil
}

// NewServiceBus creates a new service bus output binding.
func NewServiceBus(name string, options ...Option) *ServiceBus {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	return &ServiceBus{
		name: name,
		Raw:  opts.Data,
	}
}
