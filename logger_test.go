package logger

import (
	"bytes"
	"io"
	"log"
	"strings"
	"testing"
)

type mockBaseLogger struct {
	buf *bytes.Buffer
}

func (m *mockBaseLogger) Printf(format string, v ...any) {
	m.buf.WriteString(strings.TrimRight(
		strings.ReplaceAll(format, "%v", ""),
		" ",
	))
}

func (m *mockBaseLogger) Print(v ...any) {
	for _, val := range v {
		m.buf.WriteString(val.(string))
	}
}

func (m *mockBaseLogger) Writer() io.Writer {
	return m.buf
}

func newTestLogger(level LogLevel, componentLevel map[string]LogLevel) (*LoggerImpl, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	mock := &mockBaseLogger{buf: buf}

	config := &Config{
		Level:          level,
		ComponentLevel: componentLevel,
	}

	l := &LoggerImpl{
		config:     config,
		baseLogger: log.New(mock.Writer(), "", 0),
	}

	return l, buf
}

func TestSetFlags(t *testing.T) {
	l, _ := newTestLogger(LogDebug, nil)

	l.SetFlags(0)

	logger, ok := l.baseLogger.(*log.Logger)
	if !ok {
		t.Fatal("expected baseLogger to be *log.Logger")
	}

	if logger.Flags() != 0 {
		t.Errorf("expected flags 0, got %d", logger.Flags())
	}

	l.SetFlags(log.Ldate | log.Ltime)

	if logger.Flags() != log.Ldate|log.Ltime {
		t.Errorf("expected flags %d, got %d", log.Ldate|log.Ltime, logger.Flags())
	}

	l.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if logger.Flags() != log.Ldate|log.Ltime|log.Lmicroseconds {
		t.Errorf("expected flags %d, got %d", log.Ldate|log.Ltime|log.Lmicroseconds, logger.Flags())
	}
}

func TestSetFlags_NonStdLogger(t *testing.T) {
	mock := &mockBaseLogger{buf: &bytes.Buffer{}}

	l := &LoggerImpl{
		config:     &Config{Level: LogDebug},
		baseLogger: mock,
	}

	l.SetFlags(log.Ldate)
}

func TestInfo(t *testing.T) {
	l, buf := newTestLogger(LogDebug, nil)
	l.SetFlags(0)

	l.Info("hello %d", 42)
	got := buf.String()
	if !strings.Contains(got, "[inf] hello 42") {
		t.Errorf("expected info message, got %q", got)
	}
}

func TestInfo_NoArgs(t *testing.T) {
	l, buf := newTestLogger(LogDebug, nil)
	l.SetFlags(0)

	l.Info("hello")
	got := buf.String()
	if !strings.Contains(got, "[inf] hello") {
		t.Errorf("expected info message, got %q", got)
	}
}

func TestError(t *testing.T) {
	l, buf := newTestLogger(LogError, nil)
	l.SetFlags(0)

	l.Error("fail %d", 1)
	got := buf.String()
	if !strings.Contains(got, "[err] fail 1") {
		t.Errorf("expected error message, got %q", got)
	}
}

func TestDebug(t *testing.T) {
	l, buf := newTestLogger(LogDebug, nil)
	l.SetFlags(0)

	l.Debug("trace %d", 1)
	got := buf.String()
	if !strings.Contains(got, "[dbg] trace 1") {
		t.Errorf("expected debug message, got %q", got)
	}
}

func TestLogLevelFiltering(t *testing.T) {
	l, buf := newTestLogger(LogError, nil)
	l.SetFlags(0)

	l.Info("should not appear")
	l.Debug("should not appear")

	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}
}

func TestWithPrefix(t *testing.T) {
	l, buf := newTestLogger(LogDebug, nil)
	l.SetFlags(0)

	prefixed := l.WithPrefix("db")
	prefixed.Info("connected")

	got := buf.String()
	if !strings.Contains(got, "[inf] [db] connected") {
		t.Errorf("expected prefixed message, got %q", got)
	}
}

func TestComponentLevel(t *testing.T) {
	l, buf := newTestLogger(LogDebug, map[string]LogLevel{
		"db": LogError,
	})
	l.SetFlags(0)

	dbL := l.WithPrefix("db")
	dbL.Info("should not appear")

	if buf.Len() != 0 {
		t.Errorf("expected no output for db Info, got %q", buf.String())
	}

	dbL.Error("connection failed")
	got := buf.String()
	if !strings.Contains(got, "[err] [db] connection failed") {
		t.Errorf("expected error message, got %q", got)
	}
}

func TestSetConfig(t *testing.T) {
	l, buf := newTestLogger(LogError, nil)
	l.SetFlags(0)

	l.Info("before")
	if buf.Len() != 0 {
		t.Errorf("expected no output, got %q", buf.String())
	}

	l.SetConfig(&Config{Level: LogDebug})
	l.Info("after")
	got := buf.String()
	if !strings.Contains(got, "[inf] after") {
		t.Errorf("expected info message after config change, got %q", got)
	}
}

func TestSetConfig_Nil(t *testing.T) {
	l, _ := newTestLogger(LogError, nil)
	l.SetConfig(nil)

	if l.config.Level != LogError {
		t.Errorf("expected level to remain LogError, got %s", l.config.Level)
	}
}
