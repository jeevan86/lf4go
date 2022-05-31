package factory

import (
	"go.uber.org/zap"
	"reflect"
)

type ZapLogger struct {
	name    string
	level   LevelNum
	config  *zap.Config
	sink    *zap.Logger
	factory *ZapLoggerFactory
}

func (l *ZapLogger) SetLevel(level string) {
	var atomicLevel zap.AtomicLevel
	var levelNum LevelNum
	atomicLevel, levelNum = l.factory.logLevel(level)
	l.sink = l.factory.setLevel(l.name, atomicLevel)
	l.level = levelNum
}

func (l *ZapLogger) Trace(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Debug)
		return
	}
	l.log(msg, l.sink.Debug, kvs)
}
func (l *ZapLogger) Debug(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Debug)
		return
	}
	l.log(msg, l.sink.Debug, kvs)
}
func (l *ZapLogger) Info(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Info)
		return
	}
	l.log(msg, l.sink.Info, kvs)
}
func (l *ZapLogger) Warn(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Warn)
		return
	}
	l.log(msg, l.sink.Warn, kvs)
}
func (l *ZapLogger) Error(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Error)
		return
	}
	l.log(msg, l.sink.Error, kvs)
}
func (l *ZapLogger) Fatal(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Fatal)
		return
	}
	l.log(msg, l.sink.Fatal, kvs)
}
func (l *ZapLogger) DPanic(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.DPanic)
		return
	}
	l.log(msg, l.sink.DPanic, kvs)
}
func (l *ZapLogger) Panic(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Panic)
		return
	}
	l.log(msg, l.sink.Panic, kvs)
}

func (l *ZapLogger) log(msg string, f func(string, ...zap.Field), kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		f(msg)
		return
	}
	fields := l.convert(kvs)
	if fields == nil {
		f(msg)
		return
	}
	f(msg, fields...)
}

func (l *ZapLogger) convert(elements ...interface{}) []zap.Field {
	if elements == nil || len(elements) == 0 {
		return nil
	}
	fields := make([]zap.Field, len(elements))
	for i, o := range elements {
		if reflect.TypeOf(o).Name() == "KeyVal" {
			kv := o.(KeyVal)
			fields[i] = zap.Any(kv.Key, kv.Val)
		} else {
			fields[i] = zap.Any("", o)
		}
	}
	return fields
}
