package azfunc

// Context represents the function context and contains output,
// bindings, services and clients.
type Context struct {
	// log contains a logger.
	log logger
	// services contains services defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	services services
	// clients contains clients defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	clients clients
	// outputs contains output bindings.
	outputs *outputs
}

// Outputs returns the outputs set in the Context.
func (c Context) Outputs() *outputs {
	return c.outputs
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
