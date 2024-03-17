# azfunc

[![Go Reference](https://pkg.go.dev/badge/github.com/KarlGW/azfunc.svg)](https://pkg.go.dev/github.com/KarlGW/azfunc)

> Module to assist with Azure Functions with Custom handlers and Go

The purpose of this module is to provide functions and structures that helps with handling the incoming and outgoing requests to
and from the Azure Function host when developing Azure Functions with Custom handlers and Go.


## Contents

* [Why use this module?](#why-use-this-module)
* [Install](#install)
* [Usage](#usage)
  * [Concepts](#concepts)
    * [Triggers (input bindings)](#triggers-input-bindings)
    * [Output (output bindings)](#output-output-bindings)
    * [Context](#context)
    * [Error handling](#error-handling)
  * [HTTP trigger and HTTP output binding](#http-trigger-and-http-output-binding)


## Why use this module?

After writing several Azure Functions with Custom handlers for Go I came to realize that the process of handling the incoming request payload from the Function host a tedious task (the overall structure and some content being escaped) and I found myself rewriting this request handling time and time again, every time with different ways and end result.

The idea of this module awoke to address this and to create a uniform way to do it across my current and future projects.

## Install

**Prerequisites**

* Go 1.18

```sh
go get github.com/KarlGW/azfunc
```

## Usage

The framework takes care of setting up the server and handlers, and you as a user register functions and bindings. Each function must have a corresponding `function.json` with binding specifications, in a folder named after the function. For more information about bindings, check out the [documentation](https://learn.microsoft.com/en-us/azure/azure-functions/functions-bindings-expressions-patterns).

The triggers and bindings registered to
the `FunctionApp` must match with the names of the bindings in this file, the exception being with HTTP triggers and bindings, `req` and `res` where this is handled without the names for convenience.

Creating this application structure can be done with ease with the help of the Function [Core tools](https://learn.microsoft.com/en-us/azure/azure-functions/functions-run-local). In addition to scaffolding functions it can be used to run and test your functions locally.

An example on how to create a function with a HTTP trigger and a HTTP output binding (response) is provided further [below](#http-trigger-and-http-output-binding).


### Concepts

When working with the `FunctionApp` there are some concepts to understand and work with. The `FunctionApp` represents the entire Function App, and it is to this structure the functions (with their triggers and output bindings) that should be run are registered to. Each function that is registered contains a `*azfunc.Context` and a trigger (`*triggers.HTTP`, `*triggers.Timer`, `*triggers.Queue` and `*triggers.Base`).

The triggers is the triggering event and the data it contains, and the context contains output bindings (and writing to them), output error and logging.


#### Triggers (input bindings)

**[HTTP trigger](https://pkg.go.dev/github.com/KarlGW/azfunc/triggers#HTTP)**

Triggered by an incoming HTTP event. The trigger contains the HTTP data (headers, url, query, params and body).

```go
func(ctx *azfunc.Context, trigger *triggers.HTTP)
```

**[Timer trigger](https://pkg.go.dev/github.com/KarlGW/azfunc/triggers#Timer)**

Triggered by a schedule. The trigger contains the timer data (next and last run etc).

```go
func(ctx *azfunc.Context, trigger *triggers.Timer)
```

**[Queue trigger](https://pkg.go.dev/github.com/KarlGW/azfunc/triggers#Queue)**

Triggered by a message to an Azure Queue Storage queue.

```go
func(ctx *azfunc.Context, trigger *triggers.Queue)
```

**[Service Bus trigger](https://pkg.go.dev/github.com/KarlGW/azfunc/triggers#ServiceBus)**

Triggered by a message to an Azure Service Bus queue or topic subscription.

```go
func(ctx *azfunc.Context, trigger *triggers.ServiceBus)
```

**[Event Grid trigger](https://pkg.go.dev/github.com/KarlGW/azfunc/triggers#EventGrid)**

Triggered by an event to an Azure Event Grid topic subscription.

```go
func(ctx *azfunc.Context, trigger *triggers.EventGrid)
```

**[Base trigger](https://pkg.go.dev/github.com/KarlGW/azfunc/triggers#Base)**

Base trigger is a base that can be used for all not yet supported triggers. The data it contains
needs to be parsed into a `struct` matching the expected incoming payload.

```go
func(ctx *azfunc.Context, trigger *triggers.Base)
```


#### Output (output bindings)

**[HTTP binding](https://pkg.go.dev/github.com/KarlGW/azfunc/bindings#HTTP)**

Writes an HTTP response back to the caller (only works together with an **HTTP trigger**).

**[Queue binding](https://pkg.go.dev/github.com/KarlGW/azfunc/bindings#Queue)**

Writes a message to a queue in Azure Queue Storage.

**[Service Bus binding](https://pkg.go.dev/github.com/KarlGW/azfunc/bindings#Queue)**

Writes a message to a queue or topic subscription in Azure Service Bus.

**[Base binding](https://pkg.go.dev/github.com/KarlGW/azfunc/bindings#Base)**

Base binding is a base that can be used for all not yet supported bindings.

#### [Context](https://pkg.go.dev/github.com/KarlGW/azfunc#Context)

The context is the Function context, named so due to it being called so in the Azure Function implementation of other languages (foremost the old way of handling JavaScript/Node.js functions).

Assuming the `*azfunc.Context` is bound to the name `ctx`:

* `ctx.Log()`:
    * `ctx.Log().Info()` for info level logs.
    * `ctx.Log.Error()` for error level logs.
* `ctx.Binding("<binding-name>")` - Provides access to the binding by name. If the binding it hasn't been provided together with the function at registration, it will created as a `*bindings.Base` (will work as long as a binding with that same name is defined in the functions `function.json`).
* `ctx.Err()` - Get the error set to the context.
* `ctx.SetError(err)` - Set an error to the context. This will signal to the underlying host that an error has occured and the execution will count as a failure. Use this followed by a `return` in the provided function to represent an "unrecoverable" error and exit early.


#### Error handling

The functions provided to the `FunctionApp` has no return value, and as such does not return errors. These needs be handled by either setting them to the `*azfunc.Context` and passed on to be handled by the function executor in the `FunctionApp`, or handling them in some way within the function code.

As an example: A function is triggered and run, and encounters an error for one of it's calls. This error is deemed to be fatal and the function cannot carry on further. Then setting the error to the context and returning is the preferred way, like so:

```go
func run(ctx *azfunc.Context, trigger *triggers.Queue) {
    if err := someFunc(); err != nil {
        // This error is fatal, and must be signaled to the function host.
        ctx.SetError(err)
        return
    }
}
```

An example of when an error is not regarded as fatal, but not application breaking (like a malformed HTTP request), it can be handled like so:

```go
func run(ctx *azfunc.Context, trigger *triggers.HTTP) {
    var incoming IncomingRequest
    if err := trigge.Parse(&incoming); err != nil {
        // The incoming request body did not match the expected one,
        // this is equivalent of a HTTP 400.
        ctx.HTTP().WriteHeader(http.StatusBadRequest)
        return
    }
}
```



### HTTP trigger and HTTP output binding

*Create `hello-http/function.json` with a HTTP trigger and HTTP output binding*

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

    app.AddFunction("hello-http", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger *triggers.HTTP) {
        // Parse the incoming trigger body into the custom type.
        // To get the raw data of the body, use trigger.Data instead.
        var t test
        if err := trigger.Parse(&t); err != nil {
            // Send response back to caller.
            ctx.Output.HTTP().WriteHeader(http.StatusBadRequest)
            return
        }
        // Do something with t.
        // Create the response.
        ctx.Output.HTTP().WriteHeader(http.StatusOK)
        ctx.Output.HTTP().Header().Add("Content-Type", "application/json")
        ctx.Output.HTTP().Write([]byte(`{"message":"received"}`))
    }))

    if err := app.Start(); err != nil {
        // Handle error.
    }
}

type test struct {
    Message string `json:"message"`
}
```

More examples with different triggers and output can be found [here](./_examples/).

