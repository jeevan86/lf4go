package factory

import (
	"strings"
	"time"
)

const (
	DTFormatNormal = "2006-01-02 15:04:05.000"
)

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
	newLogger(config *LoggerConfig) *Logger
}

func (f *LoggerFactory) GetLevels(prefix string) map[string]string {
	return f.delegate.getLevels(prefix)
}
func (f *LoggerFactory) SetLevels(prefix string, level string) {
	f.delegate.setLevels(prefix, level)
}

func (f *LoggerFactory) NewLogger(
	callerFile string, formatter string, outPaths []string,
	maxFileSize int, maxFileBackups int, maxFileAge time.Duration, localTime bool, compress bool) *Logger {
	//fmt.Println(fmt.Sprintf("%s, %s, %s, %s", pc, file, line, ok))
	callerPackage := f.callerPackage(callerFile)
	writerConfig := &WriterConfig{
		MaxFileSize:    maxFileSize,
		MaxFileBackups: maxFileBackups,
		MaxFileAge:     maxFileAge,
		LocalTime:      localTime,
		Compress:       compress,
	}
	loggerConfig := &LoggerConfig{
		Name:      callerPackage,
		Level:     logLevelNum("info"),
		Formatter: formatter,
		Writer:    writerConfig,
		OutPaths:  outPaths,
	}
	logger := f.delegate.newLogger(loggerConfig)
	logger.factory = f
	loggers[logger.Config.Name] = logger
	return loggers[logger.Config.Name]
}

var ZapLoggerFactoryImpl = ZapLoggerFactory("zap")
var LogrusLoggerFactoryImpl = LogrusLoggerFactory("logrus")

func NewLoggerFactory(impl string, callerPackageDetector func(caller string) string) *LoggerFactory {
	switch strings.ToLower(impl) {
	case string(ZapLoggerFactoryImpl):
		return ZapLoggerFactoryImpl.NewFactory(callerPackageDetector)
	case string(LogrusLoggerFactoryImpl):
		return LogrusLoggerFactoryImpl.NewFactory(callerPackageDetector)
	default:
		return LogrusLoggerFactoryImpl.NewFactory(callerPackageDetector)
	}
}
