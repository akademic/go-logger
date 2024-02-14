package logger

import (
	"log"
)

type baseLogger interface {
	Printf(format string, v ...any)
	Print(v ...any)
}

type LoggerImpl struct {
	config     *Config
	prefix     string
	baseLogger baseLogger
}

func New(prefix string, config *Config) *LoggerImpl {
	return &LoggerImpl{
		prefix:     prefix,
		config:     config,
		baseLogger: log.New(log.Writer(), "", log.LstdFlags),
	}
}

func (l *LoggerImpl) WithPrefix(prefix string) Logger {
	return &LoggerImpl{
		prefix:     prefix,
		config:     l.config,
		baseLogger: l.baseLogger,
	}
}

func (l *LoggerImpl) Info(pattern string, args ...any) {
	if !l.logOn(LogInfo) {
		return
	}
	pattern = l.pattern(pattern)
	if len(args) == 0 {
		l.baseLogger.Print("[inf] " + pattern)
		return
	}
	l.baseLogger.Printf("[inf] "+pattern, args...)
}

func (l *LoggerImpl) Error(pattern string, args ...any) {
	if !l.logOn(LogError) {
		return
	}

	pattern = l.pattern(pattern)
	if len(args) == 0 {
		l.baseLogger.Print("[err] " + pattern)
		return
	}
	l.baseLogger.Printf("[err] "+pattern, args...)
}

func (l *LoggerImpl) Debug(pattern string, args ...any) {
	if !l.logOn(LogDebug) {
		return
	}
	pattern = l.pattern(pattern)
	if len(args) == 0 {
		l.baseLogger.Print("[dbg] " + pattern)
		return
	}
	l.baseLogger.Printf("[dbg] "+pattern, args...)
}

func (l *LoggerImpl) SetConfig(config *Config) {
	l.config.Level = config.Level
}

func (l *LoggerImpl) logOn(level LogLevel) bool {
	if l.prefix == "" {
		return l.config.Level.CanLog(level)
	}

	confLevel, ok := l.config.ComponentLevel[l.prefix]
	if !ok {
		return l.config.Level.CanLog(level)
	}

	return confLevel.CanLog(level)
}

func (l *LoggerImpl) pattern(pattern string) string {
	if l.prefix != "" {
		pattern = "[" + l.prefix + "] " + pattern
	}

	return pattern
}
