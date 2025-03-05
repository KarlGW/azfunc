package trigger

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/potatoattack/azfunc/data"
)

// ServiceBus represents a Service Bus trigger. It supports both
// queues and topics with subscriptions.
type ServiceBus struct {
	Data     data.Raw
	Metadata ServiceBusMetadata
}

// ServiceBusOptions contains options for a Service Bus trigger.
type ServiceBusOptions struct{}

// ServiceBusOption is a function that sets options on a Service Bus
// trigger.
type ServiceBusOption func(o *ServiceBusOptions)

// ServiceBusMetadata represents the metadata for a Service Bus trigger.
type ServiceBusMetadata struct {
	MessageReceiver       map[string]any
	MessageSession        map[string]any
	MessageActions        map[string]any
	SessionActions        map[string]any
	ReceiveActions        map[string]any
	ApplicationProperties map[string]any
	UserProperties        map[string]any
	DeliveryCount         string
	LockToken             string
	MessageID             string
	ContentType           string
	SequenceNumber        string
	Metadata
	ExpiresAtUTC    TimeISO8601 `json:"ExpiresAtUtc"`
	ExpiresAt       TimeISO8601
	EnqueuedTimeUTC TimeISO8601 `json:"EnqueuedTimeUtc"`
	EnqueuedTime    TimeISO8601
	Client          ServiceBusMetadataClient
}

// ServiceBusMetadataClient represents client of the service bus trigger
// metadata.
type ServiceBusMetadataClient struct {
	FullyQualifiedNamespace string
	Identifier              string
	TransportType           int
	IsClosed                bool
}

// Parse the data for the Service Bus trigger into the provided
// value.
func (t ServiceBus) Parse(v any) error {
	return json.Unmarshal(t.Data, &v)
}

// NewServiceBus creates and returns a new Service Bus trigger from the
// provided *http.Request.
func NewServiceBus(r *http.Request, name string, options ...ServiceBusOption) (*ServiceBus, error) {
	opts := ServiceBusOptions{}
	for _, option := range options {
		option(&opts)
	}

	var t serviceBusTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, ErrTriggerPayloadMalformed
	}
	defer r.Body.Close()

	d, ok := t.Data[name]
	if !ok {
		return nil, ErrTriggerNameIncorrect
	}

	t.Metadata.LockToken = strings.Trim(t.Metadata.LockToken, "\"")
	t.Metadata.MessageID = strings.Trim(t.Metadata.MessageID, "\"")
	t.Metadata.ContentType = strings.Trim(t.Metadata.ContentType, "\"")

	return &ServiceBus{
		Data:     d,
		Metadata: t.Metadata,
	}, nil
}

// serviceBusTrigger is the incoming request from the function host.
type serviceBusTrigger struct {
	Data     map[string]data.Raw
	Metadata ServiceBusMetadata
}
