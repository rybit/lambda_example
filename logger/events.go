package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type rawEvent struct {
	AWSLogs struct {
		Data string
	}
}

type decodedEvent struct {
	Owner string
	// LogGroup is like /aws/lambda/logging-dev
	LogGroup string
	// MessageType Data messages will use the "DATA_MESSAGE" type.
	// Sometimes CloudWatch Logs may emit Kinesis records
	// with a "CONTROL_MESSAGE" type, mainly for checking if the
	// destination is reachable.
	MessageType string
	// LogEvents are the actual lines sent
	LogEvents []logEvent

	// LogStream string // ignored
	// SubscriptionFilters []string // ignored
}

type logEvent struct {
	ID        string
	Timestamp int
	Message   string
}

// decode decodes the base64/gzip compressed json data into a usable form
func (in *rawEvent) decode() (*decodedEvent, error) {
	compressed, err := base64.StdEncoding.DecodeString(in.AWSLogs.Data)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode raw input")
	}

	gread, err := gzip.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create GZIP reader")
	}

	out, err := ioutil.ReadAll(gread)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read in gzip data")
	}
	res := new(decodedEvent)
	if err := json.Unmarshal(out, res); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshall raw data")
	}

	return res, nil
}
