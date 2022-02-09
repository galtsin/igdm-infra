package main

import (
	"fmt"
	"io"
	"log"

	"channels-instagram-dm/domain"
)

type Logger struct {
	path   string
	logger *log.Logger
}

func NewLogger(out io.Writer, prefix string) domain.Logger {
	return &Logger{
		logger: log.New(out, prefix, log.LstdFlags|log.LUTC),
	}
}

func (l *Logger) Copy(pathPrefix string) domain.Logger {
	loggerNew := NewLogger(l.Writer(), l.logger.Prefix())
	ln, _ := loggerNew.(*Logger)
	ln.path = l.path + " " + pathPrefix
	return ln
}

func (l *Logger) Writer() io.Writer {
	return l.logger.Writer()
}

func (l *Logger) Critical(message string, payload interface{}) {
	l.logger.Println("C:"+l.path+" "+message, fmt.Sprintf("%+v", payload))
}

func (l *Logger) Error(message string, payload interface{}) {
	l.logger.Println("E:"+l.path+" "+message, fmt.Sprintf("%+v", payload))
}

func (l *Logger) Debug(message string, payload interface{}) {
	l.logger.Println("D:"+l.path+" "+message, fmt.Sprintf("%+v", payload))
}

func (l *Logger) Info(message string, payload interface{}) {
	l.logger.Println("I:"+l.path+" "+message, fmt.Sprintf("%+v", payload))
}
