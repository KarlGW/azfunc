package triggers

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewBase(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			name    string
			options []Option
		}
		want    *Base
		wantErr error
	}{
		{
			name: "NewBase",
			input: struct {
				req     *http.Request
				name    string
				options []Option
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(baseRequest1)),
				},
				name: "base",
			},
			want: &Base{
				data: data.Raw(`{"message":"hello","number":2}`),
				Metadata: map[string]any{
					"sys": map[string]any{
						"MethodName": "helloBase",
						"UtcNow":     "2023-10-12T20:13:49.640002Z",
						"RandGuid":   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewBase(test.input.req, test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Base{})); diff != "" {
				t.Errorf("NewBase() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("NewBase() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

var baseRequest1 = []byte(`{
	"Data": {
		"base": "{\"message\":\"hello\",\"number\":2}"
	  },
	  "Metadata": {
		"sys": {
		  "MethodName": "helloBase",
		  "UtcNow": "2023-10-12T20:13:49.640002Z",
		  "RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
		}
	  }
}`)
