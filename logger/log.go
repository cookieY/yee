package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
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
	version  string
	producer *Color
	out      io.Writer
	noColor  bool
}

// Logger ...
type Logger interface {
	Critical(msg interface{})
	Error(msg interface{})
	Warn(msg interface{})
	Info(msg interface{})
	Debug(msg interface{})
	Criticalf(error string, msg ...interface{})
	Errorf(error string, msg ...interface{})
	Warnf(error string, msg ...interface{})
	Infof(error string, msg ...interface{})
	Debugf(error string, msg ...interface{})
	Custom(msg string)
	SetLevel(level uint8)
	SetOut(out io.Writer)
	IsLogger(isOk bool)
}

// LogCreator ...
func LogCreator(args ...int) Logger {
	l := new(logger)
	l.producer = New()
	l.producer.Enable()
	l.level = 1
	if len(args) > 0 {
		l.level = uint8(args[0])
	}
	l.out = os.Stdout
	return l
}

var DefaultLogger = LogCreator()

func (l *logger) SetOut(out io.Writer) {
	l.Lock()
	defer l.Unlock()
	l.out = out
	l.noColor = true
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

func (l *logger) logWrite(msg interface{}, level uint8) (string, bool) {
	var msgText string
	switch v := msg.(type) {
	case error:
		msgText = v.Error()
	case string:
		msgText = v
	}

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
		msgText = fmt.Sprintf("%s [%s] %s (%s) %s", l.version, mappingLevel[level], time.Now().Format(timeFormat), src, msgText)
	} else {
		msgText = fmt.Sprintf("%s [%s] %s %s", l.version, mappingLevel[level], time.Now().Format(timeFormat), msgText)
	}

	return msgText, true
}

func (l *logger) print(msg string) {
	l.Lock()
	defer l.Unlock()
	_, err := l.out.Write(append([]byte(msg), '\n'))
	if err != nil {
		_, _ = os.Stdout.Write(append([]byte(msg), '\n'))
	}
}

func (l *logger) Custom(msg string) {
	l.print(msg)
}

func (l *logger) Critical(msg interface{}) {
	if msg, ok := l.logWrite(msg, Critical); ok {
		l.dyer(Critical, &msg)
	}
}

func (l *logger) Criticalf(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Critical); ok {
		l.dyer(Critical, &msg)
	}
}

func (l *logger) Error(msg interface{}) {
	if msg, ok := l.logWrite(msg, Error); ok {
		l.dyer(Error, &msg)
	}
}

func (l *logger) Errorf(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Error); ok {
		l.dyer(Error, &msg)
	}
}

func (l *logger) Warn(msg interface{}) {
	if msg, ok := l.logWrite(msg, Warning); ok {
		l.dyer(Warning, &msg)
	}
}

func (l *logger) Warnf(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Warning); ok {
		l.dyer(Warning, &msg)
	}
}

func (l *logger) Info(msg interface{}) {
	if msg, ok := l.logWrite(msg, Info); ok {
		l.dyer(Info, &msg)
	}
}

func (l *logger) Infof(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Info); ok {
		l.dyer(Info, &msg)
	}
}

func (l *logger) Debugf(error string, msg ...interface{}) {
	if msg, ok := l.logWrite(fmt.Sprintf(error, msg...), Debug); ok {
		l.dyer(Debug, &msg)
	}
}

func (l *logger) Debug(msg interface{}) {
	if msg, ok := l.logWrite(msg, Debug); ok {
		l.dyer(Debug, &msg)
	}
}

func (l *logger) dyer(level int, msg *string) {
	if l.noColor {
		l.print(*msg)
		return
	}
	switch level {
	case Critical:
		l.print(l.producer.Red(*msg))
	case Error:
		l.print(l.producer.Magenta(*msg))
	case Warning:
		l.print(l.producer.Yellow(*msg))
	case Info:
		l.print(l.producer.Blue(*msg))
	case Debug:
		l.print(l.producer.Cyan(*msg))
	}
}
