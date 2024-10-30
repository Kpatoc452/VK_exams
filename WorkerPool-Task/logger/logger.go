package logger

import "github.com/sirupsen/logrus"


type Logger interface{
    Logf(format string, args ...interface{})
    Debugf(format string, args ...interface{})
}


type logger struct{}


func New() Logger{
	return logger{}
}


func (l logger) Debugf(format string, args ...interface{}){
	logrus.Debugf(format, args...)
}


func (l logger) Logf(format string, args ...interface{}){
	logrus.Printf(format, args...)
}
