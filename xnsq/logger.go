package xnsq

import (
	"strings"

	nsq "github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
)

var (
	nsqDebugLevel = nsq.LogLevelDebug.String()
	nsqInfoLevel  = nsq.LogLevelInfo.String()
	nsqWarnLevel  = nsq.LogLevelWarning.String()
	nsqErrLevel   = nsq.LogLevelError.String()
)

// LogrusLogger is an adaptor between the weird go-nsq Logger and our
// standard logrus logger.
type LogrusLogger struct{}

// NewLogrusLogger returns a new LogrusLogger and the current log level.
// This is a format to easily plug into nsq.SetLogger.
func NewLogrusLogger() (logger LogrusLogger, level nsq.LogLevel) {
	return NewLogrusLoggerAtLevel(log.GetLevel())
}

// NewLogrusLoggerAtLevel returns a new LogrusLogger with the provided log level mapped to nsq.LogLevel for easily plugging into nsq.SetLogger.
func NewLogrusLoggerAtLevel(l log.Level) (logger LogrusLogger, level nsq.LogLevel) {
	logger = LogrusLogger{}
	level = nsq.LogLevelWarning
	switch l {
	case log.DebugLevel:
		level = nsq.LogLevelDebug
	case log.InfoLevel:
		level = nsq.LogLevelInfo
	case log.WarnLevel:
		level = nsq.LogLevelWarning
	case log.ErrorLevel:
		level = nsq.LogLevelError
	}
	return
}

// Output implements stdlib log.Logger.Output using logrus
// Decodes the go-nsq log messages to figure out the log level
func (n LogrusLogger) Output(_ int, s string) error {
	if len(s) > 3 {
		msg := strings.TrimSpace(s[3:])
		switch s[:3] {
		case nsqDebugLevel:
			log.Debug(msg)
		case nsqInfoLevel:
			log.Info(msg)
		case nsqWarnLevel:
			log.Warn(msg)
		case nsqErrLevel:
			log.Error(msg)
		default:
			log.Info(msg)
		}
	}
	return nil
}
