package triggers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/KarlGW/azfunc/data"
)

// QueueStorage represents a Queue Storage trigger.
type QueueStorage struct {
	data     data.Raw
	Metadata QueueStorageMetadata
}

// QueueStorageMetadata represents the metadata for a Queue Storage trigger.
type QueueStorageMetadata struct {
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
func (t QueueStorage) Parse(v any) error {
	return json.Unmarshal(t.data, &v)
}

// Data returns the data of the Queue Storage trigger.
func (t QueueStorage) Data() data.Raw {
	return t.data
}

// NewQueueStorage creates and returns a new QueueStorage trigger from the
// provided *http.Request.
func NewQueueStorage(r *http.Request, name string, options ...Option) (*QueueStorage, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	var t queueStorageTrigger
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

	return &QueueStorage{
		data:     d,
		Metadata: t.Metadata,
	}, nil
}

// queueStorageTrigger is the incoming request from the function host.
type queueStorageTrigger struct {
	Data     map[string]data.Raw
	Metadata QueueStorageMetadata
}
