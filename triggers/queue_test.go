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

func TestNewQueue(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			name    string
			options []QueueOption
		}
		want    *Queue
		wantErr error
	}{
		{
			name: "NewQueue",
			input: struct {
				req     *http.Request
				name    string
				options []QueueOption
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueRequest1)),
				},
				name: "queue",
			},
			want: &Queue{
				data: data.Raw(`{"message":"hello","number":2}`),
				Metadata: QueueMetadata{
					DequeueCount:    "1",
					ID:              "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
					PopReceipt:      "STRING",
					ExpirationTime:  _testQueueTime1,
					InsertionTime:   _testQueueTime1,
					NextVisibleTime: _testQueueTime1,
					Metadata: Metadata{
						Sys: MetadataSys{
							MethodName: "helloQueue",
							UTCNow:     _testQueueTime1,
							RandGuid:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewQueue(test.input.req, test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Queue{})); diff != "" {
				t.Errorf("NewQueue() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("NewQueue() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

var queueRequest1 = []byte(`{
	"Data": {
		"queue": "{\"message\":\"hello\",\"number\":2}"
	  },
	  "Metadata": {
		"DequeueCount": "1",
		"ExpirationTime": "2023-10-12T20:13:49.640002Z",
		"Id": "\"4e773554-f6b7-4ea2-b07d-4c5fd5aba741\"",
		"InsertionTime": "2023-10-12T20:13:49.640002Z",
		"NextVisibleTime": "2023-10-12T20:13:49.640002Z",
		"PopReceipt": "\"STRING\"",
		"sys": {
		  "MethodName": "helloQueue",
		  "UtcNow": "2023-10-12T20:13:49.640002Z",
		  "RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
		},
		"heyoooo": "\"hello\""
	  }
}`)

var _testQueueTime1, _ = time.Parse("2006-01-02T15:04:05.999999Z", "2023-10-12T20:13:49.640002Z")
