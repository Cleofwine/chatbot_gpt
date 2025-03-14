package log

import (
	"chatgpt-proxy/pkg/config"
	"errors"
	"fmt"
	"io"
	"runtime"

	"github.com/sirupsen/logrus"
)

type logger struct {
	entry *logrus.Entry
	// panic,fatal,error,warn,warning,info,debug,trace
	level       string
	printCaller bool
	caller      func() (file string, line int, funcName string, err error)
}

func (l *logger) SetLevel(lvl string) {
	if lvl == "" {
		return
	}
	level, err := logrus.ParseLevel(lvl)
	if err == nil {
		l.level = lvl
		l.entry.Logger.Level = level
	}
}
func (l *logger) SetOutput(writer io.Writer) {
	l.entry.Logger.SetOutput(writer)
}
func (l *logger) SetPrintCaller(printCaller bool) {
	l.printCaller = printCaller
}
func (l *logger) getCallerInfo(level logrus.Level) map[string]interface{} {
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

func (l *logger) log(level logrus.Level, args ...interface{}) {
	l.entry.WithFields(l.getCallerInfo(level)).Log(level, args...)
}
func (l *logger) logf(level logrus.Level, format string, args ...interface{}) {
	l.entry.WithFields(l.getCallerInfo(level)).Logf(level, format, args...)
}

func (l *logger) WithFields(fields map[string]interface{}) *logger {
	entry := l.entry.WithFields(fields)
	return &logger{entry: entry, level: l.level, printCaller: l.printCaller, caller: l.caller}
}

var log *logger

func NewLogger() *logger {
	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)
	log.AddHook(&errorHook{})
	log.SetOutput(getRotateWriter())
	logger := &logger{
		entry:  logrus.NewEntry(log),
		caller: defaultCaller,
	}
	return logger
}

func init() {
	cnf := config.GetConf()
	log = NewLogger()
	log.SetLevel(cnf.Log.Level)
	log.SetOutput(getRotateWriter())
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

func WithFields(fields map[string]interface{}) *logger {
	entry := log.entry.WithFields(fields)
	return &logger{entry: entry, level: log.level, printCaller: log.printCaller, caller: log.caller}
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
