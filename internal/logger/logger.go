package logger

import (
	c "github.com/fatih/color"
)

type Logger struct{}

func (l *Logger) Info(msg string, a ...any) {
	c.Blue(msg, a...)
}

func (l *Logger) Error(msg string, a ...any) {
	c.Red(msg, a...)
}

func (l *Logger) Print(msg string, a ...any) {
	c.HiWhite(msg, a...)
}
