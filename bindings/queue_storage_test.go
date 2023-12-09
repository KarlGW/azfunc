package bindings

import (
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewQueueStorage(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			name    string
			options []Option
		}
		want *QueueStorage
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
			want: &QueueStorage{
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
			want: &QueueStorage{
				name: "queue",
				Raw:  data.Raw(`{"message":"hello"}`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewQueueStorage(test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(QueueStorage{})); diff != "" {
				t.Run(test.name, func(t *testing.T) {
					t.Errorf("NewQueueStorage() = unexpected (-want +got)\n%s\n", diff)
				})
			}
		})
	}
}

func TestQueueStorage_Write(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		got := &QueueStorage{}
		got.Write([]byte(`{"message":"hello"}`))
		want := &QueueStorage{Raw: data.Raw(`{"message":"hello"}`)}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(QueueStorage{})); diff != "" {
			t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestQueueStorage_Name(t *testing.T) {
	var tests = []struct {
		name  string
		input *Base
		want  string
	}{
		{
			name:  "default",
			input: &Base{},
			want:  "",
		},
		{
			name:  "with name",
			input: &Base{name: "queue"},
			want:  "queue",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Name()

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Base{})); diff != "" {
				t.Errorf("Name() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
