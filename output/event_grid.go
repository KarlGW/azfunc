package output

import (
	"github.com/potatoattack/azfunc/data"
)

// EventGrid represents an Event Grid output binding.
type EventGrid struct {
	name string
	data data.Raw
}

// EventGridOptions contains options for an Event Grid output binding.
type EventGridOptions struct {
	// Name sets the name of the binding.
	Name string
	// Data sets the data of the binding.
	Data data.Raw
}

// EventGridOption is a function that sets options on an Event Grid output binding.
type EventGridOption func(o *EventGridOptions)

// Data returns the data of the binding.
func (o EventGrid) Data() data.Raw {
	return o.data
}

// Name returns the name of the binding.
func (o EventGrid) Name() string {
	return o.name
}

// Write data to the binding.
func (o *EventGrid) Write(d []byte) (int, error) {
	o.data = data.Raw(d)
	return len(o.data), nil
}

// NewEventGrid creates a new Event Grid output binding.
func NewEventGrid(name string, options ...EventGridOption) *EventGrid {
	opts := EventGridOptions{}
	for _, option := range options {
		option(&opts)
	}
	return &EventGrid{
		name: name,
		data: opts.Data,
	}
}
