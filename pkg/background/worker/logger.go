package worker

import "github.com/charmbracelet/log"

type Logger struct {
	log *log.Logger
}

func NewLogger(log *log.Logger) *Logger {
	return &Logger{log: log}
}

func (l *Logger) Debug(v ...interface{}) {
	l.log.Debug(v[0], v[1:]...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.log.Debugf(format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.log.Error(v[0], v[1:]...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.log.Errorf(format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.log.Info(v[0], v[1:]...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.log.Infof(format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.log.Warn(v[0], v[1:]...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.log.Warnf(format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log.Fatal(v[0], v[1:]...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.log.Fatalf(format, v...)
}
