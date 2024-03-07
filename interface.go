// Package logger contains project-wide logging interface
package logger

import "io"

type Logger interface {
	WithPrefix(prefix string) Logger
	Info(pattern string, args ...interface{})
	Error(pattern string, args ...interface{})
	Debug(pattern string, args ...interface{})
	Writer() io.Writer
}
