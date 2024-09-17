package trigger

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KarlGW/azfunc/data"
	"github.com/KarlGW/azfunc/eventgrid"
)

// EventGrid represents an Event Grid trigger. It handles both
// Cloud Events and Event Grid events. CloudEvents contains
// source where as an Event Grid event contains topic.
type EventGrid struct {
	Time     time.Time
	Metadata EventGridMetadata
	ID       string
	// Topic is the topic of the event. It is the same as the source
	// for Cloud Events.
	Topic string
	// Source is the source of the event. It is the same as the topic
	// for Event Grid events.
	Source  string
	Subject string
	Type    string
	Data    data.Raw
	Schema  eventgrid.Schema
}

// EventGridOptions contains options for an Event Grid trigger.
type EventGridOptions struct {
	Schema eventgrid.Schema
}

// EventGridOption is a function that sets options on an Event Grid
// trigger.
type EventGridOption func(o *EventGridOptions)

// EventGridMetadata represents the metadata for an Event Grid trigger.
type EventGridMetadata struct {
	Metadata
	Data data.Raw `json:"data"`
}

// Parse the data from the Event Grid trigger into the provided value.
func (t EventGrid) Parse(v any) error {
	return json.Unmarshal(t.Data, &v)
}

// NewEventGrid creates and returns an Event Grid trigger from the provided
// *http.Request.
func NewEventGrid(r *http.Request, name string, options ...EventGridOption) (*EventGrid, error) {
	opts := EventGridOptions{}
	for _, option := range options {
		option(&opts)
	}

	var t eventGridTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, ErrTriggerPayloadMalformed
	}

	d, ok := t.Data[name]
	if !ok {
		return nil, ErrTriggerNameIncorrect
	}

	var eventType string
	var eventTime time.Time
	var schema eventgrid.Schema
	if len(d.SpecVersion) > 0 {
		eventType = d.Type
		eventTime = d.Time
		schema = eventgrid.SchemaCloudEvents
		d.Topic = d.Source
	} else if len(d.EventType) > 0 {
		eventType = d.EventType
		eventTime = d.EventTime
		schema = eventgrid.SchemaEventGrid
		d.Source = d.Topic
	} else {
		return nil, ErrTriggerPayloadMalformed
	}

	return &EventGrid{
		ID:       d.ID,
		Topic:    d.Topic,
		Source:   d.Source,
		Subject:  d.Subject,
		Type:     eventType,
		Time:     eventTime,
		Data:     d.Data,
		Metadata: t.Metadata,
		Schema:   schema,
	}, nil
}

// eventGridTrigger is the incoming request from the Function host.
type eventGridTrigger struct {
	Data     map[string]event
	Metadata EventGridMetadata
}

// event is the incoming event from the Function host. It contains
// all the properties that are included in both the cloud events
// schema and the event grid schema.
type event struct {
	ID              string    `json:"id"`
	Topic           string    `json:"topic"`
	Source          string    `json:"source"`
	Subject         string    `json:"subject"`
	Type            string    `json:"type"`
	EventType       string    `json:"eventType"`
	Time            time.Time `json:"time"`
	EventTime       time.Time `json:"eventTime"`
	SpecVersion     string    `json:"specversion"`
	DataVersion     string    `json:"dataVersion"`
	MetadataVersion string    `json:"metadataVersion"`
	Data            data.Raw  `json:"data"`
}
