package store

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

type LogFields map[string]interface{}

type Logger struct {
	logrus.Logger
}

func LogLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}

func Logger_new(l string) *Logger {
	logger := &Logger{logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}}
	if l == "" {
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = LogLevel(l)
	}
	return logger
}

func (p *Logger) Log(logtype string, title string, fields LogFields) {
	if fields == nil {
		return
	}
	switch {
	case logtype == "debug":
		p.WithFields(logrus.Fields(fields)).Debug(title)
	case logtype == "warn":
		p.WithFields(logrus.Fields(fields)).Warn(title)
	case logtype == "error":
		p.WithFields(logrus.Fields(fields)).Error(title)
	case logtype == "fatal":
		p.WithFields(logrus.Fields(fields)).Fatal(title)
	case logtype == "panic":
		p.WithFields(logrus.Fields(fields)).Panic(title)
	default:
		p.WithFields(logrus.Fields(fields)).Info(title)
	}
}

// This can be used as the destination for a logger and it'll
// map them into calls to testing.T.Log, so that you only see
// the logging for failed tests.
type testLoggerAdapter struct {
	t      *testing.T
	prefix string
}

func (a *testLoggerAdapter) Write(d []byte) (int, error) {
	if d[len(d)-1] == '\n' {
		d = d[:len(d)-1]
	}
	if a.prefix != "" {
		l := a.prefix + ": " + string(d)
		a.t.Log(l)
		return len(l), nil
	}
	a.t.Log(string(d))
	return len(d), nil
}

func NewTestLogger(t *testing.T) *Logger {
	logger := Logger_new("debug")
	logger.Out = &testLoggerAdapter{t: t}
	return logger
}

type benchmarkLoggerAdapter struct {
	b      *testing.B
	prefix string
}

func (b *benchmarkLoggerAdapter) Write(d []byte) (int, error) {
	if d[len(d)-1] == '\n' {
		d = d[:len(d)-1]
	}
	if b.prefix != "" {
		l := b.prefix + ": " + string(d)
		b.b.Log(l)
		return len(l), nil
	}

	b.b.Log(string(d))
	return len(d), nil
}

func NewBenchmarkLogger(b *testing.B) *Logger {
	logger := Logger_new("debug")
	logger.Out = &benchmarkLoggerAdapter{b: b}
	logger.Level = logrus.DebugLevel
	return logger
}
