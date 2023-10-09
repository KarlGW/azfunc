package triggers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KarlGW/azfunc/data"
)

// Generic represents a generic Function App trigger. With custom handlers many
// triggers that are not HTTP triggers share the same data structure.
type Generic struct {
	data.Raw
	Metadata map[string]any
}

// Parse the Raw data of the Generic trigger into the provided
// value.
func (t Generic) Parse(v any) error {
	return json.Unmarshal(t.Raw, &v)
}

// Data returns the Raw data of the Generic trigger.
func (t Generic) Data() data.Raw {
	return t.Raw
}

// NewGeneric creates an returns a Generic trigger from the provided
// *http.Request.
func NewGeneric(r *http.Request, name string, options ...Option) (Generic, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	metadata, raw, err := triggerData(r, name)
	if err != nil {
		return Generic{}, err
	}

	return Generic{
		Raw:      raw,
		Metadata: metadata,
	}, nil
}

// genericTrigger is the incoming request from the Function host.
type genericTrigger struct {
	Data     map[string]data.Raw
	Metadata map[string]any
}

// triggerData handles the incoming request and returns a trigger
// metadata, raw data and an error (if any).
func triggerData(r *http.Request, name string) (map[string]any, data.Raw, error) {
	var t genericTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrTriggerPayloadMalformed, err)
	}
	defer r.Body.Close()

	d, ok := t.Data[name]
	if !ok {
		return t.Metadata, nil, ErrTriggerNameIncorrect
	}

	return t.Metadata, d, nil
}
