package trigger

import (
	"encoding/json"
	"net/http"

	"github.com/potatoattack/azfunc/data"
)

// Generic represents a generic Function App trigger. With custom handlers many
// triggers that are not HTTP triggers share the same data structure.
type Generic struct {
	Metadata map[string]any
	Data     data.Raw
}

// GenericOptions contains options for a Generic trigger.
type GenericOptions struct{}

// GenericOption is a function that sets options on a Generic trigger.
type GenericOption func(o *GenericOptions)

// Parse the data of the Generic trigger into the provided
// value.
func (t Generic) Parse(v any) error {
	return json.Unmarshal(t.Data, &v)
}

// NewGeneric creates an returns a Generic trigger from the provided
// *http.Request.
func NewGeneric(r *http.Request, name string, options ...GenericOption) (*Generic, error) {
	opts := GenericOptions{}
	for _, option := range options {
		option(&opts)
	}

	var t genericTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, ErrTriggerPayloadMalformed
	}
	defer r.Body.Close()

	d, ok := t.Data[name]
	if !ok {
		return nil, ErrTriggerNameIncorrect
	}

	return &Generic{
		Data:     d,
		Metadata: t.Metadata,
	}, nil
}

// genericTrigger is the incoming request from the Function host.
type genericTrigger struct {
	Data     map[string]data.Raw
	Metadata map[string]any
}
