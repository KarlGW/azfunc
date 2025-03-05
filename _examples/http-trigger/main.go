package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/potatoattack/azfunc"
	"github.com/potatoattack/azfunc/trigger"
)

func main() {
	app := azfunc.NewFunctionApp()

	app.AddFunction("hello-http", azfunc.HTTPTrigger(func(ctx *azfunc.Context, trigger *trigger.HTTP) error {
		// Parse the incoming trigger body into the custom type.
		// To get the raw data of the body, use trigger.Data instead.
		var t test
		if err := trigger.Parse(&t); err != nil {
			// Send HTTP response back to the caller if parsing fails
			// and exit the function.
			ctx.Outputs.HTTP().WriteHeader(http.StatusBadRequest)
			return nil
		}
		// Log parsed t.
		ctx.Log().Info("request received", "body", t)
		// Create the HTTP response.
		ctx.Outputs.HTTP().WriteHeader(http.StatusOK)
		ctx.Outputs.HTTP().Header().Add("Content-Type", "application/json")
		ctx.Outputs.HTTP().Write([]byte(`{"message":"request received"}`))
		return nil
	}))

	if err := app.Start(); err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

type test struct {
	Message string `json:"message"`
}
