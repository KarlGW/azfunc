package triggers

import (
	"bytes"
	"net/http"
	"net/url"
)

// NewRequest takes the request from the Function Host and creates
// a new *http.Request from it. Suitable in scenarios like a middleware
// to extract data from an HTTP trigger request (such as headers etc),
// or pass it on to the next handler as an ordinarily formatted
// *http.Request.
func NewRequest(r *http.Request) (*http.Request, error) {
	trigger, err := New[HTTP](r)
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
