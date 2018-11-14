package util

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	LogLevel  string `split_words:"true"`
	LogFormat string `split_words:"true"`
}

func NewLogger(cfg *LogConfig) *logrus.Entry {
	log := logrus.New()
	log.Out = os.Stdout

	if cfg.LogFormat == "json" {
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}
	} else {
		log.Formatter = &logrus.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  time.RFC3339Nano,
			QuoteEmptyFields: true,
		}
	}
	if cfg.LogLevel != "" {
		if lvl, err := logrus.ParseLevel(cfg.LogLevel); err != nil {
			log.WithError(err).Warnf("Failed to parse '%s' into a log level - ignoring completely", cfg.LogLevel)
		} else {
			log.SetLevel(lvl)
		}
	}

	entry := log.WithFields(logrus.Fields{
		"sha": SHA,
		"tag": Tag,
	})

	entry.WithFields(logrus.Fields{
		"log_format": cfg.LogFormat,
		"log_level":  cfg.LogLevel,
	}).Debugf("Logger created")

	return entry
}

func RequestID(ctx context.Context) string {
	if lctx, ok := lambdacontext.FromContext(ctx); ok {
		return lctx.AwsRequestID
	}
	return "unknown"
}
