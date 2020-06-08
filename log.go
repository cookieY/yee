package yee

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	Critical = iota
	Error
	Warning
	Trace
	Info
	Debug
)

const timeFormat = "2006-01-02 15:04:05"

type logger struct {
	sync.Mutex
	level    uint8
	isLogger bool
}

type Logger interface {
	Critical(msg string)
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
	Trace(msg string)
	SetLevel(level uint8)
}

type coloring func(string) string

func LogCreator() *logger {
	return new(logger)
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

func newBrush(color string) coloring {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []coloring{
	newBrush("0;31"), // Critical          红色
	newBrush("0;35"), // Error              紫色
	newBrush("0;33"), // Warn           黄色
	newBrush("0;34"), // Trace               白色
	newBrush("0;34"), // Info              蓝色
	newBrush("0;37"), // Debug               白色
}

var mappingLevel = map[uint8]string{
	Critical: "Critical",
	Error:    "Error",
	Warning:  "Warn",
	Trace:    "Trace",
	Info:     "Info",
	Debug:    "Debug",
}

func (l *logger) logWrite(msgText string, level uint8) {
	if level > l.level && !l.isLogger {
		return
	}

	if !l.isLogger {
		_, file, lineno, ok := runtime.Caller(2)

		src := ""

		if ok {
			src = strings.Replace(
				fmt.Sprintf("%s:%d", file, lineno), "%2e", ".", -1)
		}
		msgText = fmt.Sprintf("%s [%s] %s (%s) %s", Version, mappingLevel[level], time.Now().Format(timeFormat), src, msgText)
	} else {
		msgText = fmt.Sprintf("%s [%s] %s %s", Version, mappingLevel[level], time.Now().Format(timeFormat), msgText)
	}

	if runtime.GOOS != "windows" {
		msgText = colors[level](msgText)
	}

	l.print(msgText)

	return
}

func (l *logger) print(msg string) {
	l.Lock()
	defer l.Unlock()
	_, _ = os.Stdout.Write(append([]byte(msg), '\n'))
}

func (l *logger) Critical(msg string) {
	l.logWrite(msg, Critical)
}

func (l *logger) Error(msg string) {
	l.logWrite(msg, Error)
}

func (l *logger) Warn(msg string) {
	l.logWrite(msg, Warning)
}

func (l *logger) Info(msg string) {
	l.logWrite(msg, Info)
}

func (l *logger) Debug(msg string) {
	l.logWrite(msg, Debug)
}

func (l *logger) Trace(msg string) {
	l.logWrite(msg, Trace)
}
