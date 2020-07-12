package main

import (
	"log"

	"github.com/elonzh/mirrorman/pkg/commands"
)

func main() {
	e := commands.NewExecutor()
	if err := e.Execute(); err != nil {
		log.Fatalln(err)
	}
}
