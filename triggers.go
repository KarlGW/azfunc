package azfunc

import (
	"net/http"

	"github.com/KarlGW/azfunc/triggers"
)

// runnable is the interface that wraps around the method run.
type runnable interface {
	run(r *http.Request, ctx *Context, options ...triggers.Option) error
}

// TriggerFunc represents a base function to be executed by the function app.
type TriggerFunc func(ctx *Context, trigger *triggers.Base)

// trigger contains the trigger func and name of the trigger.
type trigger struct {
	fn   TriggerFunc
	name string
}

// run creates the trigger and runs the trigger func.
func (t trigger) run(r *http.Request, ctx *Context, options ...triggers.Option) error {
	tr, err := triggers.NewBase(r, t.name, options...)
	if err != nil {
		return err
	}
	t.fn(ctx, tr)
	return nil
}

// HTTPTriggerFunc represents an HTTP based function to be executed by the function app.
type HTTPTriggerFunc func(ctx *Context, trigger *triggers.HTTP)

// httpTrigger contains the trigger func and name of the trigger.
type httpTrigger struct {
	fn   HTTPTriggerFunc
	name string
}

// run creates the trigger and runs the trigger func.
func (t httpTrigger) run(r *http.Request, ctx *Context, options ...triggers.Option) error {
	tr, err := triggers.NewHTTP(r, options...)
	if err != nil {
		return nil
	}
	t.fn(ctx, tr)
	return nil
}

// TimerTriggerFunc represents a Timer based function tp be executed by the function app.
type TimerTriggerFunc func(ctx *Context, trigger *triggers.Timer)

// timerTrigger contains the trigger func and name of the trigger.
type timerTrigger struct {
	fn   TimerTriggerFunc
	name string
}

// run creates the trigger and runs the trigger func.
func (t timerTrigger) run(r *http.Request, ctx *Context, options ...triggers.Option) error {
	tr, err := triggers.NewTimer(r, options...)
	if err != nil {
		return nil
	}
	t.fn(ctx, tr)
	return nil
}

// QueueTriggerFunc represents a Queue Storage based function to be exexuted
// by the function app.
type QueueTriggerFunc func(ctx *Context, trigger *triggers.Queue)

// queueTrigger contains the trigger func and name of the trigger.
type queueTrigger struct {
	fn   QueueTriggerFunc
	name string
}

// run creates the trigger and runs the trigger func.
func (t queueTrigger) run(r *http.Request, ctx *Context, options ...triggers.Option) error {
	tr, err := triggers.NewQueue(r, t.name, options...)
	if err != nil {
		return nil
	}
	t.fn(ctx, tr)
	return nil
}

// ServiceBusTriggerFunc represents a Service Bus based function to be exexuted
// by the function app.
type ServiceBusTriggerFunc func(ctx *Context, trigger *triggers.ServiceBus)

// serviceBusTrigger contains the trigger func and name of the trigger.
type serviceBusTrigger struct {
	fn   ServiceBusTriggerFunc
	name string
}

// run creates the trigger and runs the trigger func.
func (t serviceBusTrigger) run(r *http.Request, ctx *Context, options ...triggers.Option) error {
	tr, err := triggers.NewServiceBus(r, t.name, options...)
	if err != nil {
		return nil
	}
	t.fn(ctx, tr)
	return nil
}
