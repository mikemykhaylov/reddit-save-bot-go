package main

import (
	"os"

	"github.com/mikemykhaylov/reddit-save-bot-go/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
