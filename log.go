package yee

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cookieY/yee/color"
)

// logger types
const (
	Critical = iota
	Error
	Warning
	Info
	Debug
)

const timeFormat = "2006-01-02 15:04:05"

type logger struct {
	sync.Mutex
	level    uint8
	isLogger bool
	producer *color.Color
}

// Logger ...
type Logger interface {
	Critical(msg string)
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
	Criticalf(error string, msg ...interface{})
	Errorf(error string, msg ...interface{})
	Warnf(error string, msg ...interface{})
	Infof(error string, msg ...interface{})
	Debugf(error string, msg ...interface{})
	SetLevel(level uint8)
}

// LogCreator ...
func LogCreator() *logger {
	l := new(logger)
	l.producer = color.New()
	l.producer.Enable()
	l.level = 1
	return l
}

func (l *logger) SetLevel(level uint8) {
	l.Lock()
	defer l.Unlock()
	l.level = level
}

func (l *logger) IsLogger(p bool) {
	l.Lock()
	defer l.Unlock()
	l.isLogger = p
}

var mappingLevel = map[uint8]string{
	Critical: "Critical",
	Error:    "Error",
	Warning:  "Warn",
	Info:     "Info",
	Debug:    "Debug",
}

func (l *logger) logWrite(msgText string, level uint8) (string, bool) {
	if level > l.level && !l.isLogger {
		return "", false
	}

	if !l.isLogger {
		_, file, lineno, ok := runtime.Caller(2)

		src := ""

		if ok {
			src = strings.Replace(
				fmt.Sprintf("%s:%d", file, lineno), "%2e", ".", -1)
		}
		msgText = fmt.Sprintf("%s [%s] %s (%s) %s", version, mappingLevel[level], time.Now().Format(timeFormat), src, msgText)
	} else {
		msgText = fmt.Sprintf("%s [%s] %s %s", version, mappingLevel[level], time.Now().Format(timeFormat), msgText)
	}

	return msgText, true
}

func (l *logger) print(msg string) {
	l.Lock()
	defer l.Unlock()
	_, _ = os.Stdout.Write(append([]byte(msg), '\n'))
}

func (l *logger) Critical(msg string) {
	if msg, ok := l.logWrite(msg, Critical); ok {
		l.print(l.producer.Red(msg))
	}
}

func (l *logger) Criticalf(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Critical); ok {
		l.print(l.producer.Red(msg))
	}
}

func (l *logger) Error(msg string) {
	if msg, ok := l.logWrite(msg, Error); ok {
		l.print(l.producer.Magenta(msg))
	}
}

func (l *logger) Errorf(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Error); ok {
		l.print(l.producer.Magenta(msg))
	}
}

func (l *logger) Warn(msg string) {
	if msg, ok := l.logWrite(msg, Warning); ok {
		l.print(l.producer.Yellow(msg))
	}
}

func (l *logger) Warnf(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Warning); ok {
		l.print(l.producer.Yellow(msg))
	}
}

func (l *logger) Info(msg string) {
	if msg, ok := l.logWrite(msg, Info); ok {
		l.print(l.producer.Blue(msg))
	}
}

func (l *logger) Infof(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Info); ok {
		l.print(l.producer.Blue(msg))
	}
}

func (l *logger) Debugf(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Debug); ok {
		l.print(l.producer.Cyan(msg))
	}
}

func (l *logger) Debug(msg string) {
	if msg, ok := l.logWrite(msg, Debug); ok {
		l.print(l.producer.Cyan(msg))
	}
}
