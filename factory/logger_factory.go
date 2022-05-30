package factory

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"strings"
)

type LevelName string

const (
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
	LvlDebug  LevelNum = -1
	LvlInfo   LevelNum = 0
	LvlWarn   LevelNum = 1
	LvlError  LevelNum = 2
	LvlDPanic LevelNum = 3
	LvlPanic  LevelNum = 4
	LvlFatal  LevelNum = 5
)

type Logging struct {
	RootName      string            `yaml:"root-name"`
	RootLevel     string            `yaml:"root-level"`
	PackageLevels map[string]string `yaml:"package-levels"`
	Encoder       string            `yaml:"encoder"`
	LogFileDir    string            `yaml:"log-file-dir"`
}

type Logger struct {
	Name     string
	Level    LevelNum
	Config   *zap.Config
	Delegate *zap.Logger
	OutPaths []string
	ErrPaths []string
	factory  *LoggerFactory
}

type KeyVal struct {
	Key string
	Val interface{}
}

func fields(kvs ...KeyVal) []zap.Field {
	if kvs != nil {
		sz := len(kvs)
		if sz > 0 {
			fields := make([]zap.Field, sz, sz)
			for i, kv := range kvs {
				field := zap.Field{}
				zap.Any(kv.Key, kv.Val)
				fields[i] = field
			}
			return fields
		}
	}
	return nil
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

func (l *Logger) Debug(msg string, kvs ...KeyVal) {
	if kvs == nil {
		l.Delegate.Debug(msg)
		return
	}
	fields := fields(kvs...)
	if fields == nil {
		l.Delegate.Debug(msg)
		return
	}
	l.Delegate.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, kvs ...KeyVal) {
	if kvs == nil {
		l.Delegate.Info(msg)
		return
	}
	fields := fields(kvs...)
	if fields == nil {
		l.Delegate.Info(msg)
		return
	}
	l.Delegate.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, kvs ...KeyVal) {
	if kvs == nil {
		l.Delegate.Warn(msg)
		return
	}
	fields := fields(kvs...)
	if fields == nil {
		l.Delegate.Warn(msg)
		return
	}
	l.Delegate.Warn(msg, fields...)
}
func (l *Logger) Error(msg string, kvs ...KeyVal) {
	if kvs == nil {
		l.Delegate.Error(msg)
		return
	}
	fields := fields(kvs...)
	if fields == nil {
		l.Delegate.Error(msg)
		return
	}
	l.Delegate.Error(msg, fields...)
}
func (l *Logger) DPanic(msg string, kvs ...KeyVal) {
	if kvs == nil {
		l.Delegate.DPanic(msg)
		return
	}
	fields := fields(kvs...)
	if fields == nil {
		l.Delegate.DPanic(msg)
		return
	}
	l.Delegate.DPanic(msg, fields...)
}
func (l *Logger) Panic(msg string, kvs ...KeyVal) {
	if kvs == nil {
		l.Delegate.Panic(msg)
		return
	}
	fields := fields(kvs...)
	if fields == nil {
		l.Delegate.Debug(msg)
		return
	}
	l.Delegate.Panic(msg, fields...)
}
func (l *Logger) Fatal(msg string, kvs ...KeyVal) {
	if kvs == nil {
		l.Delegate.Fatal(msg)
		return
	}
	fields := fields(kvs...)
	if fields == nil {
		l.Delegate.Fatal(msg)
		return
	}
	l.Delegate.Fatal(msg, fields...)
}

func (l *Logger) SetLevel(level string) {
	var atomicLevel zap.AtomicLevel
	var levelNum LevelNum
	atomicLevel, levelNum = logLevel(level)
	l.Delegate = l.factory.setLevel(l.Name, atomicLevel)
	l.Level = levelNum
}

func (l *Logger) GetLevels(prefix string) map[string]string {
	return l.factory.GetLevels(prefix)
}

func (l *Logger) SetLevels(prefix string, level string) {
	l.factory.SetLevels(prefix, level)
}

type LoggerFactory struct {
	loggers       map[string]*Logger
	callerPackage func(caller string) string
}

func NewLoggerFactory(callerPackageDetector func(caller string) string) *LoggerFactory {
	factory := &LoggerFactory{
		loggers:       make(map[string]*Logger, 32),
		callerPackage: callerPackageDetector,
	}
	return factory
}

func (logging *LoggerFactory) setLevel(name string, level zap.AtomicLevel) *zap.Logger {
	logger := logging.loggers[name]
	config := logger.Config
	config.Level = level
	return newZapLogger(config)
}

func (logging *LoggerFactory) GetLevels(prefix string) map[string]string {
	levels := make(map[string]string, 16)
	if "ROOT" == strings.ToUpper(prefix) {
		for k, logger := range logging.loggers {
			levels[k] = logLevelName(logger.Level)
		}
	} else {
		for k, logger := range logging.loggers {
			if strings.HasPrefix(k, prefix) {
				levels[k] = logLevelName(logger.Level)
			}
		}
	}
	return levels
}

func (logging *LoggerFactory) SetLevels(prefix string, level string) {
	if "ROOT" == strings.ToUpper(prefix) {
		for k, logger := range logging.loggers {
			var atomicLevel zap.AtomicLevel
			var levelNum LevelNum
			atomicLevel, levelNum = logLevel(level)
			delegate := logging.setLevel(k, atomicLevel)
			logger.Level = levelNum
			logger.Delegate = delegate
		}
		return
	}
	for k, logger := range logging.loggers {
		if strings.HasPrefix(k, prefix) {
			var atomicLevel zap.AtomicLevel
			var levelNum LevelNum
			atomicLevel, levelNum = logLevel(level)
			delegate := logging.setLevel(k, atomicLevel)
			logger.Level = levelNum
			logger.Delegate = delegate
		}
	}
}

func (logging *LoggerFactory) NewLogger(outPaths []string, errPaths []string) *Logger {
	//fmt.Println(fmt.Sprintf("%s, %s, %s, %s", pc, file, line, ok))
	logger := newLogger("info", outPaths, errPaths)
	_, f, _, _ := runtime.Caller(1)
	callerPackage := logging.callerPackage(f)
	logger.Name = callerPackage
	logger.factory = logging
	logging.loggers[logger.Name] = logger
	return logging.loggers[logger.Name]
}

// newLogger
// []string{"stdout"},
// []string{"stderr"},
func newLogger(level string, outPaths []string, errPaths []string) *Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(DTFormatNormal)
	atomicLevel, levelNum := logLevel(level)
	config := &zap.Config{
		Level:       atomicLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         string(EncodingNormal),
		EncoderConfig:    encoderConfig,
		OutputPaths:      outPaths,
		ErrorOutputPaths: errPaths,
	}
	delegate := newZapLogger(config)
	return &Logger{
		Config:   config,
		Level:    levelNum,
		Delegate: delegate,
		OutPaths: outPaths,
		ErrPaths: errPaths,
	}
}

// newZapLogger
// []string{"stdout"},
// []string{"stderr"},
func newZapLogger(config *zap.Config) *zap.Logger {
	delegate, _ := config.Build(zap.AddCallerSkip(2))
	return delegate
}

type EncodingName string

const (
	EncodingNormal EncodingName = "console"
	EncodingJson   EncodingName = "json"
)

const (
	DTFormatNormal = "2006-01-02 15:04:05.000"
)

func logLevel(level string) (zap.AtomicLevel, LevelNum) {
	atomicLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	var levelNum = LvlInfo
	switch strings.ToUpper(level) {
	case "DEBUG":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		levelNum = LvlDebug
		break
	case "INFO":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		levelNum = LvlInfo
		break
	case "WARN":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		levelNum = LvlWarn
		break
	case "ERROR":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		levelNum = LvlError
		break
	case "DPANIC":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.DPanicLevel)
		levelNum = LvlDPanic
		break
	case "PANIC":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.PanicLevel)
		levelNum = LvlPanic
		break
	case "FATAL":
		atomicLevel = zap.NewAtomicLevelAt(zapcore.FatalLevel)
		levelNum = LvlFatal
		break
	}
	return atomicLevel, levelNum
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
