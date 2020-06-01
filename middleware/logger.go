package middleware

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"yee"
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

type Logger struct {
	sync.Mutex
	level int
}

type coloring func(string) string

func (l *Logger) SetLevel(level int) {
	l.Lock()
	defer l.Unlock()
	l.level = level
}

func newBrush(color string) coloring {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []coloring{
	newBrush("1;31"), // Critical          红色
	newBrush("1;35"), // Error              紫色
	newBrush("1;33"), // Warn           黄色
	newBrush("1;37"), // Trace               白色
	newBrush("1;34"), // Info              蓝色
	newBrush("1;37"), // Debug               白色
}

var mappingLevel = map[int]string{
	Critical: "Critical",
	Error:    "Error",
	Warning:  "Warn",
	Trace:    "Trace",
	Info:     "Info",
	Debug:    "Debug",
}

func (l *Logger) logWrite(msgText string, level int) error {
	if level > l.level {
		return nil
	}

	if level != Trace {
		_, file, lineno, ok := runtime.Caller(2)

		src := ""

		if ok {
			src = strings.Replace(
				fmt.Sprintf("%s:%d", file, lineno), "%2e", ".", -1)
		}
		msgText = fmt.Sprintf("%s [%s] %s %s: %s", yee.YeeVersion, mappingLevel[level], time.Now().Format(timeFormat), src, msgText)
	} else {
		msgText = fmt.Sprintf("%s [%s] %s %s", yee.YeeVersion, mappingLevel[level], time.Now().Format(timeFormat), msgText)
	}

	if runtime.GOOS != "windows" {
		msgText = colors[level](msgText)
	}

	l.print(msgText)

	return nil
}

func (l *Logger) print(msg string) {
	l.Lock()
	defer l.Unlock()
	_, _ = os.Stdout.Write(append([]byte(msg), '\n'))
}

func (l *Logger) Critical(msg string) {
	l.logWrite(msg, Critical)
}

func (l *Logger) Error(msg string) {
	l.logWrite(msg, Error)
}

func (l *Logger) Warn(msg string) {
	l.logWrite(msg, Warning)
}

func (l *Logger) Info(msg string) {
	l.logWrite(msg, Info)
}

func (l *Logger) Debug(msg string) {
	l.logWrite(msg, Debug)
}

func (l *Logger) Trace(msg string) {
	l.logWrite(msg, Trace)
}
