package triggers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/KarlGW/azfunc/data"
)

var (
	// ErrHTTPInvalidContentType is returned when an invalid Content-Type provided.
	ErrHTTPInvalidContentType = errors.New("invalid Content-Type")
	// ErrHTTPInvalidBody is returned when the HTTP body is invalid.
	ErrHTTPInvalidBody = errors.New("invalid body")
)

var (
	// defaultMultipartFormMaxMemory is the default memory to use
	// when parsing multipart form data.
	defaultMultipartFormMaxMemory int64 = 32 << 20
)

// HTTP represents an HTTP trigger.
type HTTP struct {
	Headers    http.Header
	Params     map[string]string
	Query      map[string]string
	URL        string `json:"Url"`
	Method     string
	Body       data.Raw
	Metadata   HTTPMetadata
	Identities []HTTPIdentity
}

// HTTPOptions contains options for an HTTP trigger.
type HTTPOptions struct {
	Name string
}

// HTTPOption is a function that sets options on an HTTP trigger.
type HTTPOption func(o *HTTPOptions)

// HTTPMetadata represents the metadata for an HTTP trigger.
type HTTPMetadata struct {
	Headers map[string]string
	Params  map[string]string
	Query   map[string]string
	Metadata
}

// HTTPIdentity represent a part of the Identities field
// of the incoming trigger request.
type HTTPIdentity struct {
	Actor              any
	BootstrapContext   any
	Label              any
	Name               any
	AuthenticationType string
	NameClaimType      string
	RoleClaimType      string
	Claims             []HTTPIdentityClaims
	IsAuthenticated    bool
}

// HTTPIdentityClaims represent the claims of an HTTPIdentity.
type HTTPIdentityClaims struct {
	Properties     map[string]string
	Issuer         string
	OriginalIssuer string
	Type           string
	Value          string
	ValueType      string
}

// Parse the body from the HTTP trigger into the provided value.
func (t HTTP) Parse(v any) error {
	return json.Unmarshal(t.Body, &v)
}

// Form parses the HTTP trigger for form data sent with Content-Type
// application/x-www-form-urlencoded and returns it as url.Values.
func (t HTTP) Form() (url.Values, error) {
	contentType := t.Headers.Get("Content-Type")
	if strings.ToLower(contentType) != "application/x-www-form-urlencoded" {
		return nil, fmt.Errorf("%w: %s", ErrHTTPInvalidContentType, contentType)
	}

	data, err := url.ParseQuery(string(t.Body))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrHTTPInvalidBody, string(t.Body))
	}
	if len(data) == 1 {
		for _, v := range data {
			if len(v[0]) == 0 {
				return nil, fmt.Errorf("%w: %s", ErrHTTPInvalidBody, string(t.Body))
			}
		}
	}

	return data, nil
}

// MultipartForm parses the HTTP trigger for multipart form data and returns the
// resulting *multipart.Form. The whole request body is parsed and up to a total
// of maxMemory bytes of its file parts are stored in memory, with the remainder
// stored on disk in temporary files. If 0 or less is provided it will
// default to 32 MB.
func (t HTTP) MultipartForm(maxMemory int64) (*multipart.Form, error) {
	r, err := http.NewRequest(t.Method, t.URL, bytes.NewReader(t.Body))
	if err != nil {
		return nil, err
	}
	hdr := t.Headers.Get("Content-Type")
	if len(hdr) == 0 {
		return nil, ErrHTTPInvalidContentType
	}
	r.Header.Add("Content-Type", hdr)

	if maxMemory <= 0 {
		maxMemory = defaultMultipartFormMaxMemory
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrHTTPInvalidBody, err.Error())
	}
	return r.MultipartForm, nil
}

// NewHTTP creates and returns an HTTP trigger from the provided
// *http.Request. By default it will use the name "req" for the
// trigger. This can be overridden with providing a name
// in the options.
func NewHTTP(r *http.Request, options ...HTTPOption) (*HTTP, error) {
	opts := HTTPOptions{
		Name: "req",
	}
	for _, option := range options {
		option(&opts)
	}

	var t httpTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, ErrTriggerPayloadMalformed
	}
	defer r.Body.Close()

	d, ok := t.Data[opts.Name]
	if !ok {
		return nil, ErrTriggerNameIncorrect
	}
	d.Metadata = t.Metadata

	return &d, nil
}

// httpTrigger is the incoming request from the Function host.
type httpTrigger struct {
	Data     map[string]HTTP
	Metadata HTTPMetadata
}
