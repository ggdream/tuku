package main

import (
	"os"

	"github.com/ggdream/tuku/app"
)

var configFile = "config.yaml"

func main() {
	args := os.Args
	if len(args) > 1 {
		configFile = args[1]
	}

	tuKu, err := app.New(configFile)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	if err = tuKu.Run(); err != nil {
		println(err.Error())
	}
}
