package triggers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KarlGW/azfunc/data"
)

// EventGridSchema is the type of schema that the Event Grid trigger
// is using.
type EventGridSchema int

const (
	// EventGridSchemaCloudEvents is the CloudEvents schema.
	EventGridSchemaCloudEvents EventGridSchema = iota
	// EventGridSchemaEventGrid is the Event Grid schema.
	EventGridSchemaEventGrid
)

// String returns the string representation of the EventGridSchema.
func (s EventGridSchema) String() string {
	switch s {
	case EventGridSchemaCloudEvents:
		return "CloudEvents"
	case EventGridSchemaEventGrid:
		return "EventGrid"
	}
	return ""
}

// EventGrid represents an Event Grid trigger. It handles both
// Cloud Events and Event Grid events. For both types of events,
// the 'Topic' is used to describe the source of the event. This
// to match the nomenclature of the Event Grid service.
type EventGrid struct {
	Time     time.Time
	Metadata EventGridMetadata
	ID       string
	Topic    string
	Subject  string
	Type     string
	Data     data.Raw
	Schema   EventGridSchema
}

// EventGridOptions contains options for an Event Grid trigger.
type EventGridOptions struct{}

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

	var topic, typ string
	var eventTime time.Time
	var schema EventGridSchema
	if len(d.SpecVersion) > 0 {
		topic, typ = d.Source, d.Type
		eventTime = d.Time
		schema = EventGridSchemaCloudEvents
	} else if len(d.EventType) > 0 {
		topic, typ = d.Topic, d.EventType
		eventTime = d.EventTime
		schema = EventGridSchemaEventGrid
	} else {
		return nil, ErrTriggerPayloadMalformed
	}

	return &EventGrid{
		ID:       d.ID,
		Topic:    topic,
		Subject:  d.Subject,
		Type:     typ,
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
