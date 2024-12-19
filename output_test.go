package azfunc

import (
	"net/http"
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/KarlGW/azfunc/output"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewOutput(t *testing.T) {
	var tests = []struct {
		name  string
		input []outputsOption
		want  *outputs
	}{
		{
			name:  "Output with defaults",
			input: nil,
			want: &outputs{
				outputs:     map[string]outputable{},
				returnValue: nil,
				log:         newInvocationLogger(),
			},
		},
		{
			name: "Output with options",
			input: []outputsOption{
				func(o *outputsOptions) {
					o.outputs = []outputable{
						output.NewHTTP(),
						output.NewGeneric("queue"),
					}
				},
			},
			want: &outputs{
				outputs: map[string]outputable{
					"queue": output.NewGeneric("queue"),
				},
				http: output.NewHTTP(),
				log:  newInvocationLogger(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newOutputs(test.input...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(outputs{}, output.HTTP{}, output.Generic{}, invocationLogger{}), cmpopts.IgnoreFields(invocationLogger{}, "l", "w", "mu")); diff != "" {
				t.Errorf("NewOutput() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestOutput_JSON(t *testing.T) {
	var tests = []struct {
		name  string
		input outputs
		want  []byte
	}{
		{
			name: "Parse output to JSON",
			input: outputs{
				outputs: map[string]outputable{
					"queue": output.NewGeneric("queue", func(o *output.GenericOptions) {
						o.Data = []byte(`{"message":"hello","number":3}`)
					}),
				},
				http: output.NewHTTP(func(o *output.HTTPOptions) {
					o.StatusCode = http.StatusOK
					o.Body = data.Raw(`{"message":"hello","number":2}`)
					o.Header = http.Header{
						"Content-Type": {"application/json"},
					}
				}),
				returnValue: nil,
				log:         newInvocationLogger(),
			},
			want: output1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.json()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("JSON() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

var output1 = []byte(`{"Outputs":{"queue":"{\"message\":\"hello\",\"number\":3}","res":{"headers":{"Content-Type":"application/json"},"statusCode":"200","body":"{\"message\":\"hello\",\"number\":2}"}},"ReturnValue":null,"Logs":null}`)
