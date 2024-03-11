package azfunc

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	functionsCustomHandlerPort = "FUNCTIONS_CUSTOMHANDLER_PORT"
)

var (
	// ErrNoFunction is returned when no function has been set to the
	// FunctionApp.
	ErrNoFunction = errors.New("at least one function must be set")
)

// FunctionApp represents a Function App with its configuration
// and functions.
type FunctionApp struct {
	httpServer *http.Server
	router     *http.ServeMux
	// functions that are set on the FunctionApp.
	functions map[string]function
	// services contains services defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	services services
	// clients contains clients defined by the user. It is up to the
	// user to perform type assertion to handle these clients.
	clients clients
	// log provides logging for the FunctionApp. Defaults to a no-op
	// logger.
	log logger
}

// FunctionAppOption is a function that sets options to a
// FunctionApp.
type FunctionAppOption func(*FunctionApp)

// NewFunction app creates and configures a FunctionApp.
func NewFunctionApp(options ...FunctionAppOption) *FunctionApp {
	port, ok := os.LookupEnv(functionsCustomHandlerPort)
	if !ok {
		port = "8080"
	}

	router := http.NewServeMux()
	app := &FunctionApp{
		httpServer: &http.Server{
			Addr:         os.Getenv(functionsCustomHandlerPort) + ":" + port,
			Handler:      router,
			ReadTimeout:  time.Second * 30,
			WriteTimeout: time.Second * 30,
			IdleTimeout:  time.Second * 60,
		},
		functions: make(map[string]function),
		router:    router,
	}
	for _, option := range options {
		option(app)
	}
	if app.log == nil {
		app.log = noOpLogger{}
	}

	return app
}

// Start the FunctionApp.
func (a FunctionApp) Start() error {
	if len(a.functions) == 0 {
		return ErrNoFunction
	}
	for name, function := range a.functions {
		a.router.Handle("/"+name, a.handler(function))
	}

	errCh := make(chan error, 1)
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(time.Millisecond * 10):
		a.log.Info("function app started.")
		close(errCh)
	}

	_, err := a.shutdown()
	if err != nil {
		return err
	}

	return nil
}

// shutdown the function app.
func (a FunctionApp) shutdown() (os.Signal, error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	a.httpServer.SetKeepAlivesEnabled(false)
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return nil, err
	}
	return sig, nil
}

// AddFunction adds a function to the FunctionApp.
func (a *FunctionApp) AddFunction(name string, options ...FunctionOption) {
	if a.functions == nil {
		a.functions = make(map[string]function)
	}
	f := function{
		name: name,
	}
	for _, option := range options {
		option(&f)
	}
	a.functions[name] = f
}

// handler takes the provided function, creates a *Context and a trigger
// and executes the function on the route it has been configured
// with (the function name).
func (a FunctionApp) handler(fn function) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			Output:   NewOutput(WithBindings(fn.bindings...)),
			log:      a.log,
			services: a.services,
			clients:  a.clients,
		}

		if err := fn.trigger.run(r, ctx); err != nil {
			a.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(ctx.Output.JSON())
	})
}

// WithService sets the provided service to the FunctionApp. Can be
// called multiple times. If a service with the same name
// has been set it will be overwritten.
func WithService(name string, service any) FunctionAppOption {
	return func(f *FunctionApp) {
		f.services.Add(name, service)
	}
}

// WithClient sets the provided client to the FunctionApp. Can be
// called multiple times. If a client with the same name
// has been set it will be overwritten.
func WithClient(name string, client any) FunctionAppOption {
	return func(f *FunctionApp) {
		f.clients.Add(name, client)
	}
}

// WithLogger sets the provided logger to the FunctionApp.
// The logger must satisfy the logger interface.
func WithLogger(log logger) FunctionAppOption {
	return func(f *FunctionApp) {
		f.log = log
	}
}

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
