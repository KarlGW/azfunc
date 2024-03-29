package azfunc

import (
	"net/http"

	"github.com/KarlGW/azfunc/triggers"
)

// triggerable is the interface that wraps around the method run.
type triggerable interface {
	run(ctx *Context, r *http.Request) error
}

// GenericTriggerFunc represents a generic function to be executed by the function app.
type GenericTriggerFunc func(ctx *Context, trigger *triggers.Generic)

// genericTrigger contains the trigger func, name and options of the trigger.
type genericTrigger struct {
	fn      GenericTriggerFunc
	name    string
	options []triggers.GenericOption
}

// run creates the trigger and runs the trigger func.
func (t genericTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewGeneric(r, t.name, t.options...)
	if err != nil {
		return err
	}
	t.fn(ctx, tr)
	return nil
}

// GenericTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func GenericTrigger(name string, fn GenericTriggerFunc, options ...triggers.GenericOption) FunctionOption {
	return func(f *function) {
		f.trigger = genericTrigger{
			fn:      fn,
			name:    name,
			options: options,
		}
	}
}

// HTTPTriggerFunc represents an HTTP trigger function to be executed by the function app.
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
		return err
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

// TimerTriggerFunc represents a Timer trigger function tp be executed by the function app.
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
		return err
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

// QueueTriggerFunc represents a Queue Storage trigger function to be exexuted
// by the function app.
type QueueTriggerFunc func(ctx *Context, trigger *triggers.Queue)

// queueTrigger contains the trigger func, name and options of the trigger.
type queueTrigger struct {
	fn      QueueTriggerFunc
	name    string
	options []triggers.QueueOption
}

// run creates the trigger and runs the trigger func.
func (t queueTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewQueue(r, t.name, t.options...)
	if err != nil {
		return err
	}
	t.fn(ctx, tr)
	return nil
}

// QueueTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func QueueTrigger(name string, fn QueueTriggerFunc, options ...triggers.QueueOption) FunctionOption {
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
type ServiceBusTriggerFunc func(ctx *Context, trigger *triggers.ServiceBus)

// serviceBusTrigger contains the trigger func, name and options of the trigger.
type serviceBusTrigger struct {
	fn      ServiceBusTriggerFunc
	name    string
	options []triggers.ServiceBusOption
}

// run creates the trigger and runs the trigger func.
func (t serviceBusTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewServiceBus(r, t.name, t.options...)
	if err != nil {
		return err
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

// EventGridTriggerFunc represents an Event Grid trigger function to be executed by
// the function app.
type EventGridTriggerFunc func(ctx *Context, trigger *triggers.EventGrid)

// eventGridTrigger contains the trigger func, name and options of the trigger.
type eventGridTrigger struct {
	fn      EventGridTriggerFunc
	name    string
	options []triggers.EventGridOption
}

// run creates the trigger and runs the trigger func.
func (t eventGridTrigger) run(ctx *Context, r *http.Request) error {
	tr, err := triggers.NewEventGrid(r, t.name, t.options...)
	if err != nil {
		return err
	}
	t.fn(ctx, tr)
	return nil
}

// EventGridTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func EventGridTrigger(name string, fn EventGridTriggerFunc, options ...triggers.EventGridOption) FunctionOption {
	return func(f *function) {
		f.trigger = eventGridTrigger{
			fn:      fn,
			name:    name,
			options: options,
		}
	}
}
