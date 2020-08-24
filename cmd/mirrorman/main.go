package main

import (
	"log"

	"github.com/elonzh/mirrorman/pkg/commands"
)

var (
	// Populated by goreleaser during build
	version = "master"
	commit  = "?"
	date    = ""
)

func main() {
	e := commands.NewExecutor(version, commit, date)
	if err := e.Execute(); err != nil {
		log.Fatalln(err)
	}
}
