package triggers

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewTimer(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			options []Option
		}
		want    *Timer
		wantErr error
	}{
		{
			name: "NewTimer",
			input: struct {
				req     *http.Request
				options []Option
			}{
				req: &http.Request{
					Body: io.NopCloser(bytes.NewBuffer(timerRequest1)),
				},
			},
			want: &Timer{
				Schedule: TimerSchedule{
					AdjustForDST: true,
				},
				ScheduleStatus: TimerScheduleStatus{
					Last:        _testTimerTime1,
					Next:        _testTimerTime1,
					LastUpdated: _testTimerTime1,
				},
				IsPastDue: false,
				Metadata: Metadata{
					Sys: MetadataSys{
						MethodName: "helloTimer",
						UtcNow:     _testTimerTime1,
						RandGuid:   "4e773554-f6b7-4ea2-b07d-4c5fd5aba741",
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewTimer(test.input.req)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("NewTimer() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("NewTimer() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

var timerRequest1 = []byte(`{
	"Data": {
	  "timer": {
		"Schedule": {
		  "AdjustForDST": true
		},
		"ScheduleStatus": {
		  "Last": "2023-10-12T20:13:49.640002Z",
		  "Next": "2023-10-12T20:13:49.640002Z",
		  "LastUpdated": "2023-10-12T20:13:49.640002Z"
		},
		"IsPastDue": false
	  }
	},
	"Metadata": {
	  "sys": {
		"MethodName": "helloTimer",
		"UtcNow": "2023-10-12T20:13:49.640002Z",
		"RandGuid": "4e773554-f6b7-4ea2-b07d-4c5fd5aba741"
	  }
	}
  }
`)

var _testTimerTime1, _ = time.Parse("2006-01-02T15:04:05.999999Z", "2023-10-12T20:13:49.640002Z")
