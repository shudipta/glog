package glog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger interface {
	Meta(meta int) *logger
	LogType(lType int) *logger
	Skip(skip int) *logger

	Log(lType int, str string) (n int, err error)

	Print(a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
	Printf(format string, a ...interface{}) (n int, err error)

	Info(a ...interface{}) (n int, err error)
	Infoln(a ...interface{}) (n int, err error)
	Infof(format string, a ...interface{}) (n int, err error)

	Warn(a ...interface{}) (n int, err error)
	Warnln(a ...interface{}) (n int, err error)
	Warnf(format string, a ...interface{}) (n int, err error)

	Debug(a ...interface{}) (n int, err error)
	Debugln(a ...interface{}) (n int, err error)
	Debugf(format string, a ...interface{}) (n int, err error)

	Error(err1 error) (n int, err error)
	Errorln(err1 error) (n int, err error)
	Errorf(format string, a ...interface{}) (n int, err error)
}

var _ Logger = &logger{}

const (
	MetaDate int = 1 << iota
	MetaTime
	MetaZone
	MetaCallerFull
	MetaCallerShort
	MetaFull    = MetaDate | MetaTime | MetaZone | MetaCallerFull
	MetaShort   = MetaTime | MetaCallerShort
	MetaDefault = MetaDate | MetaTime | MetaZone | MetaCallerShort
)

const (
	LogTypePrint int = 1 << iota
	LogTypeInfo
	LogTypeWarn
	LogTypeDebug
	LogTypeError
	LogTypeDefault = LogTypePrint | LogTypeInfo | LogTypeWarn | LogTypeDebug | LogTypeError
)

var timeFormat = map[int]string{
	MetaDate: "2006-01-02",
	MetaTime: "T15:04:05",
	MetaZone: "Z07:00",
}

var logType = map[int]string{
	LogTypePrint: "",
	LogTypeInfo:  "[INFO] ",
	LogTypeWarn:  "[WARN] ",
	LogTypeDebug: "[DEBUG] ",
	LogTypeError: "[ERROR] ",
}

const (
	White  = "\033[3;30m%s\033[0m"
	Cyan   = "\033[1;36m%s\033[0m"
	Yellow = "\033[1;33m%s\033[0m"
	Purple = "\033[1;35m%s\033[0m"
	Red    = "\033[1;31m%s\033[0m"
	Green  = "\033[1;32m%s\033[0m"
	Blue   = "\033[1;34m%s\033[0m"
	Ash    = "\033[1;37m%s\033[0m"
)

var coloredStr = map[int]string{
	LogTypePrint: White,
	LogTypeInfo:  Cyan,
	LogTypeWarn:  Yellow,
	LogTypeDebug: Purple,
	LogTypeError: Red,
}

var myLogger Logger = &logger{
	meta:    MetaDefault,
	logType: LogTypeDefault,
	skip:    0,
}

func Default() Logger {
	return myLogger
}

func New(meta, logType int) Logger {
	return &logger{
		meta:    meta,
		logType: logType,
		skip:    0,
	}
}

type logger struct {
	meta    int
	logType int
	skip    int
}

func Meta(meta int) *logger { return myLogger.Meta(meta) }
func (l *logger) Meta(meta int) *logger {
	l.meta = meta
	return l
}

func LogType(lType int) *logger { return myLogger.LogType(lType) }
func (l *logger) LogType(lType int) *logger {
	l.logType = lType
	return l
}

func Skip(skip int) *logger { return myLogger.Skip(skip) }
func (l *logger) Skip(skip int) *logger {
	l.skip = skip
	return l
}

func Log(lType int, str string) (n int, err error) { return myLogger.Skip(1).Log(lType, str) }
func (l *logger) Log(lType int, str string) (n int, err error) {
	return l.log(lType, str)
}

func Print(a ...interface{}) (n int, err error) { return myLogger.Skip(1).Print(a...) }
func (l *logger) Print(a ...interface{}) (n int, err error) {
	return l.log(LogTypePrint, fmt.Sprint(a...))
}

func Println(a ...interface{}) (n int, err error) { return myLogger.Skip(1).Println(a...) }
func (l *logger) Println(a ...interface{}) (n int, err error) {
	return l.log(LogTypePrint, fmt.Sprintln(a...))
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return myLogger.Skip(1).Printf(format, a...)
}
func (l *logger) Printf(format string, a ...interface{}) (n int, err error) {
	return l.log(LogTypePrint, fmt.Sprintf(format, a...))
}

func Info(a ...interface{}) (n int, err error) { return myLogger.Skip(1).Info(a...) }
func (l *logger) Info(a ...interface{}) (n int, err error) {
	return l.log(LogTypeInfo, fmt.Sprint(a...))
}

func Infoln(a ...interface{}) (n int, err error) { return myLogger.Skip(1).Infoln(a...) }
func (l *logger) Infoln(a ...interface{}) (n int, err error) {

	return l.log(LogTypeInfo, fmt.Sprintln(a...))
}

func Infof(format string, a ...interface{}) (n int, err error) {
	return myLogger.Skip(1).Infof(format, a...)
}
func (l *logger) Infof(format string, a ...interface{}) (n int, err error) {
	return l.log(LogTypeInfo, fmt.Sprintf(format, a...))
}

func Warn(a ...interface{}) (n int, err error) { return myLogger.Skip(1).Warn(a...) }
func (l *logger) Warn(a ...interface{}) (n int, err error) {
	return l.log(LogTypeWarn, fmt.Sprint(a...))
}

func Warnln(a ...interface{}) (n int, err error) { return myLogger.Skip(1).Warnln(a...) }
func (l *logger) Warnln(a ...interface{}) (n int, err error) {
	return l.log(LogTypeWarn, fmt.Sprintln(a...))
}

func Warnf(format string, a ...interface{}) (n int, err error) {
	return myLogger.Skip(1).Warnf(format, a...)
}
func (l *logger) Warnf(format string, a ...interface{}) (n int, err error) {
	return l.log(LogTypeWarn, fmt.Sprintf(format, a...))
}

func Debug(a ...interface{}) (n int, err error) { return myLogger.Skip(1).Debug(a...) }
func (l *logger) Debug(a ...interface{}) (n int, err error) {
	return l.log(LogTypeDebug, fmt.Sprint(a...))
}

func Debugln(a ...interface{}) (n int, err error) { return myLogger.Skip(1).Debugln(a...) }
func (l *logger) Debugln(a ...interface{}) (n int, err error) {
	return l.log(LogTypeDebug, fmt.Sprintln(a...))
}

func Debugf(format string, a ...interface{}) (n int, err error) {
	return myLogger.Skip(1).Debugf(format, a...)
}
func (l *logger) Debugf(format string, a ...interface{}) (n int, err error) {
	return l.log(LogTypeDebug, fmt.Sprintf(format, a...))
}

func Error(err1 error) (n int, err error) { return myLogger.Skip(1).Error(err1) }
func (l *logger) Error(err1 error) (n int, err error) {
	return l.log(LogTypeError, fmt.Sprint(err1))
}

func Errorln(err1 error) (n int, err error) { return myLogger.Skip(1).Errorln(err1) }
func (l *logger) Errorln(err1 error) (n int, err error) {
	return l.log(LogTypeError, fmt.Sprintln(err1))
}

func Errorf(format string, a ...interface{}) (n int, err error) {
	return myLogger.Skip(1).Errorf(format, a...)
}
func (l *logger) Errorf(format string, a ...interface{}) (n int, err error) {
	return l.log(LogTypeError, fmt.Sprintf(format, a...))
}

func (l *logger) timeFormatStr() (format string) {
	if l.meta&MetaDate > 0 {
		format = timeFormat[MetaDate] + " "
	}
	if l.meta&MetaTime > 0 {
		format += timeFormat[MetaTime] + " "
	}
	if l.meta&MetaZone > 0 {
		format += timeFormat[MetaZone] + " "
	}

	return format
}

func (l *logger) caller() (caller string) {
	if _, path, line, ok := runtime.Caller(3 + l.skip); ok {
		if l.meta&MetaCallerShort > 0 {
			caller = fmt.Sprintf("%s:%d ", strings.TrimPrefix(path, os.Getenv("PWD"))[1:], line)
		} else if l.meta&MetaCallerFull > 0 {
			caller = fmt.Sprintf("%s:%d ", path, line)
		}
	}
	return caller
}

func (l *logger) log(lType int, str string) (n int, err error) {
	defer l.Skip(0)

	t := time.Now()

	prefix := fmt.Sprintf(Green, t.Format(l.timeFormatStr()))
	prefix += fmt.Sprintf(Blue, l.caller())
	if l.logType&lType > 0 {
		prefix += fmt.Sprintf(coloredStr[lType], logType[lType])
	}

	return fmt.Print(prefix + str)
}
