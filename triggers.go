package azfunc

import (
	"net/http"

	"github.com/KarlGW/azfunc/triggers"
)

// triggerable is the interface that wraps around the method run.
type triggerable interface {
	run(ctx *Context, r *http.Request) error
}

// TriggerFunc represents a base function to be executed by the function app.
type TriggerFunc func(ctx *Context, trigger *triggers.Base)

// trigger contains the trigger func and name of the trigger.
type trigger struct {
	fn      TriggerFunc
	name    string
	options []triggers.BaseOption
}

// run creates the trigger and runs the trigger func.
func (t trigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewBase(r, t.name, t.options...)
	if err != nil {
		return err
	}
	t.fn(ctx, tr)
	return nil
}

// Trigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func Trigger(name string, fn TriggerFunc, options ...triggers.BaseOption) FunctionOption {
	return func(f *function) {
		f.trigger = trigger{
			fn:      fn,
			name:    name,
			options: options,
		}
	}
}

// HTTPTriggerFunc represents an HTTP based function to be executed by the function app.
type HTTPTriggerFunc func(ctx *Context, trigger *triggers.HTTP)

// httpTrigger contains the trigger func and name of the trigger.
type httpTrigger struct {
	fn      HTTPTriggerFunc
	options []triggers.HTTPOption
}

// run creates the trigger and runs the trigger func.
func (t httpTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewHTTP(r, t.options...)
	if err != nil {
		return nil
	}
	t.fn(ctx, tr)
	return nil
}

// HTTPTrigger takes the provided function and sets it as
// the function to be run by the trigger.
func HTTPTrigger(fn HTTPTriggerFunc, options ...triggers.HTTPOption) FunctionOption {
	return func(f *function) {
		f.trigger = httpTrigger{
			fn:      fn,
			options: options,
		}
	}
}

// TimerTriggerFunc represents a Timer based function tp be executed by the function app.
type TimerTriggerFunc func(ctx *Context, trigger *triggers.Timer)

// timerTrigger contains the trigger func and name of the trigger.
type timerTrigger struct {
	fn      TimerTriggerFunc
	options []triggers.TimerOption
}

// run creates the trigger and runs the trigger func.
func (t timerTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewTimer(r, t.options...)
	if err != nil {
		return nil
	}
	t.fn(ctx, tr)
	return nil
}

// TimerTrigger takes the provided function and sets it as
// the function to be run by the trigger.
func TimerTrigger(fn TimerTriggerFunc, options ...triggers.TimerOption) FunctionOption {
	return func(f *function) {
		f.trigger = timerTrigger{
			fn:      fn,
			options: options,
		}
	}
}

// QueueTriggerFunc represents a Queue Storage based function to be exexuted
// by the function app.
type QueueTriggerFunc func(ctx *Context, trigger *triggers.Queue)

// queueTrigger contains the trigger func and name of the trigger.
type queueTrigger struct {
	fn      QueueTriggerFunc
	name    string
	options []triggers.QueueOption
}

// run creates the trigger and runs the trigger func.
func (t queueTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewQueue(r, t.name, t.options...)
	if err != nil {
		return nil
	}
	t.fn(ctx, tr)
	return nil
}

// QueueTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func QueueTrigger(name string, fn QueueTriggerFunc) FunctionOption {
	return func(f *function) {
		f.trigger = queueTrigger{
			fn:   fn,
			name: name,
		}
	}
}

// ServiceBusTriggerFunc represents a Service Bus based function to be exexuted
// by the function app.
type ServiceBusTriggerFunc func(ctx *Context, trigger *triggers.ServiceBus)

// serviceBusTrigger contains the trigger func and name of the trigger.
type serviceBusTrigger struct {
	fn      ServiceBusTriggerFunc
	name    string
	options []triggers.ServiceBusOption
}

// run creates the trigger and runs the trigger func.
func (t serviceBusTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewServiceBus(r, t.name, t.options...)
	if err != nil {
		return nil
	}
	t.fn(ctx, tr)
	return nil
}

// ServiceBusTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func ServiceBusTrigger(name string, fn ServiceBusTriggerFunc, options ...triggers.ServiceBusOption) FunctionOption {
	return func(f *function) {
		f.trigger = serviceBusTrigger{
			fn:      fn,
			name:    name,
			options: options,
		}
	}
}
