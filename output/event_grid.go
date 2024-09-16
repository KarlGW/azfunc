package output

import (
	"encoding/json"
	"time"

	"github.com/KarlGW/azfunc/data"
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
func (b EventGrid) Data() data.Raw {
	return b.data
}

// Name returns the name of the binding.
func (b EventGrid) Name() string {
	return b.name
}

// Write data to the binding.
func (b *EventGrid) Write(d []byte) (int, error) {
	b.data = data.Raw(d)
	return len(b.data), nil
}

// NewEventGrid creates a new event grid output binding.
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

// CloudEvent represents a CloudEvent.
type CloudEvent struct {
	Data        any       `json:"data"`
	SpecVersion string    `json:"specversion"`
	Type        string    `json:"type"`
	Source      string    `json:"source"`
	ID          string    `json:"id"`
	Time        time.Time `json:"time"`
	Subject     string    `json:"subject"`
	DataSchema  string    `json:"dataschema"`
}

// JSON returns the JSON representation of the CloudEvent.
func (e CloudEvent) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}

// NewCloudEvent creates a new CloudEvent.
func NewCloudEvent() CloudEvent {
	return CloudEvent{}
}

// EventGridEvent represents an event (EventGrid schema).
type EventGridEvent struct {
	Data            any       `json:"data"`
	Topic           string    `json:"topic"`
	Subject         string    `json:"subject"`
	EventType       string    `json:"eventType"`
	EventTime       time.Time `json:"eventTime"`
	ID              string    `json:"id"`
	DataVersion     string    `json:"dataVersion"`
	MetadataVersion string    `json:"metadataVersion"`
}

// JSON returns the JSON representation of the EventGridEvent.
func (e EventGridEvent) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}

// NewEventGridEvent creates a new EventGridEvent (EventGrid schema).
func NewEventGridEvent() EventGridEvent {
	return EventGridEvent{}
}

// NewEvent creates a new event. Default schema is CloudEvents.
func NewEvent() data.JSONMarshaler {
	return nil
}
