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
