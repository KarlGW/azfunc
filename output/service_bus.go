package output

import "github.com/KarlGW/azfunc/data"

// ServiceBus represents a service bus output binding.
type ServiceBus struct {
	name string
	data data.Raw
}

// ServiceBusOptions contains options for a ServiceBus output binding.
type ServiceBusOptions struct {
	// Name sets the name of the binding.
	Name string
	// Data sets the data of the binding.
	Data data.Raw
}

// ServiceBusOption is a function that sets options on a ServiceBus output binding.
type ServiceBusOption func(o *ServiceBusOptions)

// Data returns the data of the binding.
func (o ServiceBus) Data() data.Raw {
	return o.data
}

// Name returns the name of the binding.
func (o ServiceBus) Name() string {
	return o.name
}

// Write data to the binding.
func (o *ServiceBus) Write(d []byte) (int, error) {
	o.data = data.Raw(d)
	return len(o.data), nil
}

// NewServiceBus creates a new service bus output binding.
func NewServiceBus(name string, options ...ServiceBusOption) *ServiceBus {
	opts := ServiceBusOptions{}
	for _, option := range options {
		option(&opts)
	}
	return &ServiceBus{
		name: name,
		data: opts.Data,
	}
}
