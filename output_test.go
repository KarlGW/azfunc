package azfunc

import (
	"net/http"
	"testing"

	"github.com/KarlGW/azfunc/bindings"
	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewOutput(t *testing.T) {
	var tests = []struct {
		name  string
		input []OutputOption
		want  Output
	}{
		{
			name:  "Output with defaults",
			input: nil,
			want: Output{
				Outputs:     map[string]bindable{},
				Logs:        nil,
				ReturnValue: nil,
			},
		},
		{
			name: "Output with options",
			input: []OutputOption{
				func(o *OutputOptions) {
					o.Bindings = []bindable{
						bindings.NewHTTP(),
						bindings.NewBase("queue"),
					}
					o.Logs = []string{"Log message"}
					o.ReturnValue = 0
				},
			},
			want: Output{
				Outputs: map[string]bindable{
					"queue": bindings.NewBase("queue"),
				},
				Logs:        []string{"Log message"},
				ReturnValue: 0,
				http:        bindings.NewHTTP(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewOutput(test.input...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Output{}, bindings.HTTP{}, bindings.Base{})); diff != "" {
				t.Errorf("NewOutput() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestOutput_JSON(t *testing.T) {
	var tests = []struct {
		name  string
		input Output
		want  []byte
	}{
		{
			name: "Parse output to JSON",
			input: Output{
				Outputs: map[string]bindable{
					"queue": bindings.NewBase("queue", func(o *bindings.Options) {
						o.Data = []byte(`{"message":"hello","number":3}`)
					}),
					/* 					"queue": &Base{
						name: "queue",
						data: []byte(`{"message":"hello","number":3}`),
					}, */
				},
				http: bindings.NewHTTP(func(o *bindings.Options) {
					o.StatusCode = http.StatusOK
					o.Body = data.Raw(`{"message":"hello","number":2}`)
					o.Header = http.Header{
						"Content-Type": {"application/json"},
					}
				}),
				/* http: &HTTP{
					StatusCode: http.StatusOK,
					Body:       data.Raw(`{"message":"hello","number":2}`),
					header: http.Header{
						"Content-Type": {"application/json"},
					},
				}, */
				Logs:        nil,
				ReturnValue: nil,
			},
			want: output1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.JSON()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("JSON() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

var output1 = []byte(`{"Outputs":{"queue":"{\"message\":\"hello\",\"number\":3}","res":{"headers":{"Content-Type":"application/json"},"statusCode":"200","body":"{\"message\":\"hello\",\"number\":2}"}},"ReturnValue":null,"Logs":null}`)
