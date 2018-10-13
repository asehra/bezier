package config

import (
	"io"

	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/storage"
)

type Config struct {
	DB                  storage.Storage
	CardIDGenerator     generator.CardIDGenerator
	MerchantIDGenerator generator.MerchantIDGenerator
	StdOut              io.Writer
	StdErr              io.Writer
}
