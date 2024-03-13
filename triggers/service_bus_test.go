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

func TestNewServiceBus(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			name    string
			options []ServiceBusOption
		}
		want    *ServiceBus
		wantErr error
	}{
		{
			name: "NewServiceBus",
			input: struct {
				req     *http.Request
				name    string
				options []ServiceBusOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(serviceBusRequest1)),
				},
				name: "queue",
			},
			want: &ServiceBus{
				data: data.Raw(`{"message":"hello","number":2}`),
				Metadata: ServiceBusMetadata{
					Client: ServiceBusMetadataClient{
						FullyQualifiedNamespace: "namespace",
						Identifier:              "namespace-4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
						TransportType:           0,
						IsClosed:                false,
					},
					MessageReceiver:       map[string]any{},
					MessageSession:        map[string]any{},
					MessageActions:        map[string]any{},
					SessionActions:        map[string]any{},
					ReceiveActions:        map[string]any{},
					ApplicationProperties: map[string]any{},
					UserProperties:        map[string]any{},
					DeliveryCount:         "1",
					LockToken:             "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
					MessageID:             "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
					ContentType:           "application/json",
					SequenceNumber:        "1",
					ExpiresAtUTC:          _testServiceBusTime1ISO8601,
					ExpiresAt:             _testServiceBusTimeISO8601TZ,
					EnqueuedTimeUTC:       _testServiceBusTime1ISO8601,
					EnqueuedTime:          _testServiceBusTimeISO8601TZ,
					Metadata: Metadata{
						Sys: MetadataSys{
							MethodName: "helloQueue",
							UTCNow:     _testServiceBusTime1,
							RandGuid:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewServiceBus(test.input.req, test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(ServiceBus{}, TimeISO8601{})); diff != "" {
				t.Errorf("NewServiceBus() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("NewServiceBus() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

var serviceBusRequest1 = []byte(`{
	"Data": {
		"queue": "{\"message\":\"hello\",\"number\":2}"
	},
	"Metadata": {
	  "MessageReceiver": {},
	  "MessageSession": {},
	  "MessageActions": {},
	  "SessionActions": {},
	  "ReceiveActions": {},
	  "Client": {
		"FullyQualifiedNamespace": "namespace",
		"IsClosed": false,
		"TransportType": 0,
		"Identifier": "namespace-4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
	  },
	  "DeliveryCount": "1",
	  "LockToken": "\"4e773554-f6b7-4ea2-b07d-4c5fd5aba741\"",
	  "ExpiresAtUtc": "2023-10-12T20:13:49",
	  "ExpiresAt": "2023-10-12T20:13:49+00:00",
	  "EnqueuedTimeUtc": "2023-10-12T20:13:49",
	  "EnqueuedTime": "2023-10-12T20:13:49+00:00",
	  "MessageId": "\"4e773554-f6b7-4ea2-b07d-4c5fd5aba741\"",
	  "ContentType": "\"application/json\"",
	  "SequenceNumber": "1",
	  "ApplicationProperties": {},
	  "UserProperties": {},
	  "sys": {
		"MethodName": "helloQueue",
		"UtcNow": "2023-10-12T20:13:49.640002Z",
		"RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
	  },
	  "message": "\"hello\""
	}
  }`)

var (
	_testServiceBusTime1, _             = time.Parse("2006-01-02T15:04:05.999999Z", "2023-10-12T20:13:49.640002Z")
	_testServiceBusTime1ISO8601Raw, _   = time.Parse(iso8601, "2023-10-12T20:13:49")
	_testServiceBusTime1ISO8601         = TimeISO8601{Time: _testServiceBusTime1ISO8601Raw}
	_testServiceBusTime1ISO8601TZRaw, _ = time.Parse(iso8601TZ, "2023-10-12T20:13:49+00:00")
	_testServiceBusTimeISO8601TZ        = TimeISO8601{Time: _testServiceBusTime1ISO8601TZRaw, tz: true}
)
