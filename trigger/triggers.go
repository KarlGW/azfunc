package trigger

import (
	"errors"
	"time"
)

var (
	// ErrTriggerNameIncorrect is returned when the provided trigger
	// name does not match the payload trigger name.
	ErrTriggerNameIncorrect = errors.New("trigger name incorrect")
	// ErrTriggerPayloadMalformed is returned if there is an error
	// with the payload from the Function host.
	ErrTriggerPayloadMalformed = errors.New("trigger payload malformed")
)

// Metadata represents the metadata of the incoming trigger
// request.
type Metadata struct {
	Sys MetadataSys `json:"sys"`
}

// MetadataSys contains the sys fields of the incoming trigger
// request.
type MetadataSys struct {
	MethodName string
	UTCNow     time.Time `json:"UtcNow"`
	RandGuid   string
}
