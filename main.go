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
			DB:                     storage.NewInMemoryStore(),
			CardIDGenerator:        &generator.CardIDIncrementor{LastID: 4921000000000000},
			MerchantIDGenerator:    &generator.StringIDIncrementor{Prefix: "M", LastID: 1000},
			TransactionIDGenerator: &generator.StringIDIncrementor{Prefix: "TX", LastID: 10000},
			StdOut:                 os.Stdout,
			StdErr:                 os.Stderr,
		},
	).Run(":8000")
}
