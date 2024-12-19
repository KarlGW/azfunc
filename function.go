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
	// functionsCustomHandlerPort is the environment variable that
	// contains the port for the custom handler.
	functionsCustomHandlerPort = "FUNCTIONS_CUSTOMHANDLER_PORT"
	// functionsCustomHandlerHost is the environment variable that
	// contains the host for the custom handler.
	functionsCustomHandlerHost = "FUNCTIONS_CUSTOMHANDLER_HOST"
	// functionsDisableLogging is the environment variable that
	// disables logging for the function app.
	functionsDisableLogging = "FUNCTIONS_DISABLE_LOGGING"
)

var (
	// ErrNoFunction is returned when no function has been set to the
	// FunctionApp.
	ErrNoFunction = errors.New("at least one function must be set")
)

// function is an internal structure that represents a function
// in a FunctionApp.
type function struct {
	name    string
	trigger triggerable
	outputs []outputable
}

// FunctionOption sets options to the function.
type FunctionOption func(f *function)

// WithOutput sets the provided output binding to the function.
func WithOutput(output outputable) FunctionOption {
	return func(f *function) {
		if f.outputs == nil {
			f.outputs = []outputable{output}
			return
		}
		f.outputs = append(f.outputs, output)
	}
}

// functionApp represents a Function App with its configuration
// and functions.
type functionApp struct {
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
	log    Logger
	stopCh chan os.Signal
	errCh  chan error
}

// FunctionAppOption is a function that sets options to a
// FunctionApp.
type FunctionAppOption func(*functionApp)

// NewFunctionApp creates and configures a FunctionApp.
func NewFunctionApp(options ...FunctionAppOption) *functionApp {
	port, ok := os.LookupEnv(functionsCustomHandlerPort)
	if !ok {
		port = "8080"
	}

	router := http.NewServeMux()
	app := &functionApp{
		httpServer: &http.Server{
			Addr:         os.Getenv(functionsCustomHandlerHost) + ":" + port,
			Handler:      router,
			ReadTimeout:  time.Second * 30,
			WriteTimeout: time.Second * 30,
			IdleTimeout:  time.Second * 60,
		},
		functions: make(map[string]function),
		router:    router,
		log:       setupLogger(),
		stopCh:    make(chan os.Signal),
		errCh:     make(chan error),
	}

	for _, option := range options {
		option(app)
	}

	return app
}

// New creates and configures a FunctionApp.
func New(options ...FunctionAppOption) *functionApp {
	return NewFunctionApp(options...)
}

// Start the FunctionApp.
func (a functionApp) Start() error {
	if len(a.functions) == 0 {
		return ErrNoFunction
	}
	if a.stopCh == nil {
		a.stopCh = make(chan os.Signal)
	}
	if a.errCh == nil {
		a.errCh = make(chan error)
	}

	for name, function := range a.functions {
		a.router.Handle("/"+name, a.handler(function))
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.errCh <- err
			return
		}
	}()

	go func() {
		a.stop()
	}()

	a.log.Info("Function App started.")
	for {
		select {
		case err := <-a.errCh:
			close(a.errCh)
			return err
		case sig := <-a.stopCh:
			a.log.Info("Function App stopped.", "reason", sig.String())
			close(a.stopCh)
			return nil
		}
	}
}

// stop the FunctionApp.
func (a functionApp) stop() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	a.httpServer.SetKeepAlivesEnabled(false)
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.errCh <- err
	}

	a.stopCh <- sig
}

// AddFunction adds a function to the FunctionApp.
func (a *functionApp) AddFunction(name string, options ...FunctionOption) {
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

// Add a function to the FunctionApp.
func (a *functionApp) Add(name string, options ...FunctionOption) {
	a.AddFunction(name, options...)
}

// Register a function to the FunctionApp.
func (a *functionApp) Register(name string, options ...FunctionOption) {
	a.AddFunction(name, options...)
}

// handler takes the provided function, creates a *Context and a trigger
// and executes the function on the route it has been configured
// with (the function name).
func (a functionApp) handler(fn function) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := newContext(context.Background(), func(o *contextOptions) {
			o.outputs = newOutputs(withOutputs(fn.outputs...))
			o.log = a.log
			o.services = a.services
			o.clients = a.clients
		})

		if err := fn.trigger.run(ctx, r); err != nil {
			a.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(ctx.Outputs.json())
	})
}

// WithService sets the provided service to the FunctionApp. Can be
// called multiple times. If a service with the same name
// has been set it will be overwritten.
func WithService(name string, service any) FunctionAppOption {
	return func(f *functionApp) {
		f.services.Add(name, service)
	}
}

// WithClient sets the provided client to the FunctionApp. Can be
// called multiple times. If a client with the same name
// has been set it will be overwritten.
func WithClient(name string, client any) FunctionAppOption {
	return func(f *functionApp) {
		f.clients.Add(name, client)
	}
}

// WithLogger sets the provided Logger to the FunctionApp.
// The Logger must satisfy the Logger interface.
func WithLogger(log Logger) FunctionAppOption {
	return func(f *functionApp) {
		f.log = log
	}
}

// WithDisableLogging disables logging for the FunctionApp.
func WithDisableLogging() FunctionAppOption {
	return func(f *functionApp) {
		f.log = noOpLogger{}
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

// parseBool parses a string to a boolean.
// Everything but "true" and "1" will return false.
func parseBool(s string) bool {
	return s == "true" || s == "1"
}

// setupLogger determines if logging should be disabled or not
// based on an environment variable.
func setupLogger() Logger {
	if parseBool(os.Getenv(functionsDisableLogging)) {
		return noOpLogger{}
	}
	return NewLogger()
}
