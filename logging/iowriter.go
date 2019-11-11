// Copyright (c) IN2P3 Computing Centre, IN2P3, CNRS
// 
// Author(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2019
// 
// This software is governed by the CeCILL-C license under French law and
// abiding by the rules of distribution of free software.  You can  use, 
// modify and/ or redistribute the software under the terms of the CeCILL-C
// license as circulated by CEA, CNRS and INRIA at the following URL
// "http://www.cecill.info". 
// 
// As a counterpart to the access to the source code and  rights to copy,
// modify and redistribute granted by the license, users are provided only
// with a limited warranty  and the software's author,  the holder of the
// economic rights,  and the successive licensors  have only  limited
// liability. 
// 
// In this respect, the user's attention is drawn to the risks associated
// with loading,  using,  modifying and/or developing or reproducing the
// software by the user in light of its specific status of free software,
// that may mean  that it is complicated to manipulate,  and  that  also
// therefore means  that it is reserved for developers  and  experienced
// professionals having in-depth computer knowledge. Users are therefore
// encouraged to load and test the software's suitability as regards their
// requirements in conditions enabling the security of their systems and/or 
// data to be ensured and,  more generally, to use and operate it in the 
// same conditions as regards security. 
// 
// The fact that you are presently reading this means that you have had
// knowledge of the CeCILL-C license and that you accept its terms.

package logging

import (
	"fmt"
	"io"
	"sync"
)

type IOWriterLogger struct {
	Destination io.Writer
	debug       bool
	logPrefix   string
	lock        sync.Mutex
}

func NewIOWriterLogger(destination io.Writer, debug bool) *IOWriterLogger {
	return &IOWriterLogger{
		Destination: destination,
		debug:       debug,
	}
}

func (l *IOWriterLogger) SetLogPrefix(prefix string) {
	l.logPrefix = prefix
}

func (l *IOWriterLogger) SetDebug(enable bool) {
	l.debug = enable
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
