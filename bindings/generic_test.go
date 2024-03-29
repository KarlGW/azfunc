package bindings

import (
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewGeneric(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			name    string
			options []GenericOption
		}
		want *Generic
	}{
		{
			name: "defaults",
			input: struct {
				name    string
				options []GenericOption
			}{
				name:    "queue",
				options: nil,
			},
			want: &Generic{
				name: "queue",
				data: nil,
			},
		},
		{
			name: "with options",
			input: struct {
				name    string
				options []GenericOption
			}{
				name: "queue",
				options: []GenericOption{
					func(o *GenericOptions) {
						o.Data = data.Raw(`{"message":"hello"}`)
					},
				},
			},
			want: &Generic{
				name: "queue",
				data: data.Raw(`{"message":"hello"}`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewGeneric(test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Generic{})); diff != "" {
				t.Errorf("NewGeneric() = unexpected (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestGeneric_Write(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		got := &Generic{}
		got.Write([]byte(`{"message":"hello"}`))
		want := &Generic{data: data.Raw(`{"message":"hello"}`)}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(Generic{})); diff != "" {
			t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestGeneric_Name(t *testing.T) {
	var tests = []struct {
		name  string
		input *Generic
		want  string
	}{
		{
			name:  "default",
			input: &Generic{},
			want:  "",
		},
		{
			name:  "with name",
			input: &Generic{name: "queue"},
			want:  "queue",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Name()

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Generic{})); diff != "" {
				t.Errorf("Name() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
