package eventgrid

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event represents an event (Event Grid schema).
type Event struct {
	Data        any       `json:"data"`
	Topic       string    `json:"topic"`
	Subject     string    `json:"subject"`
	Type        string    `json:"eventType"`
	Time        time.Time `json:"eventTime"`
	ID          string    `json:"id"`
	DataVersion string    `json:"dataVersion"`
}

// Schema returns the schema of the event.
func (s Event) Schema() Schema {
	return SchemaEventGrid
}

// JSON returns the JSON representation of the Event.
func (e Event) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}

// EventOptions contains options for an Event.
type EventOptions struct {
	Topic       string
	ID          string
	Time        time.Time
	DataVersion string
}

// EventOption is a function that sets options on an Event.
type EventGridEventOption func(o *EventOptions)

// NewEvent creates a new event (EventGrid schema).
func NewEvent(subject, eventType string, data any, options ...EventGridEventOption) (Event, error) {
	if len(subject) == 0 {
		return Event{}, fmt.Errorf("subject is required")
	}
	if len(eventType) == 0 {
		return Event{}, fmt.Errorf("eventType is required")
	}
	if data == nil {
		return Event{}, fmt.Errorf("data is required")
	}

	opts := EventOptions{
		DataVersion: "1.0",
	}
	for _, option := range options {
		option(&opts)
	}

	if len(opts.ID) == 0 {
		id, err := newUUID()
		if err != nil {
			return Event{}, err
		}
		opts.ID = id
	}
	if opts.Time.IsZero() {
		opts.Time = nowUTC()
	}

	return Event{
		Data:        data,
		Topic:       opts.Topic,
		Subject:     subject,
		Type:        eventType,
		Time:        opts.Time,
		ID:          opts.ID,
		DataVersion: opts.DataVersion,
	}, nil
}
