package azfunc

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KarlGW/azfunc/bindings"
	"github.com/KarlGW/azfunc/triggers"
)

var (
	// ErrNoFunction is returned when no function has been set to the
	// FunctionApp.
	ErrNoFunction = errors.New("at least one function must be set")
	// ErrInvalidTrigger is returned when an invalid trigger has been
	// provided.
	ErrInvalidTrigger = errors.New("invalid trigger")
)

// logger is the interface that wraps around methods Info
// and Error.
type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

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
	port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !ok {
		port = "8080"
	}

	router := http.NewServeMux()
	app := &FunctionApp{
		httpServer: &http.Server{
			Addr:         os.Getenv("FUNCTIONS_CUSTOMHANDLER_HOST") + ":" + port,
			Handler:      router,
			ReadTimeout:  time.Second * 30,
			WriteTimeout: time.Second * 30,
			IdleTimeout:  time.Second * 60,
		},
		functions: make(map[string]function),
		router:    router,
		services:  make(services),
		clients:   make(clients),
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
			Output:   bindings.NewOutput(bindings.WithBindings(fn.bindings...)),
			log:      a.log,
			services: a.services,
			clients:  a.clients,
		}

		if fn.triggerFunc != nil {
			trigger, err := triggers.NewBase(r, fn.triggerName)
			if err != nil {
				a.log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := fn.triggerFunc(ctx, trigger); err != nil {
				a.log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if fn.httpTriggerFunc != nil {
			trigger, err := triggers.NewHTTP(r)
			if err != nil {
				a.log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := fn.httpTriggerFunc(ctx, trigger); err != nil {
				a.log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			a.log.Error(ErrInvalidTrigger.Error())
			http.Error(w, ErrInvalidTrigger.Error(), http.StatusInternalServerError)
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

// noOpLogger is a placeholder for when no logger is provided to the
// function app.
type noOpLogger struct{}

// Info together with Error satisfies the logger interface.
func (l noOpLogger) Info(msg string, args ...any) {}

// Error together with Info satisfies the logger interface.
func (l noOpLogger) Error(msg string, args ...any) {}
