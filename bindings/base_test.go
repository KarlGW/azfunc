package bindings

import (
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewBase(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			name    string
			options []Option
		}
		want *Base
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
			want: &Base{
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
			want: &Base{
				name: "queue",
				Raw:  data.Raw(`{"message":"hello"}`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewBase(test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Base{})); diff != "" {
				t.Errorf("NewBase() = unexpected (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestBase_Write(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		got := &Base{}
		got.Write([]byte(`{"message":"hello"}`))
		want := &Base{Raw: data.Raw(`{"message":"hello"}`)}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(Base{})); diff != "" {
			t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestBase_Name(t *testing.T) {
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
