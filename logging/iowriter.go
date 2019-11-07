package logging

import (
	"fmt"
	"io"
	"sync"
)

type IOWriterLogger struct {
	Destination io.Writer
	debug       bool
	lock        sync.Mutex
}

func NewIOWriterLogger(destination io.Writer, debug bool) *IOWriterLogger {
	return &IOWriterLogger{
		Destination: destination,
		debug:       debug,
	}
}

func (l *IOWriterLogger) Info(msg string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	fmt.Fprintf(l.Destination, "[INFO] %s\n", msg)
	return nil
}

func (l *IOWriterLogger) Infof(format string, vargs ...interface{}) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	fmt.Fprintf(l.Destination, "[INFO] %s\n", fmt.Sprintf(format, vargs...))
	return nil
}

func (l *IOWriterLogger) Error(msg string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	fmt.Fprintf(l.Destination, "[ERROR] %s\n", msg)
	return nil
}

func (l *IOWriterLogger) Errorf(format string, vargs ...interface{}) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	fmt.Fprintf(l.Destination, "[ERROR] %s\n", fmt.Sprintf(format, vargs...))
	return nil
}

func (l *IOWriterLogger) Warning(msg string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	fmt.Fprintf(l.Destination, "[WARNING] %s\n", msg)
	return nil
}

func (l *IOWriterLogger) Warningf(format string, vargs ...interface{}) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	fmt.Fprintf(l.Destination, "[WARNING] %s\n", fmt.Sprintf(format, vargs...))
	return nil
}

func (l *IOWriterLogger) Debug(msg string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.debug {
		fmt.Fprintf(l.Destination, "[DEBUG] %s\n", msg)
	}
	return nil
}

func (l *IOWriterLogger) Debugf(format string, vargs ...interface{}) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.debug {
		fmt.Fprintf(l.Destination, "[DEBUG] %s\n", fmt.Sprintf(format, vargs...))
	}
	return nil
}
