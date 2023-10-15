# azfunc

> Module to assist with Azure Functions with Custom handlers and Go

The purpose of this module is to provide functions and structures that helps with handling the incoming and outgoing requests to
and from the Azure Function host when developing Azure Functions with Custom handlers and Go.

The module provides several ways of working with Azure Functions:

* [`FunctionApp`](#function-app) provides a framework that abstract away the need for setting up an HTTP server with handlers. It handles input
and output from and to the Function host.
* [Triggers](#triggers-input-bindings) and [Output](#output-output-bindings) to just provide structures and functions for handling input and output from
the function host.

## Contents

* [Why use this module?](#why-use-this-module)
  * [HTTP only Functions](#http-only-functions)
* [Install](#install)
* [Usage](#usage)
  * [`FunctionApp`](#function-app)
    * [HTTP trigger and HTTP output binding](#http-trigger-and-http-output-binding)
    * [HTTP trigger with a HTTP and a queue output binding](#http-trigger-with-a-http-and-a-queue-output-binding)
    * [Queue trigger with queue output binding](#queue-trigger-with-queue-output-binding)
    * [Timer trigger](#timer-trigger)
  * [Triggers (input bindings)](#triggers-input-bindings)
  * [Outputs (output bindings)](#output-output-bindings)
* [Roadmap](#roadmap)

## Why use this module?

After writing several Azure Functions with Custom handlers for Go I came to realize that the process of handling the incoming request payload from the Function host a tedious task (the overall structure and some content being escaped) and I found myself rewriting this request handling time and time again, every time with different ways and end result.

The idea of this module awoke to address this and to create a uniform way to do it across my current and future projects.

### HTTP only Functions

If developing a Function that will only handle incoming and outgoing HTTP (one HTTP trigger and one HTTP output binding) this module can be skipped and `host.json` can be configured like so:

```json
{
  "version": "2.0",
  // ...
  // ...
  "customHandler": {
    "description": {
      "defaultExecutablePath": ""
    },
    "enableForwardingHttpRequest": true // This setting.
  }
  // ...
  // ...
}
```

This enables the function to be developed like any HTTP server/handler.
**Note**: This is not suitable if a non HTTP output binding, or multiple output bindings are used.

## Install

**Prerequisites**

* Go 18

```sh
go get github.com/KarlGW/azfunc
```

## Usage

As a user of this module you can decide to either use the framework provided to work with the input and output from the Function host,
or to use the provided functions and structures to to parse the incoming and outgoing requests.

### `FunctionApp`

The framework takes care of setting up the server and handlers, and you as a user register functions and bindings. Each function must have a corresponding `function.json` with binding specifications, in a folder named after the function. For more information about bindings, check out the [documentation](https://learn.microsoft.com/en-us/azure/azure-functions/functions-bindings-expressions-patterns).

The triggers and bindings registered to
the `FunctionApp` must match with the names of the bindings in this file, the exception being with HTTP triggers and bindings, `req` and `res` where this is handled without the names for convenience.

Creating this structure can be done with ease with the help of the Function [Core tools](https://learn.microsoft.com/en-us/azure/azure-functions/functions-run-local). In addition to scaffolding functions ot can be used to run and test your functions locally.


#### HTTP trigger and HTTP output binding

*Create `helloHTTP/function.json` with a HTTP trigger and HTTP output binding*

```json
{
  "bindings": [
    {
      "authLevel": "function",
      "type": "httpTrigger",
      "direction": "in",
      "name": "req",
      "methods": [
        "get",
        "post"
      ]
    },
    {
      "type": "http",
      "direction": "out",
      "name": "res"
    },
  ]
}
```

```go
package main

import (
	"github.com/KarlGW/azfunc"
	"github.com/KarlGW/azfunc/triggers"
)

func main() {
    app := azfunc.NewFunctionApp()

	app.AddFunction("helloHTTP", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger triggers.HTTP) error {
        // Parse the incoming trigger body into the custom type.
        // To get the raw data of the body, use trigger.Data instead.
		var t test
		if err := trigger.Parse(&t); err != nil {
            // Send response back to caller.
			ctx.Output.HTTP().WriteHeader(http.StatusBadRequest)
            // To imply an application error, return the error here
            // instead of nil.
			return nil
		}
        // Do something with t.

        // Create the response.
		ctx.Output.HTTP().WriteHeader(http.StatusOK)
		ctx.Output.HTTP().Header().Add("Content-Type", "application/json")
		ctx.Output.HTTP().Write([]byte(`{"message":"received"}`))
		return nil
	}))

    if err := app.Start(); err != nil {
        // Handle error.
    }
}
```

#### HTTP trigger with a HTTP and a queue output binding

*Create `helloHTTPQueue/function.json` with a HTTP trigger and a HTTP and a queue output binding*:

```json
{
  "bindings": [
    {
      "authLevel": "function",
      "type": "httpTrigger",
      "direction": "in",
      "name": "req",
      "methods": [
        "get",
        "post"
      ]
    },
    {
      "type": "http",
      "direction": "out",
      "name": "res"
    },
    // Add queue output binding.
    {
      "name": "queue",
      "type": "queue",
      "direction": "out",
      "queueName": "out",
      "connection": "AzureWebJobsStorage"
    }
  ]
}
```

```go
package main

import (
	"github.com/KarlGW/azfunc"
	"github.com/KarlGW/azfunc/triggers"
)

func main() {
    app := azfunc.NewFunctionApp()

    app.AddFunction("helloHTTP", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger triggers.HTTP) error {
        // Omitted for example.
    }))

	app.AddFunction("helloHTTPQueue", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger triggers.HTTP) error {
        // Parse the incoming trigger body into the custom type.
        // To get the raw data of the body, use trigger.Data instead.
		var t test
		if err := trigger.Parse(&t); err != nil {
            // Send response back to caller.
			ctx.Output.HTTP().WriteHeader(http.StatusBadRequest)
            // To imply an application error, return the error here
            // instead of nil.
			return nil
		}
        // Do something with t.

        // Create the response to the caller.
		ctx.Output.HTTP().WriteHeader(http.StatusOK)
		ctx.Output.HTTP().Header().Add("Content-Type", "application/json")
		ctx.Output.HTTP().Write([]byte(`{"message":"received"}`))
        // Send a message to the queu. Name should be the same name as specified
        // in function.json for the binding.
        ctx.Output.Binding("queue").Write([]byte(`{"message":"hello queue"}`))
		return nil
	}))

    if err := app.Start(); err != nil {
        // Handle error.
    }
}
```

#### Queue trigger with queue output binding

*Create `helloQueue/function.json` with a queue trigger and a queue output binding*:

```json
{
  "bindings": [
    {
      "name": "queue",
      "type": "queueTrigger",
      "direction": "in",
      "queueName": "items",
      "connection": "AzureWebJobsStorage"
    },
    {
      "name": "outqueue",
      "type": "queue",
      "direction": "out",
      "queueName": "out",
      "connection": "AzureWebJobsStorage"
    }
  ]
}
```

```go
package main

import (
	"github.com/KarlGW/azfunc"
	"github.com/KarlGW/azfunc/triggers"
)

func main() {
    app := azfunc.NewFunctionApp()

    app.AddFunction("helloHTTP", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger triggers.HTTP) error {
        // Omitted for example.
    }))

	app.AddFunction("helloHTTPQueue", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger triggers.HTTP) error {
        // Omitted for example.
    }))

	app.AddFunction("helloQueue", azfunc.QueueTrigger("queue", func(ctx *azfunc.Context, trigger triggers.Queue) error {
		var t test
		if err := trigger.Parse(&t); err != nil {
			ctx.Output.ReturnValue = 1
			return nil
		}
		fmt.Println("message:", t)

		ctx.Output.Binding("outqueue").Write([]byte(`{"message":"from-queue"}`))
		return nil
	}))

    if err := app.Start(); err != nil {
        // Handle error.
    }
}
```

#### Timer trigger

*Create `helloTimer/function.json` with a timer trigger*:

```json
{
  "bindings": [
    {
      "name": "timer",
      "type": "timerTrigger",
      "direction": "in",
      "schedule": "0 */2 * * * *"
    }
  ]
}
```

```go
package main

import (
	"github.com/KarlGW/azfunc"
	"github.com/KarlGW/azfunc/triggers"
)

func main() {
    app := azfunc.NewFunctionApp()

    app.AddFunction("helloHTTP", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger triggers.HTTP) error {
        // Omitted for example.
    }))

	app.AddFunction("helloHTTPQueue", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger triggers.HTTP) error {
        // Omitted for example.
    }))

	app.AddFunction("helloQueue", azfunc.QueueTrigger("queue", func(ctx *azfunc.Context, trigger triggers.Queue) error {
        // Omitted for example.
	}))

	app.AddFunction("helloTimer", azfunc.TimeTrigger("timer", func(ctx *azfunc.Context, b triggers.Base) error {
		// Do something.
		return nil
	}))

    if err := app.Start(); err != nil {
        // Handle error.
    }
}
```

### Triggers (input bindings)

An Azure Function can have one and only one trigger. This module provides a couple of ways to handle the incoming trigger requests, suitable for different purposes.

The following triggers are supported:

* `HTTP`
* `Base`
* `Queue` (alias for `Base`, to provide clarity and intention of the trigger)

Custom defined trigger types can be used as long as they satisfy the `Triggerable` interface:

```go
type Triggerable interface {
	// Data returns the raw data of the trigger.
	Data() data.Raw
	// Parse the raw data of the trigger into the provided value.
	Parse(v any) error
}
```

**`Parse[t Triggerable](r *http.Request, v any, options ...Options) error`**

The most straight forward way to handle an incoming request and
parse it's data to a struct.

```go
package main

import (
    "github.com/KarlGW/azfunc/triggers"
)

// Request handler for Function "helloHTTP" that handles
// HTTP trigger requests. This will parse the body of
// the incoming HTTP request into "customType" (assuming)
// is is valid JSON.
func helloHTTPHandler(r *http.Request, w http.ResponseWriter) {
    var t customType
    if err := triggers.Parse[trigger.HTTP](r, &t); err != nil {
        // Handle error.
    }
    // ...
    // ...
}

// Request handler for Function "helloQueue" that handles
// Queue trigger requests. This will parse the data of
// the incoming Queue request into "customType" (assuming)
// is is valid JSON.
func helloQueueHandler(r *http.Request, w http.ResponseWriter) {
    var t customType
    if err := triggers.Parse[triggers.Queue][r, &t, triggers.WithName("<trigger/binding-name>")]; err != nil {
        // Handle error.
    }
    // ...
    // ...
}

// Simple http server for handling function requests.
func main() {
    port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
    if !ok {
        port = "8080"
    }

    r := http.NewServeMux()
    r.HandleFunc("/helloHTTP", helloHTTPHandler)
    r.HandleFunc("/helloQueue", helloQueueHandler)

    fmt.Println("Listening on port:", port)
    log.Fatalln(http.ListenAndServe(":"+port, r))
}
```

**`New[T Triggerable](r *http.Request, options ...Options) (Trigger[T], error)`**

Creates a new `Trigger[T]` from the incoming request that gives access to all
the data contained in the trigger.

```go
package main

import (
    "github.com/KarlGW/azfunc/triggers"
)

// Request handler for Function "helloHTTP" that handles
// HTTP trigger requests. This will parse the body of
// the incoming HTTP request and put it into "trigger".
// If it's desired to get the underlying type and handle
// the fields on the HTTP trigger use trigger.Trigger().
func helloHTTPHandler(r *http.Request, w http.ResponseWriter) {
    trigger, err := triggers.New[triggers.HTTP](req)
    if err != nil {
        // Handle error.
    }
    // Parse the body into a struct.
    var t customType
    if err := trigger.Parse(&t); err != nil {
        // Handle error.
    }
    // ...
    // ...
}

// Request handler for Function "helloQueue" that handles
// Queue trigger requests. This will parse the data of
// the incoming Queue request and put it into "trigger".
// If it's desired to the underlying type and handle the fields
// on the Queue trigger, use trigger.Trigger().
func helloQueueHandler(r *http.Request, w http.ResponseWriter) {
    trigger, err := triggers.New[triggers.Queue](req)
    if err != nil {
        // Handle error.
    }
    // Parse the body into a struct.
    var t customType
    if err := trigger.Parse(&t); err != nil {
        // Handle error.
    }
    // ...
    // ...
}

// Simple http server for handling function requests.
func main() {
    port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
    if !ok {
        port = "8080"
    }

    r := http.NewServeMux()
    r.HandleFunc("/helloHTTP", helloHTTPHandler)
    r.HandleFunc("/helloQueue", helloQueueHandler)

    fmt.Println("Listening on port:", port)
    log.Fatalln(http.ListenAndServe(":"+port, r))
}
```

**`NewRequest(r *http.Request) (*http.Request, error)`**

Creates a new `*http.Request` based on the payload from the incoming request (must be an HTTP Trigger, `HTTP`). This provides
a means to handle the incoming HTTP trigger payload as any other `*http.Request`. Suitable for middlewares, or when it is just required to be
handled this way.

```go
package main

import (
    "github.com/KarlGW/azfunc/triggers"
)

// Request handler for Function "helloHTTP" that handles
// HTTP trigger requests. This will parse the body of
// the incoming HTTP request into "customType" (assuming)
// is is valid JSON.
func helloHTTPHandler(r *http.Request, w http.ResponseWriter) {
    // r will contain the request details from the Function host
    // payload, effectively "passing" on the original request.
    r, err := triggers.NewRequest(r)
    if err != nil {
        // Handle error.
    }
    // ...
    // ...
}

// Simple http server for handling function requests.
func main() {
    port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
    if !ok {
        port = "8080"
    }

    r := http.NewServeMux()
    r.HandleFunc("/helloHTTP", helloHTTPHandler)

    fmt.Println("Listening on port:", port)
    log.Fatalln(http.ListenAndServe(":"+port, r))
}
```

### Output (output bindings)

An Azure Function can have multiple output bindings. This module
provides the `Output` struct to handle all outputs to the Function host.

The following output bindings are supported:

* `HTTP`
* `Base`
* `Queue` (alias for `Base`, to provide clarity and intention of the binding)

Custom defined binding types can be used as long as they satisfy the
`Bindable` interface:

```go
type Bindable interface {
	// Name returns the name of the binding.
	Name() string
	// Write to the binding.
	Write([]byte) (int, error)
}
```

#### Example

```go
package main

import (
    "github.com/KarlGW/azfunc/bindings"
)

// Request handler for Function "helloHTTP" that handles
// HTTP trigger requests. This will parse the body of
// the incoming HTTP request into "customType" (assuming)
// is is valid JSON. Together with an HTTP and Queue output binding.
func helloHTTPHandler(r *http.Request, w http.ResponseWriter) {
    // Parse the incoming request.

    // Create output.
    output := bindings.NewOutput(
        bindings.WithBindings(
            bindings.NewHTTP(),
            bindings.NewQueue("<queue-binding-name>")
        )
    )

    // Write to the HTTP output binding.
    output.HTTP().WriteHeader(http.StatusOK)
    output.HTTP().Write(`{"message":"hello, world"}`)

    // Write to the queue output binding.
    output.Binding("<queue-binding-name>").Write(`{"message":"hello, world"}`)

    // All custom handlers regardless of output binding type
    // must set Content-Type: application/json to the response
    // to the Function host, followed by the data.
    w.Header.Set("Content-Type", "application/json")
    w.Write(output.JSON())
}
```

Bindings can also be added after the output is created:

```go
package main

import (
    "github.com/KarlGW/azfunc/bindings"
)

// Request handler for Function "helloHTTP" that handles
// HTTP trigger requests. This will parse the body of
// the incoming HTTP request into "customType" (assuming)
// is is valid JSON. Together with an HTTP and Queue output binding.
func helloHTTPHandler(r *http.Request, w http.ResponseWriter) {
    // Parse the incoming request.

    // Create output.
    output := bindings.NewOutput()
    output.AddBindings(
        bindings.NewHTTP(),
        bindings.NewQueue("<queue-binding-name>")
    )

    // All custom handlers regardless of output binding type
    // must set Content-Type: application/json to the response
    // to the Function host, followed by the data.
    w.Header.Set("Content-Type", "application/json")
    w.Write(output.JSON())
}
```

It is possible to create a binding and use it at the same times:

```go
package main

import (
    "github.com/KarlGW/azfunc/bindings"
)

// Request handler for Function "helloHTTP" that handles
// HTTP trigger requests. This will parse the body of
// the incoming HTTP request into "customType" (assuming)
// is is valid JSON. Together with an HTTP and Queue output binding.
func helloHTTPHandler(r *http.Request, w http.ResponseWriter) {
    // Parse the incoming request.

    // Create output.
    output := bindings.NewOutput()
    // Create and write to the HTTP binding.
    output.HTTP().Write(`{"message":"hello, world"}`)
    // Create and write to the queue binding.
    output.Binding("<queue-binding-name>").Write(`{"message":"hello, world"}`)

    // All custom handlers regardless of output binding type
    // must set Content-Type: application/json to the response
    // to the Function host, followed by the data.
    w.Header.Set("Content-Type", "application/json")
    w.Write(output.JSON())
}
```