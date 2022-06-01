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

func (l *LogrusLogger) Trace(msg string) {
	//l.log(msg, l.sink.Trace)
	l.sink.Trace(msg)
}
func (l *LogrusLogger) Debug(msg string) {
	//l.log(msg, l.sink.Debug)
	l.sink.Debug(msg)
}
func (l *LogrusLogger) Info(msg string) {
	//l.log(msg, l.sink.Info)
	l.sink.Info(msg)
}
func (l *LogrusLogger) Warn(msg string) {
	// l.log(msg, l.sink.Warn)
	l.sink.Warn(msg)
}
func (l *LogrusLogger) Error(msg string) {
	//l.log(msg, l.sink.Error)
	l.sink.Error(msg)
}
func (l *LogrusLogger) Fatal(msg string) {
	//l.log(msg, l.sink.Fatal)
	l.sink.Fatal(msg)
}
func (l *LogrusLogger) DPanic(msg string) {
	//l.log(msg, l.sink.Panic)
	l.sink.Panic(msg)
}
func (l *LogrusLogger) Panic(msg string) {
	//l.log(msg, l.sink.Panic)
	l.sink.Panic(msg)
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
