package azfunc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

// trigger is the interface that wraps around method Data.
type Triggerable interface {
	Data() []byte
}

// Trigger represents an incoming request (trigger) from the
// Azure Function Host.
type Trigger[T Triggerable] struct {
	Payload  map[string]T `json:"Data"`
	Metadata map[string]any
	d        []byte
	n        string
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

// Trigger returns underlying trigger.
func (t Trigger[T]) Trigger() T {
	return t.Payload[t.n]
}

// HTTPTrigger represnts a Function App HTTP Trigger and contains
// the incoming HTTP data.
type HTTPTrigger struct {
	URL     string `json:"Url"`
	Method  string
	Body    RawMessage
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
	RawMessage
}

// Data returns data of the trigger as a []byte.
func (t GenericTrigger) Data() []byte {
	return t.RawMessage
}

// TriggerOptions contains options for functions and methods related
// to triggers.
type TriggerOptions struct {
	Name string
}

// TriggerOption[T] is function that sets options on TriggerOptions.
type TriggerOption func(*TriggerOptions)

// WithName sets the trigger name get data from. The name should
// match the incoming trigger (binding) name in function.json.
func WithName(name string) TriggerOption {
	return func(o *TriggerOptions) {
		o.Name = name
	}
}

// NewTrigger handles a request from the Function host and returns a Trigger[T].
func NewTrigger[T Triggerable](r *http.Request, options ...TriggerOption) (Trigger[T], error) {
	opts := TriggerOptions{}
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

// NewRequest takes the request from the Function Host and creates
// a new *http.Request from it. Suitable in scenarios like a middleware
// to extract data from an HTTP trigger request (such as headers etc),
// or pass it on to the next handler as an ordinarily formatted
// *http.Request.
func NewRequest(r *http.Request) (*http.Request, error) {
	trigger, err := NewTrigger[HTTPTrigger](r)
	if err != nil {
		return nil, err
	}
	request, ok := trigger.Payload["req"]
	if !ok {
		return nil, ErrNotHTTPTrigger
	}

	u, err := buildURL(request.URL, request.Params, request.Query)
	if err != nil {
		return nil, err
	}

	var body *bytes.Buffer
	if request.Body != nil {
		body = bytes.NewBuffer(request.Body)
	}

	req, err := http.NewRequest(request.Method, u, body)
	if err != nil {
		return nil, err
	}
	req.Header = request.Headers

	return req, nil
}

// Parse the incoming Function host request (trigger) and set
// the data to the provided value.
func Parse[T Triggerable](r *http.Request, v any, options ...TriggerOption) error {
	trigger, err := NewTrigger[T](r, options...)
	if err != nil {
		return err
	}
	return trigger.Parse(v)
}

// Data returns the data from the incoming Function host
// request (trigger).
func Data[T Triggerable](r *http.Request, options ...TriggerOption) ([]byte, error) {
	trigger, err := NewTrigger[T](r, options...)
	if err != nil {
		return nil, err
	}
	return trigger.Data(), nil
}

// buildURL from the provided url, parameters and query.
func buildURL(u string, p, q map[string]string) (string, error) {
	_url, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	for k, v := range p {
		_url.Path, err = url.JoinPath(_url.Path, k+"/"+v)
		if err != nil {
			return "", err
		}
	}

	if q != nil {
		query := _url.Query()
		for k, v := range q {
			query.Add(k, v)
		}
		_url.RawQuery = query.Encode()
	}

	return _url.String(), nil
}

// Trigger aliases.

// QueueTrigger represnts a Function App Queue Trigger and contains
// the incoming queue message data.
type QueueTrigger = GenericTrigger
