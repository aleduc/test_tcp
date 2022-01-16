package internal

import (
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"strings"
)

// Уровни логирования в численном представлении.
const (
	levelDebug = iota
	levelInfo
	levelWarning
	levelError
)

// Поддерживаемые уровни логирования для использования в конструкторе.
const (
	LevelDebug   = "DEBUG"
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"
)

func parseLevel(level string) int {
	switch strings.ToUpper(level) {
	case LevelInfo:
		return levelInfo
	case LevelDebug:
		return levelDebug
	case LevelWarning:
		return levelWarning
	case LevelError:
		return levelError
	default:
		return levelInfo
	}
}

type printer interface {
	Println(v ...interface{})
	Printf(format string, a ...interface{})
}

// WithStackTrace является структурой для логгера.
type WithStackTrace struct {
	debug   printer
	info    printer
	warning printer
	error   printer
	fatal   printer
	level   int
}

// NewWithStackTrace создает Logger со стектрейсом.
func NewWithStackTrace(infoP printer, warningP printer, errorP printer, fatalP printer) *WithStackTrace {
	return &WithStackTrace{
		info:    infoP,
		warning: warningP,
		error:   errorP,
		fatal:   fatalP,
		level:   levelInfo,
	}
}

// NewLeveledWithStackTrace создает Logger со стектрейсом.
func NewLeveledWithStackTrace(debugP printer, infoP printer, warningP printer, errorP printer, fatalP printer, level string) *WithStackTrace {
	return &WithStackTrace{
		debug:   debugP,
		info:    infoP,
		warning: warningP,
		error:   errorP,
		fatal:   fatalP,
		level:   parseLevel(level),
	}
}

// Debug логирует сообщения в debugHandle поток.
func (l *WithStackTrace) Debug(debug string) {
	if l.level <= levelDebug {
		l.debug.Println(debug)
	}
}

// Debugf логирует сообщения в debugHandle поток.
func (l *WithStackTrace) Debugf(format string, args ...interface{}) {
	if l.level <= levelDebug {
		l.debug.Printf(format, args...)
	}
}

// Info логирует сообщения в infoHandle поток.
func (l *WithStackTrace) Info(info string) {
	if l.level <= levelInfo {
		l.info.Println(info)
	}
}

// Infof логирует сообщения в infoHandle поток.
func (l *WithStackTrace) Infof(format string, args ...interface{}) {
	if l.level <= levelInfo {
		l.info.Printf(format, args...)
	}
}

// Warning логгирует сообщения в warningHandle поток.
func (l *WithStackTrace) Warning(err error) {
	if l.level <= levelWarning {
		l.warning.Println(getErrorDetails(err))
	}
}

// Warningf логирует сообщения в warningHandle поток.
func (l *WithStackTrace) Warningf(format string, args ...interface{}) {
	if l.level <= levelWarning {
		l.Warning(fmt.Errorf(format, args...))
	}
}

// Error логирует сообщения в errorHandle поток.
func (l *WithStackTrace) Error(err error) {
	l.error.Println(getErrorDetails(err))
}

// Errorf логирует сообщения в errorHandle поток.
func (l *WithStackTrace) Errorf(format string, args ...interface{}) {
	if l.level <= levelError {
		l.Error(fmt.Errorf(format, args...))
	}
}

// Fatal логирует сообщения в fatalHandle поток.
func (l *WithStackTrace) Fatal(err error) {
	l.fatal.Println(getErrorDetails(err))
}

// Fatalf логирует сообщения в fatalHandle поток.
func (l *WithStackTrace) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Errorf(format, args...))
}

// Panic падает с паникой.
func (l *WithStackTrace) Panic(err error) {
	panic(getErrorDetails(err))
}

func getErrorDetails(err error) string {
	dropboxErr := errors.New(err.Error())
	return fmt.Sprintf("%s @ %v ", dropboxErr.GetMessage(), dropboxErr.StackFrames())
}
