package util

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func LoadConfig(cfg interface{}) {
	if cfgFile := os.Getenv("_CONFIG_FILE"); cfgFile != "" {
		if err := godotenv.Load(cfgFile); err != nil {
			panic(errors.Wrapf(err, "Failed to load config file: %s", cfgFile))
		}
	}

	if err := envconfig.Process(os.Getenv("PREFIX"), cfg); err != nil {
		panic(errors.Wrap(err, "Failed to load initial config"))
	}
}
