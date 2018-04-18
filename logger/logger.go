package logger

import (
	"os"
	"time"

	"github.com/gwaylib/log/logger/proto"
)

var (
	DefaultContext = proto.Context{"default", "1.0.0", HostName}
)

// logger for private
type Logger struct {
	context    *proto.Context
	adapters   []Adapter
	loggerName string
	level      proto.Level
}

func New(loggerName string, adapter ...Adapter) *Logger {
	lg := &Logger{
		adapters:   adapter,
		loggerName: loggerName,
	}
	lg.SetLevel(0)
	return lg
}

func (l *Logger) AddAdapter(adapter ...Adapter) {
	l.adapters = append(l.adapters, adapter...)
}

func (l *Logger) SetAdapter(adapter ...Adapter) {
	l.adapters = adapter
}

func (l *Logger) SetLevel(level proto.Level) {
	l.level = level
}

// 设置服务器信息
// 不设置时，使用默认配置信息
func (l *Logger) SetContext(c *proto.Context) {
	l.context = c
}

func (l *Logger) Debug(msg ...interface{}) {
	l.put(proto.LevelDebug, proto.ToMsg(msg...))
}

func (l *Logger) Debugf(f string, msg ...interface{}) {
	l.put(proto.LevelDebug, proto.ToMsgf(f, msg...))
}

func (l *Logger) Info(msg ...interface{}) {
	l.put(proto.LevelInfo, proto.ToMsg(msg...))
}

func (l *Logger) Infof(f string, msg ...interface{}) {
	l.put(proto.LevelInfo, proto.ToMsgf(f, msg...))
}

func (l *Logger) Warn(msg ...interface{}) {
	l.put(proto.LevelWarn, proto.ToMsg(msg...))
}

func (l *Logger) Warnf(f string, msg ...interface{}) {
	l.put(proto.LevelWarn, proto.ToMsgf(f, msg...))
}

func (l *Logger) Error(msg ...interface{}) {
	l.put(proto.LevelError, proto.ToMsg(msg...))
}

func (l *Logger) Errorf(f string, msg ...interface{}) {
	l.put(proto.LevelError, proto.ToMsgf(f, msg...))
}

func (l *Logger) Fatal(msg ...interface{}) {
	m := proto.ToMsg(msg...)
	l.put(proto.LevelFatal, m)
	l.Close()
	panic(m)
}

func (l *Logger) Fatalf(f string, msg ...interface{}) {
	m := proto.ToMsgf(f, msg...)
	l.put(proto.LevelFatal, m)
	l.Close()
	panic(m)
}

// Exit
// log an info level message, and close log, then call os.Exit(code)
//
// Param
// code -- code of os exit
// msg -- exit message
func (l *Logger) Exit(code int, msg ...interface{}) {
	m := proto.ToMsg(msg...)
	l.put(proto.LevelInfo, m)
	l.Close()
	os.Exit(code)

}

func (l *Logger) put(level proto.Level, msg []byte) {
	// ignore
	if level < l.level {
		return
	}
	l.Put([]*proto.Data{&proto.Data{time.Now(), level, l.loggerName, msg}})
}

func (l *Logger) Put(data []*proto.Data) {
	if len(data) == 0 {
		panic("len(data) == 0")
	}

	// create log protocal
	var context *proto.Context
	if l.context == nil {
		context = &DefaultContext
	} else {
		context = l.context
	}

	log := &proto.Proto{
		*context,
		data,
	}

	// send log
	for _, adapter := range l.adapters {
		adapter.Put(log)
	}
}

// Close logger
// it is over time, suggest to call Exit because when logger closed, program also done.
func (l *Logger) Close() {
	for _, a := range l.adapters {
		a.Close()
	}
}
