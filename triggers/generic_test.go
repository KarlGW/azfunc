package triggers

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewGeneric(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			name    string
			options []GenericOption
		}
		want    *Generic
		wantErr error
	}{
		{
			name: "NewGeneric",
			input: struct {
				req     *http.Request
				name    string
				options []GenericOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(genericRequest1)),
				},
				name: "generic",
			},
			want: &Generic{
				Data: data.Raw(`{"message":"hello","number":2}`),
				Metadata: map[string]any{
					"sys": map[string]any{
						"MethodName": "helloGeneric",
						"UtcNow":     "2023-10-12T20:13:49.640002Z",
						"RandGuid":   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewGeneric(test.input.req, test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Generic{})); diff != "" {
				t.Errorf("NewGeneric() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("NewGeneric() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

var genericRequest1 = []byte(`{
	"Data": {
		"generic": "{\"message\":\"hello\",\"number\":2}"
	  },
	  "Metadata": {
		"sys": {
		  "MethodName": "helloGeneric",
		  "UtcNow": "2023-10-12T20:13:49.640002Z",
		  "RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
		}
	  }
}`)
