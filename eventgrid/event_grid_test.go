package eventgrid

import "testing"

func TestSchema_String(t *testing.T) {
	var tests = []struct {
		name  string
		input Schema
		want  string
	}{
		{
			name:  "cloud events",
			input: SchemaCloudEvents,
			want:  "CloudEvents",
		},
		{
			name:  "event grid",
			input: SchemaEventGrid,
			want:  "EventGrid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.String()

			if test.want != got {
				t.Errorf("String() = unexpected result, want: %s, got: %s\n", test.want, got)
			}
		})
	}
}
