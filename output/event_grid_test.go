package output

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewEventGrid(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			name    string
			options []EventGridOption
		}
		want *EventGrid
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewEventGrid(test.want.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(EventGrid{})); diff != "" {
				t.Errorf("NewEventGrid() = unexpected (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestEventGrid_Write(t *testing.T) {
	t.Run("Write", func(t *testing.T) {

	})
}

func TestEventGrid_Name(t *testing.T) {
	var tests = []struct {
		name  string
		input *EventGrid
		want  string
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Name()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("EventGrid.Name() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
