package main

import "github.com/sirupsen/logrus"


type LoggerInterface interface{
    Logf(format string, args ...interface{})
    Debugf(format string, args ...interface{})
}


type Logger struct{}


func NewLogger() Logger{
	return Logger{}
}


func (l Logger) Debugf(format string, args ...interface{}){
	logrus.Debugf(format, args...)
}


func (l Logger) Logf(format string, args ...interface{}){
	logrus.Printf(format, args...)
}
