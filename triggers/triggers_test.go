package triggers

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew_HTTP(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    trigger[HTTP]
		wantErr error
	}{
		{
			name: "new trigger[HTTP]",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
			},
			want: trigger[HTTP]{
				Payload: map[string]HTTP{
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
			name: "new trigger[HTTP] with simple body, parameters and query",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(httpTrigger2)),
			},
			want: trigger[HTTP]{
				Payload: map[string]HTTP{
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
			got, gotErr := New[HTTP](test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(trigger[HTTP]{})); diff != "" {
				t.Errorf("New[HTTP]() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("New[HTTP]() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestNew_Generic(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			options []Option
		}
		want    trigger[Generic]
		wantErr error
	}{
		{
			name: "new trigger[Generic]",
			input: struct {
				req     *http.Request
				options []Option
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
				},
				options: []Option{
					WithName("queue"),
				},
			},
			want: trigger[Generic]{
				Payload: map[string]Generic{
					"queue": {
						RawData: []byte(`{"message":"hello","number":2}`),
					},
				},
				Metadata: map[string]any{},
				d:        []byte(`{"message":"hello","number":2}`),
				n:        "queue",
			},
			wantErr: nil,
		},
		{
			name: "new trigger[GenericTRigger] - simple body",
			input: struct {
				req     *http.Request
				options []Option
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger2)),
				},
				options: []Option{
					WithName("queue"),
				},
			},
			want: trigger[Generic]{
				Payload: map[string]Generic{
					"queue": {
						RawData: []byte(`hello`),
					},
				},
				Metadata: map[string]any{},
				d:        []byte(`hello`),
				n:        "queue",
			},
			wantErr: nil,
		},
		{
			name: "new trigger[Generic] - error no name provided",
			input: struct {
				req     *http.Request
				options []Option
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger2)),
				},
				options: nil,
			},
			want:    trigger[Generic]{},
			wantErr: ErrTriggerNameIncorrect,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := New[Generic](test.input.req, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(trigger[Generic]{})); diff != "" {
				t.Errorf("New[Generic]() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("New[Generic]() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestTrigger_Parse(t *testing.T) {
	var tests = []struct {
		name    string
		input   func() trigger[HTTP]
		want    testType
		wantErr error
	}{
		{
			name: "Parse the data in trigger[HTTP]",
			input: func() trigger[HTTP] {
				req := &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
				}
				trigger, _ := New[HTTP](req)
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

func TestHTTP_Data(t *testing.T) {
	var tests = []struct {
		name    string
		input   func() trigger[HTTP]
		want    []byte
		wantErr error
	}{
		{
			name: "Parse the data in trigger[HTTP]",
			input: func() trigger[HTTP] {
				req := &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
				}
				trigger, _ := New[HTTP](req)
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

func TestGeneric_Parse(t *testing.T) {
	var tests = []struct {
		name    string
		input   func() trigger[Generic]
		want    testType
		wantErr error
	}{
		{
			name: "Parse the data in trigger[GenericTrigger]",
			input: func() trigger[Generic] {
				req := &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
				}
				trigger, _ := New[Generic](req, WithName("queue"))
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

func TestGeneric_Data(t *testing.T) {
	var tests = []struct {
		name    string
		input   func() trigger[Generic]
		want    []byte
		wantErr error
	}{
		{
			name: "Parse the data in trigger[GenericTrigger]",
			input: func() trigger[Generic] {
				req := &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
				}
				trigger, _ := New[Generic](req, WithName("queue"))
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

func TestParse_HTTP(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    testType
		wantErr error
	}{
		{
			name: "Parse - HTTP",
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
			gotErr := Parse[HTTP](test.input, &got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Parse() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestData_HTTP(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    []byte
		wantErr error
	}{
		{
			name: "Data from HTTP",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(httpTrigger1)),
			},
			want:    []byte(`{"message":"hello","number":2}`),
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := Data[HTTP](test.input)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Data() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Data() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestParse_Generic(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    testType
		wantErr error
	}{
		{
			name: "Parse - Generic",
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
			gotErr := Parse[Generic](test.input, &got, WithName("queue"))

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Parse() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestData_Generic(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    []byte
		wantErr error
	}{
		{
			name: "Data from Generic",
			input: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(queueTrigger1)),
			},
			want:    []byte(`{"message":"hello","number":2}`),
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := Data[Generic](test.input, WithName("queue"))

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Data() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Data() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
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