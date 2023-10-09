package triggers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Generic represents a generic Function App trigger. With custom handlers many
// triggers that are not HTTP triggers share the same data structure.
type Generic struct {
	RawData
	Metadata map[string]any
}

// Parse the RawData of the Generic trigger into the provided
// value.
func (t Generic) Parse(v any) error {
	return json.Unmarshal(t.RawData, &v)
}

// Data returns the RawData of the Generic trigger.
func (t Generic) Data() RawData {
	return t.RawData
}

// NewGeneric creates an returns a Generic trigger from the provided
// *http.Request.
func NewGeneric(r *http.Request, name string, options ...Option) (Generic, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	metadata, rawData, err := triggerData(r, name)
	if err != nil {
		return Generic{}, err
	}

	return Generic{
		RawData:  rawData,
		Metadata: metadata,
	}, nil
}

// genericTrigger is the incoming request from the Function host.
type genericTrigger struct {
	Data     map[string]RawData
	Metadata map[string]any
}

// triggerData handles the incoming request and returns a trigger
// metadata, raw data and an error (if any).
func triggerData(r *http.Request, name string) (map[string]any, RawData, error) {
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
