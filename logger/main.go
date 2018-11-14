package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/rybit/lambda_example/util"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

var config configuration
var rootLogger *logrus.Entry
var client *http.Client

func main() {
	util.LoadConfig(&config)
	rootLogger = util.NewLogger(&config.LogConfig)
	if err := config.validate(); err != nil {
		rootLogger.WithError(err).Fatal("Invalid configuration")
	}

	client = new(http.Client)
	if config.TimeoutSec > 0 {
		client.Timeout = time.Second * time.Duration(config.TimeoutSec)
	}

	rootLogger.Debug("Startup completed")
	lambda.Start(handleEvent)
}

// handleEvent will decode the payload and send it to humio. Errors will only be returned if we could recover on retry
func handleEvent(ctx context.Context, input rawEvent) error {
	log := rootLogger.WithField("aws_id", util.RequestID(ctx))
	decoded, err := input.decode()
	if err != nil {
		log.WithError(err).Error("Failed to decode message")
		return nil // swallow because we don't want the msg again
	}
	log.Debug("Successfully decoded message")

	out := &humioMsg{
		Tags: map[string]interface{}{
			"aws_account_id": decoded.Owner,
			"message_type":   decoded.MessageType,
		},
	}

	if config.Humio.Parser != "" {
		out.Type = config.Humio.Parser
		out.Tags["log_group"] = decoded.LogGroup
	} else {
		out.Type = strings.Replace(decoded.LogGroup, "/", "_", -1)
	}

	for _, le := range decoded.LogEvents {
		out.Messages = append(out.Messages, le.Message)
	}

	code, err := send(log, out)
	if err != nil {
		log.WithError(err).Error("Failed to post data to humio")
		return err
	}
	log.WithField("status_code", code).Info("Finished sending lines to humio")
	return nil
}
