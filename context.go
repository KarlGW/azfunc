package azfunc

import (
	"context"
)

// Context represents the function context and contains output,
// bindings, services and clients.
type Context struct {
	context.Context
	// log contains a logger.
	log logger
	// services contains services defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	services services
	// clients contains clients defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	clients clients
	// Outputs contains output bindings.
	Outputs *outputs
}

// Log returns the logger of the Context.
func (c Context) Log() logger {
	return c.log
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

// contextOptions contains options for creating a new Context.
type contextOptions struct {
	outputs  *outputs
	log      logger
	services services
	clients  clients
}

// contextOption is a function that sets options on a Context.
type contextOption func(o *contextOptions)

// newContext creates a new Context from the provided context.Context.
func newContext(ctx context.Context, options ...contextOption) *Context {
	opts := contextOptions{}
	for _, option := range options {
		option(&opts)
	}

	var c *Context
	if ctx, ok := ctx.(*Context); !ok {
		c = &Context{
			Context: ctx,
		}
	} else {
		c = ctx
	}

	c.Outputs = opts.outputs
	c.log = opts.log
	c.services = opts.services
	c.clients = opts.clients

	return c
}
