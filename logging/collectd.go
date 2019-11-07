package logging

import (
	"fmt"

	"collectd.org/plugin"
)

type CollectdLogger struct {
	debug bool
}

func NewCollectdLogger(debug bool) *CollectdLogger {
	return &CollectdLogger{
		debug: debug,
	}
}

func (l *CollectdLogger) Info(msg string) error {
	return plugin.Info(msg)
}

func (l *CollectdLogger) Infof(format string, vargs ...interface{}) error {
	return plugin.Infof(format, vargs...)
}

func (l *CollectdLogger) Error(msg string) error {
	return plugin.Error(msg)
}

func (l *CollectdLogger) Errorf(format string, vargs ...interface{}) error {
	return plugin.Errorf(format, vargs...)
}

func (l *CollectdLogger) Warning(msg string) error {
	return plugin.Warning(msg)
}

func (l *CollectdLogger) Warningf(format string, vargs ...interface{}) error {
	return plugin.Warningf(format, vargs...)
}

func (l *CollectdLogger) Debug(msg string) error {
	if l.debug {
		return plugin.Infof("[DEBUG] %s", msg)
	}
	return nil
}

func (l *CollectdLogger) Debugf(format string, vargs ...interface{}) error {
	if l.debug {
		return plugin.Infof("[DEBUG] %s", fmt.Sprintf(format, vargs...))
	}
	return nil
}
