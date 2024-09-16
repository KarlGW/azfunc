package output

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/KarlGW/azfunc/data"
	"github.com/KarlGW/azfunc/uuid"
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
	Data        any       `json:"data,omitempty"`
	SpecVersion string    `json:"specversion"`
	Type        string    `json:"type"`
	Source      string    `json:"source"`
	ID          string    `json:"id"`
	Time        time.Time `json:"time"`
	Subject     string    `json:"subject,omitempty"`
	DataSchema  string    `json:"dataschema,omitempty"`
}

// JSON returns the JSON representation of the CloudEvent.
func (e CloudEvent) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}

// CloudEventOptions contains options for a CloudEvent.
type CloudEventOptions struct {
	Subject    string
	DataSchema string
}

// CloudEventOption is a function that sets options on a CloudEvent.
type CloudEventOption func(o *CloudEventOptions)

// NewCloudEvent creates a new CloudEvent.
func NewCloudEvent(source, eventType string, data any, options ...CloudEventOption) (CloudEvent, error) {
	if len(source) == 0 {
		return CloudEvent{}, fmt.Errorf("source is required")
	}
	if len(eventType) == 0 {
		return CloudEvent{}, fmt.Errorf("type (eventType) is required")
	}

	opts := CloudEventOptions{}
	for _, option := range options {
		option(&opts)
	}

	id, err := uuid.New()
	if err != nil {
		return CloudEvent{}, err
	}

	return CloudEvent{
		Data:        data,
		SpecVersion: "1.0",
		Type:        eventType,
		Source:      source,
		ID:          id,
		Time:        time.Now().UTC(),
		Subject:     opts.Subject,
		DataSchema:  opts.DataSchema,
	}, nil
}

// EventGridEvent represents an event (EventGrid schema).
type EventGridEvent struct {
	Data        any       `json:"data"`
	Topic       string    `json:"topic"`
	Subject     string    `json:"subject"`
	EventType   string    `json:"eventType"`
	EventTime   time.Time `json:"eventTime"`
	ID          string    `json:"id"`
	DataVersion string    `json:"dataVersion"`
}

// JSON returns the JSON representation of the EventGridEvent.
func (e EventGridEvent) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}

// EventGridEventOptions contains options for an EventGridEvent.
type EventGridEventOptions struct {
	Topic string
}

// EventGridEventOption is a function that sets options on an EventGridEvent.
type EventGridEventOption func(o *EventGridEventOptions)

// NewEventGridEvent creates a new EventGridEvent (EventGrid schema).
func NewEventGridEvent(subject, eventType string, data any, options ...EventGridEventOption) (EventGridEvent, error) {
	if len(subject) == 0 {
		return EventGridEvent{}, fmt.Errorf("subject is required")
	}
	if len(eventType) == 0 {
		return EventGridEvent{}, fmt.Errorf("eventType is required")
	}
	if data == nil {
		return EventGridEvent{}, fmt.Errorf("data is required")
	}

	opts := EventGridEventOptions{}
	for _, option := range options {
		option(&opts)
	}

	id, err := uuid.New()
	if err != nil {
		return EventGridEvent{}, err
	}

	return EventGridEvent{
		Data:        data,
		Topic:       opts.Topic,
		Subject:     subject,
		EventType:   eventType,
		EventTime:   time.Now().UTC(),
		ID:          id,
		DataVersion: "1.0",
	}, nil
}
