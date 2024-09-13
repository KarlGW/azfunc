# service-bus-queue-trigger

> Example with a service bus queue trigger and service bus queue output binding

This is an example of running an Azure Function with a service bus queue triggered function and service bus queue output binding. It is recommended to make use of the Function [Core](https://learn.microsoft.com/en-us/azure/azure-functions/functions-run-local) tools to test and run this example application.

* [Run the function](#run-the-function)
* [Application structure](#application-structure)

## Run the function

Create the file `local.settings.json`:

```json
{
  "IsEncrypted": false,
  "Values": {
    "FUNCTIONS_WORKER_RUNTIME": "custom",
    "ServiceBusConnection": "<service-bus-connection-string>"
  }
}
```

This file ensures the environment variables listed in `Values` are in place for the function app for local execution.

Make sure the name of the built executable is the same as set in `host.json`:

```json
{
  // ...
  "customHandler": {
    "description": {
      "defaultExecutablePath": "service-bus-queue-trigger",
      // ...
    }
  }
}
```

```sh
go build && func start
```

## Application structure:

```sh
.
├── hello-sb-queue
│   └── function.json
├── README.md
├── go.mod
├── go.sum
├── host.json
├── main.go
```

This example only contains one function, `hello-sb-queue`. When using Azure Functions with custom handlers each function needs its own directory with the same name as the function containing a `function.json`.

```sh
.
├── hello-sb-queue
│   └── function.json
```

In `main.go` it is registered like so:

```go
package main

// ...

func main() {
    // ...
    app.AddFunction("hello-sb-queue", azfunc.QueueTrigger(/* ... */))
    // ...
}
```
