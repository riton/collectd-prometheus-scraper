package logging

type Logger interface {
	SetDebug(bool)
	SetLogPrefix(string)
	Info(string) error
	Infof(string, ...interface{}) error
	Error(string) error
	Errorf(string, ...interface{}) error
	Warning(string) error
	Warningf(string, ...interface{}) error
	Debug(string) error
	Debugf(string, ...interface{}) error
}
