package config

import (
	"github.com/asehra/bezier/generator"
	"github.com/asehra/bezier/storage"
)

type Config struct {
	DB          storage.Storage
	IDGenerator generator.IDGenerator
}
