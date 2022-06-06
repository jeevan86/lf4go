package factory

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type ZapLoggerFactory string

func (zf *ZapLoggerFactory) NewFactory(callerPackageDetector func(caller string) string) *LoggerFactory {
	internalFactory := ZapLoggerFactoryImpl
	factory := &LoggerFactory{
		callerPackage: callerPackageDetector,
		delegate:      &internalFactory,
	}
	return factory
}

func (zf *ZapLoggerFactory) getLevels(prefix string) map[string]string {
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

func (zf *ZapLoggerFactory) setLevels(prefix string, level string) {
	if "ROOT" == strings.ToUpper(prefix) {
		for k, logger := range loggers {
			var levelObj zap.AtomicLevel
			var levelNum LevelNum
			levelObj, levelNum = zf.logLevel(level)
			sink := zf.setLevel(k, levelObj)
			delegate := &ZapLogger{
				sink:    sink,
				factory: zf,
			}
			logger.Config.Level = levelNum
			logger.delegate = delegate
		}
		return
	}
	for k, logger := range loggers {
		if strings.HasPrefix(k, prefix) {
			var levelObj zap.AtomicLevel
			var levelNum LevelNum
			levelObj, levelNum = zf.logLevel(level)
			sink := zf.setLevel(k, levelObj)
			delegate := &ZapLogger{
				sink:    sink,
				factory: zf,
			}
			logger.Config.Level = levelNum
			logger.delegate = delegate
		}
	}
}

func (zf *ZapLoggerFactory) setLevel(name string, level zap.AtomicLevel) *zap.Logger {
	logger := loggers[name]
	loggerConfig := logger.Config
	internal := logger.delegate.(*ZapLogger)
	internal.config.Level = level
	return newZapLogger(loggerConfig, internal.config)
}

// newLogger
// []string{"stdout"},
// []string{"stderr"},
func (zf *ZapLoggerFactory) newLogger(loggerConfig *LoggerConfig) *Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(DTFormatNormal)
	atomicLevel, _ := zf.logLevel(logLevelName(loggerConfig.Level))
	encoding := zf.formatterToEncoding(loggerConfig.Formatter)
	config := &zap.Config{
		Level:       atomicLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:      encoding,
		EncoderConfig: encoderConfig,
	}
	sink := newZapLogger(loggerConfig, config)
	delegate := &ZapLogger{
		config:  config,
		sink:    sink,
		factory: zf,
	}
	return &Logger{
		Config:   loggerConfig,
		delegate: delegate,
	}
}

func (zf *ZapLoggerFactory) formatterToEncoding(formatter string) string {
	encoding := strings.ToLower(formatter)
	if encoding == "normal" {
		encoding = "console"
	} else if encoding == "json" {
		encoding = "json"
	} else {
		encoding = "console"
	}
	return encoding
}

type zapEncodingName string

const (
	zapEncodingNormal zapEncodingName = "console"
	zapEncodingJson   zapEncodingName = "json"
)

// newZapLogger
// []string{"stdout"},
// []string{"stderr"},
func newZapLogger(loggerConfig *LoggerConfig, config *zap.Config) *zap.Logger {
	var encoder zapcore.Encoder
	if string(zapEncodingNormal) == config.Encoding {
		encoder = zapcore.NewConsoleEncoder(config.EncoderConfig)
	} else if string(zapEncodingJson) == config.Encoding {
		encoder = zapcore.NewJSONEncoder(config.EncoderConfig)
	}
	log := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer(loggerConfig.Name, loggerConfig.Appenders)), config.Level),
	)
	delegate := log.WithOptions(
		zap.AddCallerSkip(3),
		zap.AddCaller(),
	)
	//delegate, _ := config.Build(zap.AddCallerSkip(3))
	return delegate
}

func (zf *ZapLoggerFactory) logLevel(level string) (zap.AtomicLevel, LevelNum) {
	var levelObj = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	var levelNum = LvlInfo
	switch strings.ToUpper(level) {
	case "TRACE":
		levelObj = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		levelNum = LvlDebug
		break
	case "DEBUG":
		levelObj = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		levelNum = LvlDebug
		break
	case "INFO":
		levelObj = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		levelNum = LvlInfo
		break
	case "WARN":
		levelObj = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		levelNum = LvlWarn
		break
	case "ERROR":
		levelObj = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		levelNum = LvlError
		break
	case "DPANIC":
		levelObj = zap.NewAtomicLevelAt(zapcore.DPanicLevel)
		levelNum = LvlDPanic
		break
	case "PANIC":
		levelObj = zap.NewAtomicLevelAt(zapcore.PanicLevel)
		levelNum = LvlPanic
		break
	case "FATAL":
		levelObj = zap.NewAtomicLevelAt(zapcore.FatalLevel)
		levelNum = LvlFatal
		break
	}
	return levelObj, levelNum
}
