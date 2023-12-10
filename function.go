package azfunc

import (
	"github.com/KarlGW/azfunc/bindings"
	"github.com/KarlGW/azfunc/triggers"
)

// TriggerFunc represents a base function to be executed by the function app.
type TriggerFunc func(ctx *Context, trigger *triggers.Base)

// HTTPTriggerFunc represents an HTTP based function to be executed by the function app.
type HTTPTriggerFunc func(ctx *Context, trigger *triggers.HTTP)

// QueueTriggerFunc represents a Queue Storage based function to be exexuted
// by the function app.
type QueueTriggerFunc func(ctx *Context, trigger *triggers.Queue)

// ServiceBusTriggerFunc represents a Service Bus based function to be exexuted
// by the function app.
type ServiceBusTriggerFunc func(ctx *Context, trigger *triggers.ServiceBus)

// TimerTriggerFunc represents a Timer based function tp be executed by the function app.
type TimerTriggerFunc func(ctx *Context, trigger *triggers.Timer)

// services is intended to hold custom services to be used within the
// Function App. Both services and clients both exists just for semantics,
// and either can be used.
type services map[string]any

// Add a service.
func (s services) Add(name string, service any) {
	if s == nil {
		s = make(services)
	}
	s[name] = service
}

// Get a service.
func (s services) Get(name string) any {
	return s[name]
}

// clients is intended to hold custom clients to be used within the
// Function App. Both clients and services both exists just for semantics,
// and either can be used.
type clients map[string]any

// Add a client.
func (c clients) Add(name string, client any) {
	if c == nil {
		c = make(clients)
	}
	c[name] = client
}

// Get a client.
func (c clients) Get(name string) any {
	return c[name]
}

// function is an internal structure that represents a function
// in a FunctionApp.
type function struct {
	name        string
	triggerName string
	trigger     any
	bindings    []bindings.Bindable
}

// Context represents the function context and contains output,
// bindings, services and clients.
type Context struct {
	// Output contains bindings.
	Output bindings.Output
	// log contains a logger.
	log logger
	// services contains services defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	services services
	// clients contains clients defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	clients clients
	// err contains error set to the Context.
	err error
}

// Log returns the logger of the Context.
func (c Context) Log() logger {
	return c.log
}

// Err returns the error set to the Context.
func (c Context) Err() error {
	return c.err
}

// SetError sets an error to the Context.
func (c *Context) SetError(err error) {
	c.err = err
}

// Services returns the services set in the Context.
func (c *Context) Services() services {
	return c.services
}

// Clients returns the clients set in the Context.
func (c *Context) Clients() clients {
	return c.clients
}

// SetLogger sets a logger to the Context. Should not be used in most
// use-cases due to it being set by the FunctionApp.
func (c *Context) SetLogger(l logger) {
	c.log = l
}

// FunctionOption sets options to the function.
type FunctionOption func(f *function)

// Binding sets the provided binding to the function.
func Binding(binding bindings.Bindable) FunctionOption {
	return func(f *function) {
		if f.bindings == nil {
			f.bindings = []bindings.Bindable{binding}
			return
		}
		f.bindings = append(f.bindings, binding)
	}
}

// Trigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func Trigger(name string, fn TriggerFunc) FunctionOption {
	return func(f *function) {
		f.triggerName = name
		f.trigger = fn
	}
}

// HTTPTrigger takes the provided function and sets it as
// the function to be run by the trigger.
func HTTPTrigger(fn HTTPTriggerFunc) FunctionOption {
	return func(f *function) {
		f.triggerName = "req"
		f.trigger = fn
	}
}

// TimerTrigger takes the provided function and sets it as
// the function to be run by the trigger.
func TimerTrigger(fn TimerTriggerFunc) FunctionOption {
	return func(f *function) {
		f.triggerName = "timer"
		f.trigger = fn
	}
}

// QueueTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func QueueTrigger(name string, fn QueueTriggerFunc) FunctionOption {
	return func(f *function) {
		f.triggerName = name
		f.trigger = fn
	}
}
