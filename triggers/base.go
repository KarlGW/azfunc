package triggers

import (
	"encoding/json"
	"net/http"

	"github.com/KarlGW/azfunc/data"
)

// Base represents a base Function App trigger. With custom handlers many
// triggers that are not HTTP triggers share the same data structure.
type Base struct {
	data.Raw
	Metadata map[string]any
}

// Parse the Raw data of the Base trigger into the provided
// value.
func (t Base) Parse(v any) error {
	return json.Unmarshal(t.Raw, &v)
}

// Data returns the Raw data of the Base trigger.
func (t Base) Data() data.Raw {
	return t.Raw
}

// NewBase creates an returns a Base trigger from the provided
// *http.Request.
func NewBase(r *http.Request, name string, options ...Option) (Base, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	metadata, raw, err := triggerData(r, name)
	if err != nil {
		return Base{}, err
	}

	return Base{
		Raw:      raw,
		Metadata: metadata,
	}, nil
}

// baseTrigger is the incoming request from the Function host.
type baseTrigger struct {
	Data     map[string]data.Raw
	Metadata map[string]any
}

// triggerData handles the incoming request and returns a trigger
// metadata, raw data and an error (if any).
func triggerData(r *http.Request, name string) (map[string]any, data.Raw, error) {
	var t baseTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, nil, ErrTriggerPayloadMalformed
	}
	defer r.Body.Close()

	d, ok := t.Data[name]
	if !ok {
		return t.Metadata, nil, ErrTriggerNameIncorrect
	}

	return t.Metadata, d, nil
}
