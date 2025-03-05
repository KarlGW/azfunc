package output

import (
	"net/http"
	"testing"

	"github.com/potatoattack/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewHTTP(t *testing.T) {
	var tests = []struct {
		name  string
		input []HTTPOption
		want  *HTTP
	}{
		{
			name:  "defaults",
			input: nil,
			want: &HTTP{
				name:       "res",
				statusCode: http.StatusOK,
				body:       nil,
				header:     http.Header{},
			},
		},
		{
			name: "with options",
			input: []HTTPOption{
				func(o *HTTPOptions) {
					o.Name = "httpoutput"
					o.StatusCode = http.StatusNotFound
					o.Body = data.Raw(`{"message":"not found"}`)
					o.Header = http.Header{"Content-Type": {"application/json"}}
				},
			},
			want: &HTTP{
				name:       "httpoutput",
				statusCode: http.StatusNotFound,
				body:       data.Raw(`{"message":"not found"}`),
				header:     http.Header{"Content-Type": {"application/json"}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewHTTP(test.input...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(HTTP{})); diff != "" {
				t.Errorf("NewHTTP() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestHTTP_Write(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		got := &HTTP{}
		got.Write([]byte(`{"message":"hello"}`))
		want := &HTTP{body: data.Raw(`{"message":"hello"}`)}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(HTTP{})); diff != "" {
			t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestHTTP_WriteHeader(t *testing.T) {
	t.Run("WriteHeader", func(t *testing.T) {
		got := &HTTP{}
		got.WriteHeader(http.StatusNotFound)
		want := &HTTP{statusCode: http.StatusNotFound}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(HTTP{})); diff != "" {
			t.Errorf("WriteHeader() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestHTTP_WriteResponse(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			statusCode int
			body       []byte
			options    []HTTPOption
		}
		want *HTTP
	}{
		{
			name: "write response",
			input: struct {
				statusCode int
				body       []byte
				options    []HTTPOption
			}{
				statusCode: http.StatusCreated,
				body:       []byte(`{"message":"hello"}`),
			},
			want: &HTTP{
				statusCode: http.StatusCreated,
				body:       data.Raw(`{"message":"hello"}`),
			},
		},
		{
			name: "write response with hheaders",
			input: struct {
				statusCode int
				body       []byte
				options    []HTTPOption
			}{
				statusCode: http.StatusCreated,
				body:       []byte(`{"message":"hello"}`),
				options: []HTTPOption{
					WithHeader(http.Header{"Content-Type": {"application/json"}}),
				},
			},
			want: &HTTP{
				statusCode: http.StatusCreated,
				body:       data.Raw(`{"message":"hello"}`),
				header:     http.Header{"Content-Type": {"application/json"}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := &HTTP{}
			got.WriteResponse(test.input.statusCode, test.input.body, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(HTTP{})); diff != "" {
				t.Errorf("WriteResponse() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestHTTP_Name(t *testing.T) {
	var tests = []struct {
		name  string
		input *HTTP
		want  string
	}{
		{
			name:  "default",
			input: &HTTP{},
			want:  "res",
		},
		{
			name:  "with name",
			input: &HTTP{name: "httpoutput"},
			want:  "httpoutput",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Name()

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(HTTP{})); diff != "" {
				t.Errorf("Name() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
