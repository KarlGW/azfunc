package main

import (
	"fmt"
	"os"

	"github.com/potatoattack/azfunc"
	"github.com/potatoattack/azfunc/trigger"
)

func main() {
	app := azfunc.NewFunctionApp()

	app.AddFunction("hello-timer", azfunc.TimerTrigger(func(ctx *azfunc.Context, trigger *trigger.Timer) error {
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
