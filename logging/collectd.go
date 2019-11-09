package logging

import (
	"fmt"

	"collectd.org/plugin"
)

type CollectdLogger struct {
	debug     bool
	logPrefix string
}

func NewCollectdLogger(logPrefix string) *CollectdLogger {
	return &CollectdLogger{
		logPrefix: logPrefix,
	}
}

func (l *CollectdLogger) SetLogPrefix(prefix string) {
	l.logPrefix = prefix
}

func (l *CollectdLogger) SetDebug(enable bool) {
	l.debug = enable
}

func (l *CollectdLogger) Info(msg string) error {
	return plugin.Info(l.logPrefix + msg)
}

func (l *CollectdLogger) Infof(format string, vargs ...interface{}) error {
	return plugin.Infof(l.logPrefix+format, vargs...)
}

func (l *CollectdLogger) Error(msg string) error {
	return plugin.Error(l.logPrefix + msg)
}

func (l *CollectdLogger) Errorf(format string, vargs ...interface{}) error {
	return plugin.Errorf(l.logPrefix+format, vargs...)
}

func (l *CollectdLogger) Warning(msg string) error {
	return plugin.Warning(l.logPrefix + msg)
}

func (l *CollectdLogger) Warningf(format string, vargs ...interface{}) error {
	return plugin.Warningf(l.logPrefix+format, vargs...)
}

func (l *CollectdLogger) Debug(msg string) error {
	if l.debug {
		return plugin.Infof("[DEBUG] %s", l.logPrefix+msg)
	}
	return nil
}

func (l *CollectdLogger) Debugf(format string, vargs ...interface{}) error {
	if l.debug {
		return plugin.Infof("[DEBUG] %s", fmt.Sprintf(l.logPrefix+format, vargs...))
	}
	return nil
}
