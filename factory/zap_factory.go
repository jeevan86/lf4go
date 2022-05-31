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
			levels[k] = logLevelName(logger.Level)
		}
	} else {
		for k, logger := range loggers {
			if strings.HasPrefix(k, prefix) {
				levels[k] = logLevelName(logger.Level)
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
			logger.Level = levelNum
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
			logger.Level = levelNum
			logger.delegate = delegate
		}
	}
}

func (zf *ZapLoggerFactory) setLevel(name string, level zap.AtomicLevel) *zap.Logger {
	logger := loggers[name]
	internal := logger.delegate.(*ZapLogger)
	internal.config.Level = level
	return newZapLogger(name, internal.config)
}

// newLogger
// []string{"stdout"},
// []string{"stderr"},
func (zf *ZapLoggerFactory) newLogger(name string, level string, outPaths []string) *Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(DTFormatNormal)
	atomicLevel, levelNum := zf.logLevel(level)
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
		ErrorOutputPaths: outPaths,
	}
	sink := newZapLogger(name, config)
	delegate := &ZapLogger{
		name:    name,
		config:  config,
		level:   levelNum,
		sink:    sink,
		factory: zf,
	}
	loggers[name] = &Logger{
		Name:     name,
		Level:    levelNum,
		delegate: delegate,
		outPaths: outPaths,
	}
	return loggers[name]
}

// newZapLogger
// []string{"stdout"},
// []string{"stderr"},
func newZapLogger(name string, config *zap.Config) *zap.Logger {
	var encoder zapcore.Encoder
	if string(EncodingNormal) == config.Encoding {
		encoder = zapcore.NewConsoleEncoder(config.EncoderConfig)
	} else if string(EncodingJson) == config.Encoding {
		encoder = zapcore.NewJSONEncoder(config.EncoderConfig)
	}
	log := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer(name, config.OutputPaths)), config.Level),
	)
	delegate := log.WithOptions(zap.AddCallerSkip(3))
	//delegate, _ := config.Build(zap.AddCallerSkip(3))
	return delegate
}

func (lf *ZapLoggerFactory) logLevel(level string) (zap.AtomicLevel, LevelNum) {
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
