package factory

import (
	"strings"
)

type EncodingName string

const (
	EncodingNormal EncodingName = "console"
	EncodingJson   EncodingName = "json"
)

const (
	DTFormatNormal = "2006-01-02 15:04:05.000"
)

var loggers = make(map[string]*Logger)

type LoggerFactory struct {
	callerPackage func(caller string) string
	delegate      internalFactory
}

type LevelName string

const (
	TRACE  LevelName = "TRACE"
	DEBUG  LevelName = "DEBUG"
	INFO   LevelName = "INFO"
	WARN   LevelName = "WARN"
	ERROR  LevelName = "ERROR"
	FATAL  LevelName = "FATAL"
	DPANIC LevelName = "DPANIC"
	PANIC  LevelName = "PANIC"
)

type LevelNum int8

const (
	LvlTrace  LevelNum = -2
	LvlDebug  LevelNum = -1
	LvlInfo   LevelNum = 0
	LvlWarn   LevelNum = 1
	LvlError  LevelNum = 2
	LvlDPanic LevelNum = 3
	LvlPanic  LevelNum = 4
	LvlFatal  LevelNum = 5
)

type KeyVal struct {
	Key string
	Val interface{}
}

type internalFactory interface {
	getLevels(string) map[string]string
	setLevels(string, string)
	newLogger(string, string, []string) *Logger
}

func (f *LoggerFactory) GetLevels(prefix string) map[string]string {
	return f.delegate.getLevels(prefix)
}
func (f *LoggerFactory) SetLevels(prefix string, level string) {
	f.delegate.setLevels(prefix, level)
}

func (f *LoggerFactory) NewLogger(callerFile string, outPaths []string) *Logger {
	//fmt.Println(fmt.Sprintf("%s, %s, %s, %s", pc, file, line, ok))
	callerPackage := f.callerPackage(callerFile)
	logger := f.delegate.newLogger(callerPackage, "info", outPaths)
	logger.factory = f
	return loggers[logger.Name]
}

var ZapLoggerFactoryImpl = ZapLoggerFactory("zap")
var LogrusLoggerFactoryImpl = LogrusLoggerFactory("logrus")

func NewLoggerFactory(impl string, callerPackageDetector func(caller string) string) *LoggerFactory {
	switch strings.ToLower(impl) {
	case string(ZapLoggerFactoryImpl):
		return ZapLoggerFactoryImpl.NewFactory(callerPackageDetector)
	case string(LogrusLoggerFactoryImpl):
		return LogrusLoggerFactoryImpl.NewFactory(callerPackageDetector)
	}
	return nil
}
