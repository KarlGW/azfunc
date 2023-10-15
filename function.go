package azfunc

import (
	"github.com/KarlGW/azfunc/bindings"
	"github.com/KarlGW/azfunc/triggers"
)

// TriggerFunc represents a base function to be executed by the function app.
type TriggerFunc func(*Context, triggers.Base) error

// HTTPTriggerFunc represents an HTTP based function to be executed by the function app.
type HTTPTriggerFunc func(*Context, triggers.HTTP) error

// services is intended to hold custom services to be used within the
// Function App. Both services and clients both exists just for semantics,
// and either can be used.
type services map[string]any

// Add a service.
func (s services) Add(name string, service any) {
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
	c[name] = client
}

// Get a client.
func (c clients) Get(name string) any {
	return c[name]
}

// function is an internal structure that represents a function
// in a FunctionApp.
type function struct {
	name            string
	triggerName     string
	triggerFunc     TriggerFunc
	httpTriggerFunc HTTPTriggerFunc
	bindings        []bindings.Bindable
}

// Context represents the function context and contains output,
// bindings, services and clients.
type Context struct {
	// Output contains bindings.
	Output bindings.Output
	// services contains services defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	services services
	// clients contains clients defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	clients clients
}

// Services returns the services set in the Context.
func (c *Context) Services() services {
	return c.services
}

// Clients returns the clients set in the Context.
func (c *Context) Clients() clients {
	return c.clients
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

// Trigger takes the provided name and function and sets it to the
// function.
func Trigger(name string, fn TriggerFunc) FunctionOption {
	return func(f *function) {
		f.triggerName = name
		f.triggerFunc = fn
		f.httpTriggerFunc = nil
	}
}

// HTTPTrigger takes the provided function and sets it to the
// function.
func HTTPTrigger(fn HTTPTriggerFunc) FunctionOption {
	return func(f *function) {
		f.httpTriggerFunc = fn
		f.triggerFunc = nil
	}
}

// QueueTrigger takes the provided name and function and sets it to the
// function.
var QueueTrigger = Trigger

// TimerTrigger takes the provided name and function and sets it to the
// function.
var TimerTrigger = Trigger
