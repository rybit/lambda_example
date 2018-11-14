package main

import (
	"errors"

	"github.com/rybit/lambda_example/util"
)

type configuration struct {
	util.LogConfig

	Humio      humioConfig
	TimeoutSec int `split_words:"true"`
}

func (c *configuration) validate() error {
	if c.Humio.Token == "" {
		return errors.New("Must set HUMIO_TOKEN")
	}
	if c.Humio.Repository == "" {
		return errors.New("Must set HUMIO_REPOSITORY")
	}
	return nil
}
