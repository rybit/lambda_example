package main

import (
	"errors"

	"github.com/rybit/lambda_example/util"
)

type configuration struct {
	util.LogConfig

	Humio      humioConfig
	TimeoutSec int `split_words:"true"`
	Test       bool
}

func (c *configuration) validate() error {
	if c.Humio.Token == "" {
		return errors.New("Must set HUMIO_TOKEN")
	}
	return nil
}
