package azfunc

import (
	"net/http"

	"github.com/potatoattack/azfunc/trigger"
)

// triggerable is the interface that wraps around the method run.
type triggerable interface {
	run(ctx *Context, r *http.Request) error
}

// GenericTriggerFunc represents a generic function to be executed by the function app.
type GenericTriggerFunc func(ctx *Context, trigger *trigger.Generic) error

// genericTrigger contains the trigger func, name and options of the trigger.
type genericTrigger struct {
	fn      GenericTriggerFunc
	name    string
	options []trigger.GenericOption
}

// run creates the trigger and runs the trigger func.
func (t genericTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := trigger.NewGeneric(r, t.name, t.options...)
	if err != nil {
		return err
	}
	return t.fn(ctx, tr)
}

// GenericTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func GenericTrigger(name string, fn GenericTriggerFunc, options ...trigger.GenericOption) FunctionOption {
	return func(f *function) {
		f.trigger = genericTrigger{
			fn:      fn,
			name:    name,
			options: options,
		}
	}
}

// HTTPTriggerFunc represents an HTTP trigger function to be executed by the function app.
type HTTPTriggerFunc func(ctx *Context, trigger *trigger.HTTP) error

// httpTrigger contains the trigger func and name of the trigger.
type httpTrigger struct {
	fn      HTTPTriggerFunc
	options []trigger.HTTPOption
}

// run creates the trigger and runs the trigger func.
func (t httpTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := trigger.NewHTTP(r, t.options...)
	if err != nil {
		return err
	}
	return t.fn(ctx, tr)
}

// HTTPTrigger takes the provided function and sets it as
// the function to be run by the trigger.
func HTTPTrigger(fn HTTPTriggerFunc, options ...trigger.HTTPOption) FunctionOption {
	return func(f *function) {
		f.trigger = httpTrigger{
			fn:      fn,
			options: options,
		}
	}
}

// TimerTriggerFunc represents a Timer trigger function tp be executed by the function app.
type TimerTriggerFunc func(ctx *Context, trigger *trigger.Timer) error

// timerTrigger contains the trigger func and name of the trigger.
type timerTrigger struct {
	fn      TimerTriggerFunc
	options []trigger.TimerOption
}

// run creates the trigger and runs the trigger func.
func (t timerTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := trigger.NewTimer(r, t.options...)
	if err != nil {
		return err
	}
	return t.fn(ctx, tr)
}

// TimerTrigger takes the provided function and sets it as
// the function to be run by the trigger.
func TimerTrigger(fn TimerTriggerFunc, options ...trigger.TimerOption) FunctionOption {
	return func(f *function) {
		f.trigger = timerTrigger{
			fn:      fn,
			options: options,
		}
	}
}

// QueueTriggerFunc represents a Queue Storage trigger function to be exexuted
// by the function app.
type QueueTriggerFunc func(ctx *Context, trigger *trigger.Queue) error

// queueTrigger contains the trigger func, name and options of the trigger.
type queueTrigger struct {
	fn      QueueTriggerFunc
	name    string
	options []trigger.QueueOption
}

// run creates the trigger and runs the trigger func.
func (t queueTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := trigger.NewQueue(r, t.name, t.options...)
	if err != nil {
		return err
	}
	return t.fn(ctx, tr)
}

// QueueTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func QueueTrigger(name string, fn QueueTriggerFunc, options ...trigger.QueueOption) FunctionOption {
	return func(f *function) {
		f.trigger = queueTrigger{
			fn:      fn,
			name:    name,
			options: options,
		}
	}
}

// ServiceBusTriggerFunc represents a Service Bus trigger function to be exexuted
// by the function app.
type ServiceBusTriggerFunc func(ctx *Context, trigger *trigger.ServiceBus) error

// serviceBusTrigger contains the trigger func, name and options of the trigger.
type serviceBusTrigger struct {
	fn      ServiceBusTriggerFunc
	name    string
	options []trigger.ServiceBusOption
}

// run creates the trigger and runs the trigger func.
func (t serviceBusTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := trigger.NewServiceBus(r, t.name, t.options...)
	if err != nil {
		return err
	}
	return t.fn(ctx, tr)
}

// ServiceBusTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func ServiceBusTrigger(name string, fn ServiceBusTriggerFunc, options ...trigger.ServiceBusOption) FunctionOption {
	return func(f *function) {
		f.trigger = serviceBusTrigger{
			fn:      fn,
			name:    name,
			options: options,
		}
	}
}

// EventGridTriggerFunc represents an Event Grid trigger function to be executed by
// the function app.
type EventGridTriggerFunc func(ctx *Context, trigger *trigger.EventGrid) error

// eventGridTrigger contains the trigger func, name and options of the trigger.
type eventGridTrigger struct {
	fn      EventGridTriggerFunc
	name    string
	options []trigger.EventGridOption
}

// run creates the trigger and runs the trigger func.
func (t eventGridTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := trigger.NewEventGrid(r, t.name, t.options...)
	if err != nil {
		return err
	}
	return t.fn(ctx, tr)
}

// EventGridTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func EventGridTrigger(name string, fn EventGridTriggerFunc, options ...trigger.EventGridOption) FunctionOption {
	return func(f *function) {
		f.trigger = eventGridTrigger{
			fn:      fn,
			name:    name,
			options: options,
		}
	}
}
