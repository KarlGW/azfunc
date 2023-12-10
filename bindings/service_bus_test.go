package bindings

import (
	"testing"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewServiceBus(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			name    string
			options []Option
		}
		want *ServiceBus
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
			want: &ServiceBus{
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
			want: &ServiceBus{
				name: "queue",
				Raw:  data.Raw(`{"message":"hello"}`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewServiceBus(test.want.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(ServiceBus{})); diff != "" {
				t.Errorf("NewServiceBus() = unexpected (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestServiceBus_Write(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		got := &ServiceBus{}
		got.Write([]byte(`{"message":"hello"}`))
		want := &ServiceBus{Raw: data.Raw(`{"message":"hello"}`)}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(ServiceBus{})); diff != "" {
			t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
		}
	})
}

func TestServiceBus_Name(t *testing.T) {
	var tests = []struct {
		name  string
		input *ServiceBus
		want  string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Name()

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(ServiceBus{})); diff != "" {
				t.Errorf("Name() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
