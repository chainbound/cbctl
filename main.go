package main

import (
	"os"
)

func main() {
	app := NewApp()

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
