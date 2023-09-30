# azfunc

> Module to assist working with Azure Functions Custom handlers with Go

The purpose of this module is to provide functions and structures that helps with handling the incoming and outgoing requests to and from the Azure Function host when developing Azure Functions with Custom handlers and Go.

* [Intial advice](#initial-advice)

## Initial advice

If you are developing a Function that will only handle incoming and outgoing HTTP (HTTP trigger and HTT output binding) you can skip this module alltogether and just configure `host.json` like so:

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

That enbables you to develop it like any HTTP server/handler.
