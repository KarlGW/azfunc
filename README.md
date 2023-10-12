# azfunc

> Module to assist working with Azure Functions Custom handlers with Go

The purpose of this module is to provide functions and structures that helps with handling the incoming and outgoing requests to and from the Azure Function host when developing Azure Functions with Custom handlers and Go.

The module provides two ways of working with Functions, triggers and output bindings:

* [`FunctionApp`](#functions) which provides a framework to abstract away the need for setting up an HTTP server with handlers. This makes use of the modules various Triggers and output bindings.
* [Triggers](#triggers-input-bindings) and [Output](#output-output-bindings) to just provide structurs and functions for handling input and output from
the function host.

## Contents

* [Why use this module?](#why-use-this-module)
  * [HTTP only Functions](#http-only-functions)
* [Install](#install)
* [Usage](#usage)
  * [Functions](#functions)
  * [Triggers (input bindings)](#triggers-input-bindings)
  * [Outputs (output bindings)](#output-output-bindings)
* [Roadmap](#roadmap)

## Why use this module?

After writing several Azure Functions with Custom handlers for Go I came to realize that the process of handling the incoming request payload from the Function host a tedious task (the overall structure and some content being escaped) and I found myself rewriting this request handling time and time again, every time with different ways and end result.

For examples on how incoming and outgoing requests might look, look at the bottom of [triggers_test.go](triggers(triggers_test.go) and [bindings_test.go](bindings/bindings_test.go) respectively.

The thought of this module awoke to address this and to create a uniform
way to do it across my different and future projects.

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

### Functions



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
    var customType struct
    if err := triggers.Parse[trigger.Queue][r, &customType, triggers.WithName("<trigger/binding-name>")]; err != nil {
        // Handle error.
    }
    // Parse the data into a struct.
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
            bindings.NewHTTP(http.StatusOk, []byte(`{"message":"hello world"}`)),
            bindings.NewQueue("<queue-binding-name>", []byte(`{"message":"Hello queue"}`))
        )
    )

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
    // ...
    // ...

    // Can be added in their on subsequent calls.
    output.AddBindings(
        bindings.NewHTTP(http.StatusOk, []byte(`{"message":"hello world"}`)),
        bindings.NewQueue("<queue-binding-name>", []byte(`{"message":"Hello queue"}`))
    )

    // All custom handlers regardless of output binding type
    // must set Content-Type: application/json to the response
    // to the Function host, followed by the data.
    w.Header.Set("Content-Type", "application/json")
    w.Write(output.JSON())
}
```
