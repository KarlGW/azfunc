package triggers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/KarlGW/azfunc/data"
)

// Queue represents a Queue Storage trigger.
type Queue struct {
	Metadata QueueMetadata
	data     data.Raw
}

// QueueMetadata represents the metadata for a Queue Storage trigger.
type QueueMetadata struct {
	DequeueCount    string
	ID              string
	PopReceipt      string
	ExpirationTime  time.Time
	InsertionTime   time.Time
	NextVisibleTime time.Time
	Metadata
}

// Parse the data of the Queue Storage trigger into the provided
// value.
func (t Queue) Parse(v any) error {
	return json.Unmarshal(t.data, &v)
}

// Data returns the data of the Queue Storage trigger.
func (t Queue) Data() data.Raw {
	return t.data
}

// NewQueue creates and returns a new Queue Storage trigger from the
// provided *http.Request.
func NewQueue(r *http.Request, name string, options ...Option) (*Queue, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	var t queueTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, ErrTriggerPayloadMalformed
	}
	defer r.Body.Close()

	d, ok := t.Data[name]
	if !ok {
		return nil, ErrTriggerNameIncorrect
	}

	t.Metadata.ID = strings.Trim(t.Metadata.ID, "\"")
	t.Metadata.PopReceipt = strings.Trim(t.Metadata.PopReceipt, "\"")

	return &Queue{
		data:     d,
		Metadata: t.Metadata,
	}, nil
}

// queueTrigger is the incoming request from the function host.
type queueTrigger struct {
	Data     map[string]data.Raw
	Metadata QueueMetadata
}
