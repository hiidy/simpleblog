package main

import (
	"os"

	"github.com/hiidy/simpleblog/cmd/sb-apiserver/app"
)

func main() {
	command := app.NewSimpleBlogCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
