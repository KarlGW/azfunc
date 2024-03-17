package triggers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KarlGW/azfunc/data"
)

// Timer represent a Timer trigger.
type Timer struct {
	ScheduleStatus TimerScheduleStatus
	Metadata       Metadata
	Schedule       TimerSchedule
	IsPastDue      bool
}

// TimerOptions contains options for a Timer trigger.
type TimerOptions struct {
	Name string
}

// TimerOption is a function that sets options on a Timer trigger.
type TimerOption func(o *TimerOptions)

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
// *http.Request. By default it will use the name "timer" for the
// trigger. This can be overridden with providing a name
// in the options.
func NewTimer(r *http.Request, options ...TimerOption) (*Timer, error) {
	opts := TimerOptions{
		Name: "timer",
	}
	for _, option := range options {
		option(&opts)
	}

	var t timerTrigger
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, ErrTriggerPayloadMalformed
	}
	defer r.Body.Close()

	d, ok := t.Data[opts.Name]
	if !ok {
		return nil, ErrTriggerNameIncorrect
	}
	d.Metadata = t.Metadata

	return &d, nil
}

// timerTrigger is the incoming request from the Function host.
type timerTrigger struct {
	Data     map[string]Timer
	Metadata Metadata
}
