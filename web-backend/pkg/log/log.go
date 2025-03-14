package log

import (
	"errors"
	"fmt"
	"io"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	entry *logrus.Entry
	// panic,fatal,error,warn,warning,info,debug,trace
	level       string
	printCaller bool
	caller      func() (file string, line int, funcName string, err error)
}

func (l *Logger) SetLevel(lvl string) {
	if lvl == "" {
		return
	}
	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		l.level = lvl
		l.entry.Logger.Level = level
	}
}
func (l *Logger) SetOutput(writer io.Writer) {
	l.entry.Logger.SetOutput(writer)
}
func (l *Logger) SetPrintCaller(printCaller bool) {
	l.printCaller = printCaller
}
func (l *Logger) getCallerInfo(level logrus.Level) map[string]interface{} {
	mp := make(map[string]interface{})
	if l.printCaller || level != logrus.InfoLevel {
		file, line, funName, err := l.caller()
		if err == nil {
			mp["file"] = fmt.Sprintf("%s:%d", file, line)
			mp["func"] = funName
		}
	}
	return mp
}

func (l *Logger) log(level logrus.Level, args ...interface{}) {
	l.entry.WithFields(l.getCallerInfo(level)).Log(level, args...)
}
func (l *Logger) logf(level logrus.Level, format string, args ...interface{}) {
	l.entry.WithFields(l.getCallerInfo(level)).Logf(level, format, args...)
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	entry := l.entry.WithFields(fields)
	return &Logger{entry: entry, level: l.level, printCaller: l.printCaller, caller: l.caller}
}

func (l *Logger) Info(args ...interface{}) {
	l.log(logrus.InfoLevel, args...)
}
func (l *Logger) Trace(args ...interface{}) {
	l.log(logrus.TraceLevel, args...)
}
func (l *Logger) Debug(args ...interface{}) {
	l.log(logrus.DebugLevel, args...)
}
func (l *Logger) Warn(args ...interface{}) {
	l.log(logrus.WarnLevel, args...)
}
func (l *Logger) Error(args ...interface{}) {
	l.log(logrus.ErrorLevel, args...)
}
func (l *Logger) Fatal(args ...interface{}) {
	l.log(logrus.FatalLevel, args...)
}
func (l *Logger) Panic(args ...interface{}) {
	l.log(logrus.PanicLevel, args...)
}

func (l *Logger) InfoF(format string, args ...interface{}) {
	l.logf(logrus.InfoLevel, format, args...)
}
func (l *Logger) TraceF(format string, args ...interface{}) {
	l.logf(logrus.TraceLevel, format, args...)
}
func (l *Logger) DebugF(format string, args ...interface{}) {
	l.logf(logrus.DebugLevel, format, args...)
}
func (l *Logger) WarnF(format string, args ...interface{}) {
	l.logf(logrus.WarnLevel, format, args...)
}
func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.logf(logrus.ErrorLevel, format, args...)
}
func (l *Logger) FatalF(format string, args ...interface{}) {
	l.logf(logrus.FatalLevel, format, args...)
}
func (l *Logger) PanicF(format string, args ...interface{}) {
	l.logf(logrus.PanicLevel, format, args...)
}

var log *Logger

func NewLogger() *Logger {
	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)
	log.AddHook(&errorHook{})
	Logger := &Logger{
		entry:  logrus.NewEntry(log),
		caller: defaultCaller,
	}
	return Logger
}

func init() {
	log = NewLogger()
}

func SetLevel(lvl string) {
	if lvl == "" {
		return
	}
	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		log.level = lvl
		log.entry.Logger.Level = level
	}
}

func SetOutput(writer io.Writer) {
	log.entry.Logger.SetOutput(writer)
}

func SetPrintCaller(printCaller bool) {
	log.printCaller = printCaller
}

func Info(args ...interface{}) {
	log.log(logrus.InfoLevel, args...)
}
func Trace(args ...interface{}) {
	log.log(logrus.TraceLevel, args...)
}
func Debug(args ...interface{}) {
	log.log(logrus.DebugLevel, args...)
}
func Warn(args ...interface{}) {
	log.log(logrus.WarnLevel, args...)
}
func Error(args ...interface{}) {
	log.log(logrus.ErrorLevel, args...)
}
func Fatal(args ...interface{}) {
	log.log(logrus.FatalLevel, args...)
}
func Panic(args ...interface{}) {
	log.log(logrus.PanicLevel, args...)
}

func InfoF(format string, args ...interface{}) {
	log.logf(logrus.InfoLevel, format, args...)
}
func TraceF(format string, args ...interface{}) {
	log.logf(logrus.TraceLevel, format, args...)
}
func DebugF(format string, args ...interface{}) {
	log.logf(logrus.DebugLevel, format, args...)
}
func WarnF(format string, args ...interface{}) {
	log.logf(logrus.WarnLevel, format, args...)
}
func ErrorF(format string, args ...interface{}) {
	log.logf(logrus.ErrorLevel, format, args...)
}
func FatalF(format string, args ...interface{}) {
	log.logf(logrus.FatalLevel, format, args...)
}
func PanicF(format string, args ...interface{}) {
	log.logf(logrus.PanicLevel, format, args...)
}

func WithFields(fields map[string]interface{}) *Logger {
	entry := log.entry.WithFields(fields)
	return &Logger{entry: entry, level: log.level, printCaller: log.printCaller, caller: log.caller}
}

func defaultCaller() (file string, line int, funcName string, err error) {
	pc, f, l, ok := runtime.Caller(4)
	if !ok {
		err = errors.New("caller failure")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	file, line = f, l
	return
}
