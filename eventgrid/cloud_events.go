package eventgrid

import (
	"encoding/json"
	"fmt"
	"time"
)

// CloudEvent represents a CloudEvent.
type CloudEvent struct {
	Time        time.Time `json:"time"`
	Data        any       `json:"data,omitempty"`
	SpecVersion string    `json:"specversion"`
	Type        string    `json:"type"`
	Source      string    `json:"source"`
	ID          string    `json:"id"`
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
	Time        time.Time
	ID          string
	Subject     string
	DataSchema  string
	SpecVersion string
}

// CloudEventOption is a function that sets options on a CloudEvent.
type CloudEventOption func(o *CloudEventOptions)

// NewCloudEvent creates a new CloudEvent. By default a new UUID is generated
// for the ID, the current time is used for the Time and specversion is set to
// "1.0". This can be overridden by providing options.
func NewCloudEvent(source, eventType string, data any, options ...CloudEventOption) (CloudEvent, error) {
	if len(source) == 0 {
		return CloudEvent{}, fmt.Errorf("source is required")
	}
	if len(eventType) == 0 {
		return CloudEvent{}, fmt.Errorf("type (eventType) is required")
	}

	opts := CloudEventOptions{
		SpecVersion: "1.0",
	}
	for _, option := range options {
		option(&opts)
	}

	if len(opts.ID) == 0 {
		id, err := newUUID()
		if err != nil {
			return CloudEvent{}, err
		}
		opts.ID = id
	}
	if opts.Time.IsZero() {
		opts.Time = nowUTC()
	}

	return CloudEvent{
		Data:        data,
		SpecVersion: opts.SpecVersion,
		Type:        eventType,
		Source:      source,
		ID:          opts.ID,
		Time:        opts.Time,
		Subject:     opts.Subject,
		DataSchema:  opts.DataSchema,
	}, nil
}
