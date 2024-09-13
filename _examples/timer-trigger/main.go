package main

import (
	"fmt"
	"os"

	"github.com/KarlGW/azfunc"
	"github.com/KarlGW/azfunc/triggers"
)

func main() {
	app := azfunc.NewFunctionApp(azfunc.WithLogger(azfunc.NewLogger()))

	app.AddFunction("hello-timer", azfunc.TimerTrigger(func(ctx *azfunc.Context, trigger *triggers.Timer) error {
		ctx.Log().Info("timer ran")
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
