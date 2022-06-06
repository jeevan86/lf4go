package factory

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"strings"
)

type LogrusLoggerFactory string

func (lf *LogrusLoggerFactory) NewFactory(callerPackageDetector func(caller string) string) *LoggerFactory {
	internalFactory := LogrusLoggerFactoryImpl
	factory := &LoggerFactory{
		callerPackage: callerPackageDetector,
		delegate:      &internalFactory,
	}
	return factory
}

func (lf *LogrusLoggerFactory) fields(elements ...interface{}) []interface{} {
	if elements != nil {
		sz := len(elements)
		if sz > 0 {
			fields := make([]interface{}, sz, sz)
			for i, e := range elements {
				if reflect.TypeOf(e).Name() == "KeyVal" {
					kv := e.(KeyVal)
					field := kv.Key + ":" + fmt.Sprint(kv.Val)
					fields[i] = field
				} else {
					fields[i] = e
				}
			}
			return fields
		}
	}
	return nil
}

func (lf *LogrusLoggerFactory) setLevel(name string, level logrus.Level) *logrus.Logger {
	logger := loggers[name]
	return lf.newLogrusLogger(logger.Config, level)
}

func (lf *LogrusLoggerFactory) getLevels(prefix string) map[string]string {
	levels := make(map[string]string, 16)
	if "ROOT" == strings.ToUpper(prefix) {
		for k, logger := range loggers {
			levels[k] = logLevelName(logger.Config.Level)
		}
	} else {
		for k, logger := range loggers {
			if strings.HasPrefix(k, prefix) {
				levels[k] = logLevelName(logger.Config.Level)
			}
		}
	}
	return levels
}

func (lf *LogrusLoggerFactory) setLevels(prefix string, level string) {
	if "ROOT" == strings.ToUpper(prefix) {
		for _, logger := range loggers {
			lf.setLoggerLevel(logger, level)
		}
		return
	}
	for k, logger := range loggers {
		if strings.HasPrefix(k, prefix) {
			lf.setLoggerLevel(logger, level)
		}
	}
}

func (lf *LogrusLoggerFactory) setLoggerLevel(logger *Logger, level string) {
	var logrusLevel logrus.Level
	var levelNum LevelNum
	logrusLevel, levelNum = lf.logLevel(level)
	sink := lf.newLogrusLogger(logger.Config, logrusLevel)
	delegate := &LogrusLogger{
		sink:    sink,
		factory: lf,
	}
	logger.Config.Level = levelNum
	logger.delegate = delegate
}

// newLogger
// []string{"stdout", "logs/application.log"},
func (lf *LogrusLoggerFactory) newLogger(loggerConfig *LoggerConfig) *Logger {
	logrusLevel, _ := lf.logLevel(logLevelName(loggerConfig.Level))
	sink := lf.newLogrusLogger(loggerConfig, logrusLevel)
	delegate := &LogrusLogger{
		sink:    sink,
		factory: lf,
	}
	return &Logger{
		Config:   loggerConfig,
		delegate: delegate,
	}
}

type logrusHook string

func (h logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

const hookFireFuncSkip = 10 // function Fire may invoke 10th order.
const slash = "/"

func (h logrusHook) Fire(entry *logrus.Entry) error {
	pc, file, line, _ := runtime.Caller(hookFireFuncSkip)
	lastStashIdx := strings.LastIndex(file, slash)
	if lastStashIdx >= 0 {
		file = string(h) + ":" + file[lastStashIdx+1:]
	} else {
		file = string(h) + ":" + file
	}
	entry.Caller = &runtime.Frame{
		PC:   pc,
		File: file,
		Line: line,
	}
	return nil
}

func (lf *LogrusLoggerFactory) newHook(name string) logrus.LevelHooks {
	var allLevelHook = logrusHook(name)
	var allLevelHooks = []logrus.Hook{allLevelHook}
	return logrus.LevelHooks{
		logrus.TraceLevel: allLevelHooks,
		logrus.DebugLevel: allLevelHooks,
		logrus.InfoLevel:  allLevelHooks,
		logrus.WarnLevel:  allLevelHooks,
		logrus.ErrorLevel: allLevelHooks,
		logrus.FatalLevel: allLevelHooks,
		logrus.PanicLevel: allLevelHooks,
	}
}

// newLogrusLogger
// []string{"stdout", "logs/application.log"},
func (lf *LogrusLoggerFactory) newLogrusLogger(loggerConfig *LoggerConfig, level logrus.Level) *logrus.Logger {
	merged := writer(loggerConfig.Name, loggerConfig.Appenders)
	delegate := &logrus.Logger{
		Out:          merged,
		Hooks:        lf.newHook(loggerConfig.Name),
		Formatter:    logrusFormatter(loggerConfig.Formatter),
		ReportCaller: true, // set to false will cause entry.HasCaller() return false, wtf!
		Level:        level,
		// ExitFunc exitFunc, // Function to exit the application, defaults to `os.Exit()`
	}
	return delegate
}

func (lf *LogrusLoggerFactory) logLevel(level string) (logrus.Level, LevelNum) {
	var logrusLevel logrus.Level
	var levelNum = LvlInfo
	switch strings.ToUpper(level) {
	case "TRACE":
		logrusLevel = logrus.TraceLevel
		levelNum = LvlDebug
		break
	case "DEBUG":
		logrusLevel = logrus.DebugLevel
		levelNum = LvlDebug
		break
	case "INFO":
		logrusLevel = logrus.InfoLevel
		levelNum = LvlInfo
		break
	case "WARN":
		logrusLevel = logrus.WarnLevel
		levelNum = LvlWarn
		break
	case "ERROR":
		logrusLevel = logrus.ErrorLevel
		levelNum = LvlError
		break
	case "DPANIC":
		logrusLevel = logrus.PanicLevel
		levelNum = LvlDPanic
		break
	case "PANIC":
		logrusLevel = logrus.PanicLevel
		levelNum = LvlPanic
		break
	case "FATAL":
		logrusLevel = logrus.FatalLevel
		levelNum = LvlFatal
		break
	}
	return logrusLevel, levelNum
}
