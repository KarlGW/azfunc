# http-trigger-queue-output

> Example with a HTTP trigger and a queue output binding

This is an example of running an Azure Function with a HTTP triggered function with a queue output binding. It is recommended to make use of the Function [Core](https://learn.microsoft.com/en-us/azure/azure-functions/functions-run-local) tools together with [azurite](https://learn.microsoft.com/en-us/azure/storage/common/storage-use-azurite) to test and run this example application.

* [Run the function](#run-the-function)
* [Application structure](#application-structure)

## Run the function

Create the file `local.settings.json`:

```json
{
  "IsEncrypted": false,
  "Values": {
    "FUNCTIONS_WORKER_RUNTIME": "custom",
    "AzureWebJobsStorage": "UseDevelopmentStorage=true"
  }
}
```

This file ensures the environment variables listed in `Values` are in place for the function app for local execution.
If using a real storage account instead of `azurite`, add the connection string to the property: `AzureWebJobsStorage`.

Make sure the name of the built executable is the same as set in `host.json`:

```json
{
  // ...
  "customHandler": {
    "description": {
      "defaultExecutablePath": "http-trigger-queue-output",
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
├── hello-http-queue
│   └── function.json
├── README.md
├── go.mod
├── go.sum
├── host.json
├── main.go
```

This example only contains one function, `hello-http-queue`. When using Azure Functions with custom handlers each function needs its own directory with the same name as the function containing a `function.json`.

```sh
.
├── hello-http-queue
│   └── function.json
```

In `main.go` it is registered like so:

```go
package main

// ...

func main() {
    // ...
    app.AddFunction("hello-http-queue", azfunc.HTTPTrigger(/* ... */))
    // ...
}
```
