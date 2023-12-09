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

func TestNewQueueStorage(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			name    string
			options []Option
		}
		want    *QueueStorage
		wantErr error
	}{
		{
			name: "NewQueueStorage",
			input: struct {
				req     *http.Request
				name    string
				options []Option
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(queueStorageRequest1)),
				},
				name: "queue",
			},
			want: &QueueStorage{
				data: data.Raw(`{"message":"hello","number":2}`),
				Metadata: QueueStorageMetadata{
					DequeueCount:    "1",
					ID:              "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
					PopReceipt:      "STRING",
					ExpirationTime:  _testQueueStorageTime1,
					InsertionTime:   _testQueueStorageTime1,
					NextVisibleTime: _testQueueStorageTime1,
					Metadata: Metadata{
						Sys: MetadataSys{
							MethodName: "helloQueue",
							UtcNow:     _testQueueStorageTime1,
							RandGuid:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewQueueStorage(test.input.req, test.input.name, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(QueueStorage{})); diff != "" {
				t.Errorf("NewQueueStorage() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("NewQueueStorage() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

var queueStorageRequest1 = []byte(`{
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

var _testQueueStorageTime1, _ = time.Parse("2006-01-02T15:04:05.999999Z", "2023-10-12T20:13:49.640002Z")
