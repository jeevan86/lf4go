package factory

import (
	"fmt"
	"strings"
)

var loggers = make(map[string]*Logger)

type Logger struct {
	Config   *LoggerConfig
	delegate loggerDelegate
	factory  *LoggerFactory
}

type LoggerConfig struct {
	Name      string
	Level     LevelNum
	Formatter string
	Appenders []AppenderConfig
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
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.delegate.Trace(msg)
}
func (l *Logger) Debug(format string, args ...interface{}) {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.delegate.Debug(msg)
}
func (l *Logger) Info(format string, args ...interface{}) {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.delegate.Info(msg)
}
func (l *Logger) Warn(format string, args ...interface{}) {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.delegate.Warn(msg)
}
func (l *Logger) Error(format string, args ...interface{}) {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.delegate.Error(msg)
}
func (l *Logger) DPanic(format string, args ...interface{}) {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.delegate.DPanic(msg)
}
func (l *Logger) Panic(format string, args ...interface{}) {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.delegate.Panic(msg)
}
func (l *Logger) Fatal(format string, args ...interface{}) {
	msg := format
	if args != nil && len(args) != 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.delegate.Fatal(msg)
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
