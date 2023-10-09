package azfunc

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KarlGW/azfunc/triggers"
)

var (
	// ErrNoFunction is returned when no function has been set to the
	// FunctionApp.
	ErrNoFunction = errors.New("at least one function must be set on")
	// ErrInvalidTrigger is returned when an invalid trigger has been
	// provided.
	ErrInvalidTrigger = errors.New("invalid trigger")
)

// FunctionApp represents a Function App with its configuration
// and functions.
type FunctionApp struct {
	httpServer *http.Server
	router     *http.ServeMux
	// functions that are set on the FunctionApp.
	functions map[string]*function
	// services contains services defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	services services
	// clients contains clients defined by the user. It is up to the
	// user to perform type assertion to handle these services.
	clients clients
}

// FunctionAppOption is a function that sets options to a
// FunctionApp.
type FunctionAppOption func(*FunctionApp)

// NewFunction app creates and configures a FunctionApp.
func NewFunctionApp(options ...FunctionAppOption) *FunctionApp {
	router := http.NewServeMux()
	app := &FunctionApp{
		httpServer: &http.Server{
			Handler: router,
		},
		router:   router,
		services: make(services, 0),
		clients:  make(clients, 0),
	}
	for _, option := range options {
		option(app)
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
		// Await so that error is not returned while
		// server is starting. Add logging her.
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
	f := &function{
		name: name,
	}
	for _, option := range options {
		option(f)
	}
	a.functions[name] = f
}

// handler takes the provided *function, creates a *Context and a trigger
// and executes the function on the route it has been configured
// with (the function name).
func (a FunctionApp) handler(fn *function) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			services: a.services,
			clients:  a.clients,
			bindings: fn.bindings,
		}

		if fn.httpTriggerFunc != nil {
			trigger, err := triggers.NewHTTP(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := fn.httpTriggerFunc(trigger, ctx); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if fn.triggerFunc != nil {
			trigger, err := triggers.NewGeneric(r, fn.bindingName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := fn.triggerFunc(trigger, ctx); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, ErrInvalidTrigger.Error(), http.StatusInternalServerError)
			return
		}

		// Write OK response to the function host.
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(ctx.output.JSON())
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
