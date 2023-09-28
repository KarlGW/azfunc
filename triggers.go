package azfunc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// trigger is the interface that wraps around method Data.
type trigger interface {
	Data() []byte
}

// Input represents an incoming request (trigger) from the
// Azure Function Host.
type Input[T trigger] struct {
	Data     map[string]T
	Metadata map[string]any
	d        []byte
	n        string
}

// Parse is used to parse the data contained in a trigger into
// the provided struct.
func (t Input[T]) Parse(v any) error {
	return json.Unmarshal(t.d, &v)
}

// HTTPTrigger represnts a Function App HTTP Trigger and contains
// the incoming HTTP data.
type HTTPTrigger struct {
	URL     string `json:"Url"`
	Method  string
	Body    []byte
	Headers http.Header
	Params  map[string]string
	Query   map[string]string
}

// Data returns the body of the HTTPTrigger as a []byte.
func (t HTTPTrigger) Data() []byte {
	return []byte(t.Body)
}

// GenericTrigger represents a generic Function App trigger. With custom handlers all
// triggers that are not HTTP triggers share the same data structure.
type GenericTrigger struct {
	json.RawMessage
}

// Data returns data of the trigger as a []byte.
func (t GenericTrigger) Data() []byte {
	return t.RawMessage
}

// TriggerOption[T] is function that sets option on a Input[T]
type TriggerOption[T trigger] func(*Input[T])

// WithName sets the trigger name get data from. The name should
// match the incoming trigger (binding) name in function.json.
func WithName[T trigger](name string) TriggerOption[T] {
	return func(t *Input[T]) {
		t.n = name
	}
}

// NewInput handles a request from the Function host and returns a Input[T].
func NewInput[T trigger](r *http.Request, options ...TriggerOption[T]) (Input[T], error) {
	t := Input[T]{
		n: "req",
	}
	for _, option := range options {
		option(&t)
	}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return Input[T]{}, err
	}

	d, ok := t.Data[t.n]
	if !ok {
		return Input[T]{}, errors.New("name incorrect")
	}
	t.d = d.Data()
	return t, nil
}

// QueueTrigger represnts a Function App Queue Trigger and contains
// the incoming queue message data.
type QueueTrigger = GenericTrigger

// RequestFrom takes the request from the Function Host and creates
// a new *http.Request from it. Suitable in scenarios like a middleware
// to extract data from an HTTP trigger request (such as headers etc),
// or pass it on to the next handler as an ordinarily formatted
// *http.Request.
func RequestFrom(r *http.Request) (*http.Request, error) {
	input, err := NewInput[HTTPTrigger](r)
	if err != nil {
		return nil, err
	}
	data, ok := input.Data["req"]
	if !ok {
		return nil, errors.New("not an HTTP trigger")
	}

	u, err := url.Parse(data.URL)
	if err != nil {
		return nil, err
	}
	for k, v := range data.Params {
		u.Path, err = url.JoinPath(u.Path, k+"/"+v)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range data.Query {
		u.Query().Add(k, v)
	}
	u.RawQuery = u.Query().Encode()

	var body *bytes.Buffer
	if data.Body != nil {
		body = bytes.NewBuffer(data.Body)
	}

	return http.NewRequest(data.Method, u.String(), body)
}
