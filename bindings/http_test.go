package bindings

import (
	"net/http"
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewHTTP(t *testing.T) {
	var tests = []struct {
		name  string
		input []Option
		want  *HTTP
	}{
		{
			name:  "defaults",
			input: nil,
			want: &HTTP{
				name:       "res",
				StatusCode: http.StatusOK,
				Body:       nil,
				header:     nil,
			},
		},
		{
			name: "with options",
			input: []Option{
				func(o *Options) {
					o.Name = "httpoutput"
					o.StatusCode = http.StatusNotFound
					o.Body = data.Raw(`{"message":"not found"}`)
					o.Header = http.Header{"Content-Type": {"application/json"}}
				},
			},
			want: &HTTP{
				name:       "httpoutput",
				StatusCode: http.StatusNotFound,
				Body:       data.Raw(`{"message":"not found"}`),
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
		want := &HTTP{Body: data.Raw(`{"message":"hello"}`)}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(HTTP{})); diff != "" {
			t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestHTTP_WriteHeader(t *testing.T) {
	t.Run("WriteHeader", func(t *testing.T) {
		got := &HTTP{}
		got.WriteHeader(http.StatusNotFound)
		want := &HTTP{StatusCode: http.StatusNotFound}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(HTTP{})); diff != "" {
			t.Errorf("WriteHeader() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
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
