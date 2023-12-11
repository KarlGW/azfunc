package triggers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KarlGW/azfunc/data"
)

// Timer represent a Timer trigger.
type Timer struct {
	Schedule       TimerSchedule
	ScheduleStatus TimerScheduleStatus
	IsPastDue      bool
	Metadata       Metadata
}

// TimerSchedule represents the Schedule field from the incoming
// request.
type TimerSchedule struct {
	AdjustForDST bool
}

// TimerScheduleStatus represents the ScheduleStatus field from
// the incoming request.
type TimerScheduleStatus struct {
	Last        time.Time
	Next        time.Time
	LastUpdated time.Time
}

// Parse together with Data satisfies the Triggerable interface. It
// is a no-op method. A timer trigger contains no other data needed
// to be parsed. Use the fields directly.
func (t Timer) Parse(v any) error {
	return nil
}

// Data together with Parse satisfies the Triggerable interface. It
// is a no-op method. A timer trigger contains no other data needed
// to be parsed. Use the fields directly.
func (t Timer) Data() data.Raw {
	return nil
}

// NewTimer creates and returns a Timer trigger from the provided
// *http.Request.
func NewTimer(r *http.Request, options ...Option) (*Timer, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	var t timerTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, ErrTriggerPayloadMalformed
	}
	defer r.Body.Close()

	return &Timer{
		Schedule:       t.Data.Timer.Schedule,
		ScheduleStatus: t.Data.Timer.ScheduleStatus,
		IsPastDue:      t.Data.Timer.IsPastDue,
		Metadata:       t.Metadata,
	}, nil
}

// timerTrigger is the incoming request from the Function host.
type timerTrigger struct {
	Data struct {
		Timer struct {
			Schedule       TimerSchedule
			ScheduleStatus TimerScheduleStatus
			IsPastDue      bool
		} `json:"timer"`
	}
	Metadata Metadata
}
