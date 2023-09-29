package azfunc

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFromRequest(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    wantRequest
		wantErr error
	}{
		{
			name: "Request without params or query",
			input: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBuffer(httpTrigger1)),
			},
			want: wantRequest{
				method: http.MethodPost,
				url:    "http://localhost:7071/api/endpoint",
				body:   []byte(`"{\"message\":\"hello\",\"number\":2}"`),
				header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Request with params and query",
			input: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{},
				Body:   io.NopCloser(bytes.NewBuffer(httpTrigger2)),
			},
			want: wantRequest{
				method: http.MethodPost,
				url:    "http://localhost:7071/api/endpoint/resource/1?order=desc",
				body:   []byte(`"hello"`),
				header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, gotErr := RequestFrom(test.input)
			got := newWantRequest(req)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(wantRequest{})); diff != "" {
				t.Errorf("RequestFrom() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("RequestFrom() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

type wantRequest struct {
	method string
	url    string
	body   []byte
	header http.Header
}

func newWantRequest(r *http.Request) wantRequest {
	var b []byte
	if r.Body != nil {
		b, _ = io.ReadAll(r.Body)
		defer r.Body.Close()
	}
	return wantRequest{
		method: r.Method,
		url:    r.URL.String(),
		body:   b,
		header: r.Header,
	}
}

var httpTrigger1 = []byte(`{
  "Data": {
    "req": {
	  "Url": "http://localhost:7071/api/endpoint",
	  "Method": "POST",
	  "Query": {},
	  "Headers": {
	    "Content-Type": [
	      "application/json"
	    ]
	  },
	  "Params": {},
	  "Body": "{\"message\":\"hello\",\"number\":2}"
	}
  },
  "Metadata": {
  }
}
`)

var httpTrigger2 = []byte(`{
  "Data": {
    "req": {
      "Url": "http://localhost:7071/api/endpoint?order=desc",
      "Method": "POST",
      "Query": {
        "order": "desc"
      },
      "Headers": {
        "Content-Type": [
          "application/json"
        ]
      },
      "Params": {
        "resource": "1"
      },
      "Body": "hello"
    }
  },
  "Metadata": {
  }
}
`)
