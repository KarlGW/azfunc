package triggers

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/KarlGW/azfunc/data"
	"github.com/google/go-cmp/cmp"
)

func TestNewEventGrid(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			name    string
			options []EventGridOption
		}
		want    *EventGrid
		wantErr error
	}{
		{
			name: "NewEventGrid - cloud event",
			input: struct {
				req     *http.Request
				name    string
				options []EventGridOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(eventGridCloudEventRequest1)),
				},
				name: "event",
			},
			want: &EventGrid{
				ID:      "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
				Topic:   "topic",
				Subject: "subject",
				Type:    "created",
				Time:    _testEventGridTime1,
				Data:    data.Raw(`{"id":"4e773554-f6b7-4ea2-b07d-4c5fd5aba741","name":"test"}`),
				Schema:  EventGridSchemaCloudEvents,
				Metadata: EventGridMetadata{
					Data: data.Raw(`{"id":"4e773554-f6b7-4ea2-b07d-4c5fd5aba741","name":"test"}`),
					Metadata: Metadata{
						Sys: MetadataSys{
							MethodName: "testevent",
							UTCNow:     _testEventGridTime1,
							RandGuid:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
						},
					},
				},
			},
		},
		{
			name: "NewEventGrid - event grid",
			input: struct {
				req     *http.Request
				name    string
				options []EventGridOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(eventGridEventRequest1)),
				},
				name: "event",
			},
			want: &EventGrid{
				ID:      "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
				Topic:   "topic",
				Subject: "subject",
				Type:    "created",
				Time:    _testEventGridTime1,
				Data:    data.Raw(`{"id":"4e773554-f6b7-4ea2-b07d-4c5fd5aba741","name":"test"}`),
				Schema:  EventGridSchemaEventGrid,
				Metadata: EventGridMetadata{
					Data: data.Raw(`{"id":"4e773554-f6b7-4ea2-b07d-4c5fd5aba741","name":"test"}`),
					Metadata: Metadata{
						Sys: MetadataSys{
							MethodName: "testevent",
							UTCNow:     _testEventGridTime1,
							RandGuid:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewEventGrid(test.input.req, test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("NewEventGrid() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("NewEventGrid() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestEventGrid_Parse(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			name    string
			options []EventGridOption
		}
		want    eventGridTest
		wantErr error
	}{
		{
			name: "Parse",
			input: struct {
				req     *http.Request
				name    string
				options []EventGridOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(eventGridCloudEventRequest1)),
				},
				name: "event",
			},
			want: eventGridTest{
				ID:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
				Name: "test",
			},
		},
		{
			name: "Parse - pretty JSON",
			input: struct {
				req     *http.Request
				name    string
				options []EventGridOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(eventGridCloudEventRequest2)),
				},
				name: "event",
			},
			want: eventGridTest{
				ID:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
				Name: "test",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trigger, err := NewEventGrid(test.input.req, test.input.name, test.input.options...)
			if err != nil {
				t.Fatalf("NewEventGrid() = unexpected error: %v", err)
			}

			var got eventGridTest
			gotErr := trigger.Parse(&got)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse() = unexpected result (-want +got)\n%s\n", diff)
			}

			if test.wantErr == nil && gotErr != nil {
				t.Errorf("Parse() = unexpected error: %v\n", gotErr)
			}

		})
	}
}

func TestEventGridSchema_String(t *testing.T) {
	var tests = []struct {
		name  string
		input EventGridSchema
		want  string
	}{
		{
			name:  "cloud events",
			input: EventGridSchemaCloudEvents,
			want:  "CloudEvents",
		},
		{
			name:  "event grid",
			input: EventGridSchemaEventGrid,
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

var eventGridCloudEventRequest1 = []byte(`{
	"Data": {
	  "event": {
		"id": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
		"source": "topic",
		"specversion": "1.0",
		"type": "created",
		"subject": "subject",
		"time": "2023-10-12T20:13:49.640002Z",
		"data": {"id":"4e773554-f6b7-4ea2-b07d-4c5fd5aba741","name":"test"}
	  }
	},
	"Metadata": {
	  "data": {"id":"4e773554-f6b7-4ea2-b07d-4c5fd5aba741","name":"test"},
	  "sys": {
		"MethodName": "testevent",
		"UtcNow": "2023-10-12T20:13:49.640002Z",
		"RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
	  }
	}
  }
`)

var eventGridCloudEventRequest2 = []byte(`{
	"Data": {
	  "event": {
		"id": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
		"source": "topic",
		"specversion": "1.0",
		"type": "created",
		"subject": "subject",
		"time": "2023-10-12T20:13:49.640002Z",
		"data": {
		  "id": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
		  "name": "test"
		}
	  }
	},
	"Metadata": {
	  "data": {
		"id": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
		"name": "test"
	  },
	  "sys": {
		"MethodName": "testevent",
		"UtcNow": "2023-10-12T20:13:49.640002Z",
		"RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
	  }
	}
  }
`)

var eventGridEventRequest1 = []byte(`{
	"Data": {
	  "event": {
		"id": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
		"topic": "topic",
		"subject": "subject",
		"eventType": "created",
		"dataVersion": "1",
		"metadataVersion": "1",
		"eventTime": "2023-10-12T20:13:49.640002Z",
		"data": {"id":"4e773554-f6b7-4ea2-b07d-4c5fd5aba741","name":"test"}
	  }
	},
	"Metadata": {
	  "data": {"id":"4e773554-f6b7-4ea2-b07d-4c5fd5aba741","name":"test"},
	  "sys": {
		"MethodName": "testevent",
		"UtcNow": "2023-10-12T20:13:49.640002Z",
		"RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
	  }
	}
  }
`)

var _testEventGridTime1, _ = time.Parse("2006-01-02T15:04:05.999999Z", "2023-10-12T20:13:49.640002Z")

type eventGridTest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
