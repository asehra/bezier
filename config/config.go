package config

import (
	"io"

	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/storage"
)

type Config struct {
	DB                     storage.Storage
	CardIDGenerator        generator.CardIDGenerator
	MerchantIDGenerator    generator.StringIDGenerator
	TransactionIDGenerator generator.StringIDGenerator
	StdOut                 io.Writer
	StdErr                 io.Writer
}
