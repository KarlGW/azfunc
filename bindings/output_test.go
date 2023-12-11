package bindings

import (
	"net/http"
	"testing"

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
				Outputs:     map[string]Bindable{},
				Logs:        nil,
				ReturnValue: nil,
			},
		},
		{
			name: "Output with options",
			input: []OutputOption{
				func(o *OutputOptions) {
					o.Bindings = []Bindable{
						NewHTTP(),
						NewBase("queue"),
					}
					o.Logs = []string{"Log message"}
					o.ReturnValue = 0
				},
			},
			want: Output{
				Outputs: map[string]Bindable{
					"queue": &Base{name: "queue"},
				},
				Logs:        []string{"Log message"},
				ReturnValue: 0,
				http:        &HTTP{name: "res", StatusCode: http.StatusOK, header: http.Header{}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewOutput(test.input...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Output{}, HTTP{}, Base{})); diff != "" {
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
				Outputs: map[string]Bindable{
					"queue": &Base{
						name: "queue",
						data: []byte(`{"message":"hello","number":3}`),
					},
				},
				http: &HTTP{
					StatusCode: http.StatusOK,
					Body:       data.Raw(`{"message":"hello","number":2}`),
					header: http.Header{
						"Content-Type": {"application/json"},
					},
				},
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

var output1 = []byte(`{"Outputs":{"queue":"{\"message\":\"hello\",\"number\":3}","res":{"statusCode":"200","body":"{\"message\":\"hello\",\"number\":2}","headers":{"Content-Type":"application/json"}}},"Logs":null,"ReturnValue":null}`)
