package util

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		log.New(os.Stdout, "ifandonlyif ", 0),
	}
}

func (l *Logger) Infof(s string, i ...interface{}) {
	l.Logger.Println("[INFO]", s, i)
}
func (l *Logger) Debugf(s string, i ...interface{}) {
	l.Logger.Println("[DEBUG]", s, i)
}
func (l *Logger) Errorf(s string, i ...interface{}) {
	l.Logger.Println("[ERROR]", s, i)
}
