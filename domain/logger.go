package domain

import "io"

type Logger interface {
	Copy(prefix string) Logger
	Writer() io.Writer
	Critical(message string, payload interface{})
	Debug(message string, payload interface{})
	Info(message string, payload interface{})
	Error(message string, payload interface{})
}
