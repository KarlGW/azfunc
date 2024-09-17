package eventgrid

import (
	"time"

	"github.com/KarlGW/azfunc/uuid"
)

// Schema represents the schema of the event.
type Schema int

const (
	// CloudEvents is the CloudEvents schema.
	SchemaCloudEvents Schema = iota
	// EventGrid is the Event Grid schema.
	SchemaEventGrid Schema = iota
)

func (s Schema) String() string {
	switch s {
	case SchemaCloudEvents:
		return "CloudEvents"
	case SchemaEventGrid:
		return "EventGrid"
	}
	return ""
}

// EventProvider is an interface that represents an event provider.
type EventProvider interface {
	Schema() Schema
	JSON() []byte
}

// newUUID is a function that generates a new UUID.
var newUUID = func() (string, error) {
	return uuid.New()
}

// nowUTC is a function that returns the current time in UTC.
var nowUTC = func() time.Time {
	return time.Now().UTC()
}
