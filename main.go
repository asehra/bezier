package main

import (
	"os"

	"github.com/asehra/bezier/config"
	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/storage"
	"github.com/asehra/bezier/webserver"
)

func main() {
	webserver.Create(
		config.Config{
			DB:          storage.NewInMemoryStore(),
			IDGenerator: &generator.CardID{LastID: 4921000000000000},
			StdOut:      os.Stdout,
			StdErr:      os.Stderr,
		},
	).Run(":8000")
}
