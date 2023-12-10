package triggers

import (
	"errors"
	"time"

	"github.com/KarlGW/azfunc/data"
)

var (
	// ErrNotHTTPTrigger is returned when the provide trigger is not
	// an HTTP trigger.
	ErrNotHTTPTrigger = errors.New("not an HTTP trigger")
	// ErrTriggerNameIncorrect is returned when the provided trigger
	// name does not match the payload trigger name.
	ErrTriggerNameIncorrect = errors.New("trigger name incorrect")
	// ErrTriggerPayloadMalformed is returned if there is an error
	// with the payload from the Function host.
	ErrTriggerPayloadMalformed = errors.New("trigger payload malformed")
)

// Triggerable is the interface that wraps around methods Data and Write.
type Triggerable interface {
	// Data returns the raw data of the trigger.
	Data() data.Raw
	// Parse the raw data of the trigger into the provided value.
	Parse(v any) error
}

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
