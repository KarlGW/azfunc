package triggers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

// Triggerable is the interface that wraps around method Data.
type Triggerable interface {
	Data() data.Raw
	Parse(v any) error
}

// Trigger represents an incoming request (trigger) from the
// Azure Function Host.
type Trigger[T Triggerable] struct {
	Payload  map[string]T `json:"Data"`
	Metadata map[string]any
	d        []byte
	n        string
}

// New handles a request from the Function host and returns a Trigger[T].
func New[T Triggerable](r *http.Request, options ...Option) (Trigger[T], error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	if len(opts.Name) == 0 {
		opts.Name = "req"
	}

	t := Trigger[T]{
		n: opts.Name,
	}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return Trigger[T]{}, fmt.Errorf("%w: %w", ErrTriggerPayloadMalformed, err)
	}

	d, ok := t.Payload[t.n]
	if !ok {
		return Trigger[T]{}, ErrTriggerNameIncorrect
	}
	t.d = d.Data()
	return t, nil
}

// Parse is used to parse the data contained in a trigger into
// the provided struct.
func (t Trigger[T]) Parse(v any) error {
	return json.Unmarshal(t.d, &v)
}

// Data returns the data contained in the trigger.
func (t Trigger[T]) Data() []byte {
	return t.d
}

// Parse the incoming Function host request (trigger) and set
// the data to the provided value.
func Parse[T Triggerable](r *http.Request, v any, options ...Option) error {
	trigger, err := New[T](r, options...)
	if err != nil {
		return err
	}
	return trigger.Parse(v)
}

// Data returns the data from the incoming Function host
// request (trigger).
func Data[T Triggerable](r *http.Request, options ...Option) ([]byte, error) {
	trigger, err := New[T](r, options...)
	if err != nil {
		return nil, err
	}
	return trigger.Data(), nil
}

// Trigger aliases.

// Queue represents a Function App Queue Trigger and contains
// the incoming queue message data.
type Queue = Generic

// NewQueue creates an returns a Generic trigger from the provided
// *http.Request.
var NewQueue = NewGeneric
