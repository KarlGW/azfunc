package output

import (
	"testing"
	"time"

	"github.com/KarlGW/azfunc/data"
	"github.com/KarlGW/azfunc/eventgrid"
	"github.com/google/go-cmp/cmp"
)

func TestNewEventGrid(t *testing.T) {
	type testEvent struct {
		Message string `json:"message"`
	}

	var tests = []struct {
		name  string
		input struct {
			name    string
			options []EventGridOption
		}
		want *EventGrid
	}{
		{
			name: "defaults",
			input: struct {
				name    string
				options []EventGridOption
			}{
				name:    "event",
				options: nil,
			},
			want: &EventGrid{
				name: "event",
				data: nil,
			},
		},
		{
			name: "with options",
			input: struct {
				name    string
				options []EventGridOption
			}{
				name: "event",
				options: []EventGridOption{
					func(o *EventGridOptions) {
						o.Data = func() data.Raw {
							event, _ := eventgrid.NewCloudEvent("source", "type", testEvent{Message: "hello"}, func(o *eventgrid.CloudEventOptions) {
								o.ID = "12345"
								o.Time = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
							})
							return event.JSON()
						}()
					},
				},
			},
			want: &EventGrid{
				name: "event",
				data: data.Raw(`{"time":"2024-01-01T00:00:00Z","data":{"message":"hello"},"specversion":"1.0","type":"type","source":"source","id":"12345"}`),
			},
		},
	}

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
		got := &EventGrid{}
		got.Write([]byte(`{"data":{"message":"hello"},"specversion":"1.0","type":"type","source":"source","id":"12345","time":"2024-01-01T00:00:00+01:00"}`))
		want := &EventGrid{data: data.Raw(`{"data":{"message":"hello"},"specversion":"1.0","type":"type","source":"source","id":"12345","time":"2024-01-01T00:00:00+01:00"}`)}

		if diff := cmp.Diff(want, got, cmp.AllowUnexported(EventGrid{})); diff != "" {
			t.Errorf("Write() = unexpected result (-want +got)\n%s\n", diff)
		}
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
