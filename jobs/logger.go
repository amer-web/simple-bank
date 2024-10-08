package jobs

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}
func (l *Logger) Handle(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}
func (l *Logger) Debug(args ...interface{}) {
	l.Handle(zerolog.DebugLevel, args...)
}
func (l *Logger) Info(args ...interface{}) {
	l.Handle(zerolog.InfoLevel, args...)
}
func (l *Logger) Warn(args ...interface{}) {
	l.Handle(zerolog.WarnLevel, args...)
}
func (l *Logger) Error(args ...interface{}) {
	l.Handle(zerolog.ErrorLevel, args...)
}
func (l *Logger) Fatal(args ...interface{}) {
	l.Handle(zerolog.FatalLevel, args...)
}
