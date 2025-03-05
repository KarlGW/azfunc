package eventgrid

import (
	"time"

	"github.com/potatoattack/azfunc/uuid"
)

// Schema represents the schema of the event.
type Schema string

const (
	// CloudEvents is the CloudEvents schema.
	SchemaCloudEvents Schema = "CloudEvents"
	// EventGrid is the Event Grid schema.
	SchemaEventGrid Schema = "EventGrid"
)

// newUUID is a function that generates a new UUID.
var newUUID = func() (string, error) {
	return uuid.New()
}

// nowUTC is a function that returns the current time in UTC.
var nowUTC = func() time.Time {
	return time.Now().UTC()
}
