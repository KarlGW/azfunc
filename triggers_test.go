package azfunc

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewTrigger_HTTPTrigger(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    Trigger[HTTPTrigger]
		wantErr error
	}{
		{
			name: "new Trigger[HTTPTrigger]",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
			},
			want: Trigger[HTTPTrigger]{
				Payload: map[string]HTTPTrigger{
					"req": {
						URL:    "http://localhost:7071/api/endpoint",
						Method: http.MethodPost,
						Body:   []byte(`{"message":"hello","number":2}`),
						Headers: map[string][]string{
							"Content-Type": {"application/json"},
						},
						Params: map[string]string{},
						Query:  map[string]string{},
					},
				},
				Metadata: map[string]any{},
				d:        []byte(`{"message":"hello","number":2}`),
				n:        "req",
			},
			wantErr: nil,
		},
		{
			name: "new Trigger[HTTPTrigger] with simple body, parameters and query",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(httpTrigger2)),
			},
			want: Trigger[HTTPTrigger]{
				Payload: map[string]HTTPTrigger{
					"req": {
						URL:    "http://localhost:7071/api/endpoint",
						Method: http.MethodPost,
						Body:   []byte(`hello`),
						Headers: map[string][]string{
							"Content-Type": {"application/json"},
						},
						Params: map[string]string{
							"resource": "1",
						},
						Query: map[string]string{
							"order": "desc",
						},
					},
				},
				Metadata: map[string]any{},
				d:        []byte(`hello`),
				n:        "req",
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewTrigger[HTTPTrigger](test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Trigger[HTTPTrigger]{})); diff != "" {
				t.Errorf("NewTrigger[HTTPTrigger]() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("NewTrigger[HTTPTrigger]() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestHTTPTrigger_Trigger(t *testing.T) {
	t.Run("Get underlying HTTPTrigger", func(t *testing.T) {
		trigger, err := NewTrigger[HTTPTrigger](&http.Request{
			Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
		})
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		want := HTTPTrigger{
			URL:    "http://localhost:7071/api/endpoint",
			Method: http.MethodPost,
			Query:  map[string]string{},
			Params: map[string]string{},
			Headers: http.Header{
				"Content-Type": {"application/json"},
			},
			Body: []byte(`{"message":"hello","number":2}`),
		}

		got := trigger.Trigger()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Trigger() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestGenericTrigger_Trigger(t *testing.T) {
	t.Run("Get underlying QueueTrigger", func(t *testing.T) {
		trigger, err := NewTrigger[GenericTrigger](&http.Request{
			Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
		}, WithName("queue"))
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		want := GenericTrigger{
			RawMessage: []byte(`{"message":"hello","number":2}`),
		}

		got := trigger.Trigger()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Trigger() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestNewTrigger_GenericTrigger(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			options []TriggerOption
		}
		want    Trigger[GenericTrigger]
		wantErr error
	}{
		{
			name: "new Trigger[GenericTRigger]",
			input: struct {
				req     *http.Request
				options []TriggerOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
				},
				options: []TriggerOption{
					WithName("queue"),
				},
			},
			want: Trigger[GenericTrigger]{
				Payload: map[string]GenericTrigger{
					"queue": {
						RawMessage: []byte(`{"message":"hello","number":2}`),
					},
				},
				Metadata: map[string]any{},
				d:        []byte(`{"message":"hello","number":2}`),
				n:        "queue",
			},
			wantErr: nil,
		},
		{
			name: "new Trigger[GenericTRigger] - simple body",
			input: struct {
				req     *http.Request
				options []TriggerOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger2)),
				},
				options: []TriggerOption{
					WithName("queue"),
				},
			},
			want: Trigger[GenericTrigger]{
				Payload: map[string]GenericTrigger{
					"queue": {
						RawMessage: []byte(`hello`),
					},
				},
				Metadata: map[string]any{},
				d:        []byte(`hello`),
				n:        "queue",
			},
			wantErr: nil,
		},
		{
			name: "new Trigger[GenericTRigger] - error no name provided",
			input: struct {
				req     *http.Request
				options []TriggerOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger2)),
				},
				options: nil,
			},
			want:    Trigger[GenericTrigger]{},
			wantErr: ErrTriggerNameIncorrect,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewTrigger[GenericTrigger](test.input.req, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Trigger[GenericTrigger]{})); diff != "" {
				t.Errorf("NewTrigger[GenericTrigger]() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("NewTrigger[GenericTrigger]() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestTrigger_Parse(t *testing.T) {
	var tests = []struct {
		name    string
		input   func() Trigger[HTTPTrigger]
		want    testType
		wantErr error
	}{
		{
			name: "Parse the data in Trigger[HTTPTrigger]",
			input: func() Trigger[HTTPTrigger] {
				req := &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
				}
				trigger, _ := NewTrigger[HTTPTrigger](req)
				return trigger
			},
			want: testType{
				Message: "hello",
				Number:  2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got testType
			gotErr := test.input().Parse(&got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestHTTPTrigger_Data(t *testing.T) {
	var tests = []struct {
		name    string
		input   func() Trigger[HTTPTrigger]
		want    []byte
		wantErr error
	}{
		{
			name: "Parse the data in Trigger[GenericTrigger]",
			input: func() Trigger[HTTPTrigger] {
				req := &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
				}
				trigger, _ := NewTrigger[HTTPTrigger](req)
				return trigger
			},
			want: []byte(`{"message":"hello","number":2}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input().Data()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestGenericTrigger_Parse(t *testing.T) {
	var tests = []struct {
		name    string
		input   func() Trigger[GenericTrigger]
		want    testType
		wantErr error
	}{
		{
			name: "Parse the data in Trigger[GenericTrigger]",
			input: func() Trigger[GenericTrigger] {
				req := &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
				}
				trigger, _ := NewTrigger[GenericTrigger](req, WithName("queue"))
				return trigger
			},
			want: testType{
				Message: "hello",
				Number:  2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got testType
			gotErr := test.input().Parse(&got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestGenericTrigger_Data(t *testing.T) {
	var tests = []struct {
		name    string
		input   func() Trigger[GenericTrigger]
		want    []byte
		wantErr error
	}{
		{
			name: "Parse the data in Trigger[GenericTrigger]",
			input: func() Trigger[GenericTrigger] {
				req := &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
				}
				trigger, _ := NewTrigger[GenericTrigger](req, WithName("queue"))
				return trigger
			},
			want: []byte(`{"message":"hello","number":2}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input().Data()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestParse_HTTPTrigger(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    testType
		wantErr error
	}{
		{
			name: "Parse - HTTPTrigger",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
			},
			want: testType{
				Message: "hello",
				Number:  2,
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got testType
			gotErr := Parse[HTTPTrigger](test.input, &got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Parse() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestData_HTTPTrigger(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    []byte
		wantErr error
	}{
		{
			name: "Data from HTTPTrigger",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
			},
			want:    []byte(`{"message":"hello","number":2}`),
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := Data[HTTPTrigger](test.input)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Data() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Data() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestParse_GenericTrigger(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    testType
		wantErr error
	}{
		{
			name: "Parse - GenericTrigger",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
			},
			want: testType{
				Message: "hello",
				Number:  2,
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got testType
			gotErr := Parse[QueueTrigger](test.input, &got, WithName("queue"))

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Parse() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestData_GenericTrigger(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    []byte
		wantErr error
	}{
		{
			name: "Data from GenericTrigger",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
			},
			want:    []byte(`{"message":"hello","number":2}`),
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := Data[GenericTrigger](test.input, WithName("queue"))

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Data() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Data() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestNewRequestFrom(t *testing.T) {
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
				body:   []byte(`{"message":"hello","number":2}`),
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
				body:   []byte(`hello`),
				header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, gotErr := NewRequest(test.input)
			got := newWantRequest(req)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(wantRequest{})); diff != "" {
				t.Errorf("NewRequestFrom() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("NewRequestFrom() = unexpected error (-want +got)\n%s\n", diff)
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

type testType struct {
	Message string
	Number  float64
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
      "Url": "http://localhost:7071/api/endpoint",
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

var queueTrigger1 = []byte(`{
  "Data": {
    "queue": "{\"message\":\"hello\",\"number\":2}"
  },
  "Metadata": {
  }
}
`)

var queueTrigger2 = []byte(`{
  "Data": {
    "queue": "hello"
  },
  "Metadata": {
  }
}
`)
