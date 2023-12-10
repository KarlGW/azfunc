package bindings

import (
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewQueue(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			name    string
			options []Option
		}
		want *Queue
	}{
		{
			name: "defaults",
			input: struct {
				name    string
				options []Option
			}{
				name:    "queue",
				options: nil,
			},
			want: &Queue{
				name: "queue",
				Raw:  nil,
			},
		},
		{
			name: "with options",
			input: struct {
				name    string
				options []Option
			}{
				name: "queue",
				options: []Option{
					func(o *Options) {
						o.Data = data.Raw(`{"message":"hello"}`)
					},
				},
			},
			want: &Queue{
				name: "queue",
				Raw:  data.Raw(`{"message":"hello"}`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewQueue(test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Queue{})); diff != "" {
				t.Errorf("NewQueue() = unexpected (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestQueue_Write(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		got := &Queue{}
		got.Write([]byte(`{"message":"hello"}`))
		want := &Queue{Raw: data.Raw(`{"message":"hello"}`)}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(Queue{})); diff != "" {
			t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestQueue_Name(t *testing.T) {
	var tests = []struct {
		name  string
		input *Queue
		want  string
	}{
		{
			name:  "default",
			input: &Queue{},
			want:  "",
		},
		{
			name:  "with name",
			input: &Queue{name: "queue"},
			want:  "queue",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Name()

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Queue{})); diff != "" {
				t.Errorf("Name() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
