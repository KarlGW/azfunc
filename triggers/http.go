package triggers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KarlGW/azfunc/data"
)

// HTTP represents an HTTP trigger.
type HTTP struct {
	URL      string
	Method   string
	Body     data.Raw
	Headers  http.Header
	Params   map[string]string
	Query    map[string]string
	Metadata map[string]any
}

// Parse the body from the HTTP trigger into the provided value.
func (t HTTP) Parse(v any) error {
	return json.Unmarshal(t.Body, &v)
}

// Data returns the Raw data of the HTTP trigger.
func (t HTTP) Data() data.Raw {
	return t.Body
}

// NewHTTP creates and returns an HTTP trigger from the provided
// *http.Request.
func NewHTTP(r *http.Request, options ...Option) (HTTP, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	var t httpTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return HTTP{}, fmt.Errorf("%w: %w", ErrTriggerPayloadMalformed, err)
	}
	defer r.Body.Close()

	return HTTP{
		URL:      t.Data.Req.URL,
		Method:   t.Data.Req.Method,
		Body:     t.Data.Req.Body,
		Headers:  t.Data.Req.Headers,
		Params:   t.Data.Req.Params,
		Query:    t.Data.Req.Query,
		Metadata: t.Metadata,
	}, nil
}

// httpTrigger is the incoming request from the Function host.
type httpTrigger struct {
	Data struct {
		Req struct {
			URL     string `json:"Url"`
			Method  string
			Body    data.Raw
			Headers http.Header
			Params  map[string]string
			Query   map[string]string
		} `json:"req"`
	}
	Metadata map[string]any
}
