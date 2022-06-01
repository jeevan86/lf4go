package factory

import (
	"go.uber.org/zap"
	"reflect"
)

type ZapLogger struct {
	config  *zap.Config
	sink    *zap.Logger
	factory *ZapLoggerFactory
}

func (l *ZapLogger) Trace(msg string) {
	//l.log(msg, l.sink.Debug)
	l.sink.Debug(msg)
}
func (l *ZapLogger) Debug(msg string) {
	//l.log(msg, l.sink.Debug)
	l.sink.Debug(msg)
}
func (l *ZapLogger) Info(msg string) {
	//l.log(msg, l.sink.Info)
	l.sink.Info(msg)
}
func (l *ZapLogger) Warn(msg string) {
	// l.log(msg, l.sink.Warn)
	l.sink.Warn(msg)
}
func (l *ZapLogger) Error(msg string) {
	//l.log(msg, l.sink.Error)
	l.sink.Error(msg)
}
func (l *ZapLogger) Fatal(msg string) {
	//l.log(msg, l.sink.Fatal)
	l.sink.Fatal(msg)
}
func (l *ZapLogger) DPanic(msg string) {
	//l.log(msg, l.sink.DPanic)
	l.sink.DPanic(msg)
}
func (l *ZapLogger) Panic(msg string) {
	//l.log(msg, l.sink.Panic)
	l.sink.Panic(msg)
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
