package factory

import (
	"fmt"
	"runtime"
	"strings"
)

var loggers = make(map[string]*Logger)

type Logger struct {
	Config   *LoggerConfig
	delegate loggerDelegate
	factory  *LoggerFactory
}

type LoggerConfig struct {
	Name         string
	Level        LevelNum
	Formatter    string
	Appenders    []AppenderConfig
	ReportCaller bool
}

type loggerDelegate interface {
	Trace(string)
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	DPanic(string)
	Panic(string)
}

func (l *Logger) SetLevels(prefix string, level string) {
	l.factory.SetLevels(prefix, level)
}
func (l *Logger) GetLevels(prefix string) map[string]string {
	return l.factory.GetLevels(prefix)
}

func (l *Logger) IsTraceEnabled() bool {
	return l.Config.Level <= LvlTrace
}
func (l *Logger) IsDebugEnabled() bool {
	return l.Config.Level <= LvlDebug
}
func (l *Logger) IsInfoEnabled() bool {
	return l.Config.Level <= LvlInfo
}
func (l *Logger) IsWarnEnabled() bool {
	return l.Config.Level <= LvlWarn
}
func (l *Logger) IsErrorEnabled() bool {
	return l.Config.Level <= LvlError
}
func (l *Logger) IsDPanicEnabled() bool {
	return l.Config.Level <= LvlDPanic
}
func (l *Logger) IsPanicEnabled() bool {
	return l.Config.Level <= LvlPanic
}
func (l *Logger) IsFatalEnabled() bool {
	return l.Config.Level <= LvlFatal
}

func (l *Logger) Trace(format string, args ...interface{}) {
	l.delegate.Trace(l.doFormat(format, args...))
}
func (l *Logger) Debug(format string, args ...interface{}) {
	l.delegate.Debug(l.doFormat(format, args...))
}
func (l *Logger) Info(format string, args ...interface{}) {
	l.delegate.Info(l.doFormat(format, args...))
}
func (l *Logger) Warn(format string, args ...interface{}) {
	l.delegate.Warn(l.doFormat(format, args...))
}
func (l *Logger) Error(format string, args ...interface{}) {
	l.delegate.Error(l.doFormat(format, args...))
}
func (l *Logger) DPanic(format string, args ...interface{}) {
	l.delegate.DPanic(l.doFormat(format, args...))
}
func (l *Logger) Panic(format string, args ...interface{}) {
	l.delegate.Panic(l.doFormat(format, args...))
}
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.delegate.Fatal(l.doFormat(format, args...))
}

func (l *Logger) doFormat(format string, args ...interface{}) string {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	if l.Config.ReportCaller {
		msg = l.withCaller(3, msg)
	}
	return msg
}

func (l *Logger) SkTrace(skip int, format string, args ...interface{}) {
	l.delegate.Trace(l.skDoFormat(skip, format, args...))
}
func (l *Logger) SkDebug(skip int, format string, args ...interface{}) {
	l.delegate.Debug(l.skDoFormat(skip, format, args...))
}
func (l *Logger) SkInfo(skip int, format string, args ...interface{}) {
	l.delegate.Info(l.skDoFormat(skip, format, args...))
}
func (l *Logger) SkWarn(skip int, format string, args ...interface{}) {
	l.delegate.Warn(l.skDoFormat(skip, format, args...))
}
func (l *Logger) SkError(skip int, format string, args ...interface{}) {
	l.delegate.Error(l.skDoFormat(skip, format, args...))
}
func (l *Logger) SkDPanic(skip int, format string, args ...interface{}) {
	l.delegate.DPanic(l.skDoFormat(skip, format, args...))
}
func (l *Logger) SkPanic(skip int, format string, args ...interface{}) {
	l.delegate.Panic(l.skDoFormat(skip, format, args...))
}
func (l *Logger) SkFatal(skip int, format string, args ...interface{}) {
	l.delegate.Fatal(l.skDoFormat(skip, format, args...))
}

func (l *Logger) skDoFormat(skip int, format string, args ...interface{}) string {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	if l.Config.ReportCaller {
		msg = l.withCaller(skip, msg)
	}
	return msg
}

func (l *Logger) withCaller(skip int, format string) string {
	pc, file, line, _ := runtime.Caller(skip)
	funcName := stringAfterLast(runtime.FuncForPC(pc).Name(), SLASH)
	packName := l.Config.Name
	fileName := stringAfterLast(file, SLASH)
	fileLine := line
	return fmt.Sprintf("%s(%s/%s:%d)\t%s", funcName, packName, fileName, fileLine, format)
}

func stringAfterLast(origin, last string) string {
	idx := strings.LastIndex(origin, last)
	if idx == -1 {
		return origin
	}
	if len(origin) <= idx+1 {
		return ""
	}
	return origin[idx+1:]
}

func logLevelNum(level string) LevelNum {
	var levelNum = LvlInfo
	switch strings.ToUpper(level) {
	case "TRACE":
		levelNum = LvlTrace
		break
	case "DEBUG":
		levelNum = LvlDebug
		break
	case "INFO":
		levelNum = LvlInfo
		break
	case "WARN":
		levelNum = LvlWarn
		break
	case "ERROR":
		levelNum = LvlError
		break
	case "DPANIC":
		levelNum = LvlDPanic
		break
	case "PANIC":
		levelNum = LvlPanic
		break
	case "FATAL":
		levelNum = LvlFatal
		break
	}
	return levelNum
}

func logLevelName(num LevelNum) string {
	name := "Info"
	switch num {
	case LvlTrace:
		name = "Trace"
		break
	case LvlDebug:
		name = "Debug"
		break
	case LvlInfo:
		name = "Info"
		break
	case LvlWarn:
		name = "Warn"
		break
	case LvlError:
		name = "Error"
		break
	case LvlDPanic:
		name = "DPanic"
		break
	case LvlPanic:
		name = "Panic"
		break
	case LvlFatal:
		name = "Fatal"
		break
	}
	return name
}
