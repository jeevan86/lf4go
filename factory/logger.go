package factory

import (
	"strings"
)

type Logger struct {
	Name     string
	Level    LevelNum
	delegate loggerDelegate
	outPaths []string
	factory  *LoggerFactory
}

type loggerDelegate interface {
	Trace(string, ...interface{})
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	DPanic(string, ...interface{})
	Panic(string, ...interface{})
}

func (l *Logger) SetLevels(prefix string, level string) {
	l.factory.SetLevels(prefix, level)
}
func (l *Logger) GetLevels(prefix string) map[string]string {
	return l.factory.GetLevels(prefix)
}

func (l *Logger) IsTraceEnabled() bool {
	return l.Level <= LvlTrace
}
func (l *Logger) IsDebugEnabled() bool {
	return l.Level <= LvlDebug
}
func (l *Logger) IsInfoEnabled() bool {
	return l.Level <= LvlInfo
}
func (l *Logger) IsWarnEnabled() bool {
	return l.Level <= LvlWarn
}
func (l *Logger) IsErrorEnabled() bool {
	return l.Level <= LvlError
}
func (l *Logger) IsDPanicEnabled() bool {
	return l.Level <= LvlDPanic
}
func (l *Logger) IsPanicEnabled() bool {
	return l.Level <= LvlPanic
}
func (l *Logger) IsFatalEnabled() bool {
	return l.Level <= LvlFatal
}

func (l *Logger) Trace(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.delegate.Trace(msg)
	}
	l.delegate.Trace(msg, kvs)
}
func (l *Logger) Debug(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.delegate.Debug(msg)
	}
	l.delegate.Debug(msg, kvs)
}
func (l *Logger) Info(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.delegate.Info(msg)
	}
	l.delegate.Info(msg, kvs)
}
func (l *Logger) Warn(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.delegate.Warn(msg)
	}
	l.delegate.Warn(msg, kvs)
}
func (l *Logger) Error(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.delegate.Error(msg)
	}
	l.delegate.Error(msg, kvs)
}
func (l *Logger) DPanic(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.delegate.DPanic(msg)
	}
	l.delegate.DPanic(msg, kvs)
}
func (l *Logger) Panic(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.delegate.Panic(msg)
	}
	l.delegate.Panic(msg, kvs)
}
func (l *Logger) Fatal(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.delegate.Fatal(msg)
	}
	l.delegate.Fatal(msg, kvs)
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
