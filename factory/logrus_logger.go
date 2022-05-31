package factory

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type LogrusLogger struct {
	name    string
	level   LevelNum
	sink    *logrus.Logger
	factory *LogrusLoggerFactory
}

func (l *LogrusLogger) SetLevel(level string) {
	var levelObj logrus.Level
	var levelNum LevelNum
	levelObj, levelNum = l.factory.logLevel(level)
	l.sink = l.factory.setLevel(l.name, levelObj)
	l.level = levelNum
}

func (l *LogrusLogger) Trace(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Trace)
		return
	}
	l.sink.Tracef(msg, kvs)
}
func (l *LogrusLogger) Debug(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Debug)
		return
	}
	l.sink.Debugf(msg, kvs)
}
func (l *LogrusLogger) Info(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Info)
		return
	}
	l.sink.Infof(msg, kvs)
}
func (l *LogrusLogger) Warn(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Warn)
		return
	}
	l.sink.Warnf(msg, kvs)
}
func (l *LogrusLogger) Error(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Error)
		return
	}
	l.sink.Errorf(msg, kvs)
}
func (l *LogrusLogger) Fatal(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Fatal)
		return
	}
	l.sink.Fatalf(msg, kvs)
}
func (l *LogrusLogger) DPanic(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Panic)
		return
	}
	l.sink.Panicf(msg, kvs)
}
func (l *LogrusLogger) Panic(msg string, kvs ...interface{}) {
	if kvs == nil || len(kvs) == 0 {
		l.log(msg, l.sink.Panic)
		return
	}
	l.sink.Panicf(msg, kvs)
}

func (l *LogrusLogger) log(msg string, f func(...interface{}), args ...interface{}) {
	if args == nil || len(args) == 0 {
		f(msg)
		return
	}
	appended := make([]interface{}, len(args)+1)
	appended[0] = msg
	for i, arg := range args {
		appended[i+1] = arg
	}
	f(appended...)
}

func (l *LogrusLogger) convert(elements ...interface{}) []interface{} {
	if elements == nil || len(elements) == 0 {
		return nil
	}
	args := make([]interface{}, len(elements))
	for i, e := range elements {
		if reflect.TypeOf(e).Name() == "KeyVal" {
			kv := e.(KeyVal)
			arg := kv.Key + ":" + fmt.Sprint(kv.Val)
			args[i] = arg
		} else {
			args[i] = e
		}
	}
	return args
}
