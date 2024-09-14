package main

import (
	"fmt"
	"os"

	"github.com/KarlGW/azfunc"
	"github.com/KarlGW/azfunc/trigger"
)

func main() {
	app := azfunc.NewFunctionApp(azfunc.WithLogger(azfunc.NewLogger()))

	app.AddFunction("hello-queue", azfunc.QueueTrigger("queue", func(ctx *azfunc.Context, trigger *trigger.Queue) error {
		// Parse the incoming queue trigger body into the custom type.
		// To get the raw data of the queue message, use trigger.Data instead.
		var t test
		if err := trigger.Parse(&t); err != nil {
			return err
		}
		// Log parsed t.
		ctx.Log().Info("queue message received", "content", t)
		// Create output to queue.
		ctx.Outputs.Binding("outqueue").Write([]byte(`{"message":"message received"}`))
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
