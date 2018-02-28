package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/gwaylib/log/logger/adapter"
	"github.com/gwaylib/log/logger/proto"
)

var (
	DefaultContext = proto.Context{ServerName, "1.0.0", HostIp}
)

// logger for private
type Logger struct {
	context    *proto.Context
	adapters   map[string]adapter.Adapter
	loggerName string
	IsDebug    bool
	IsInfo     bool
	IsWarn     bool
	IsErr      bool
	IsFatal    bool
}

func New(loggerName string, level proto.Level, adapterName ...string) proto.Log {
	adapters := map[string]adapter.Adapter{}
	for _, name := range adapterName {
		adapter, err := adapter.GetAdapter(name)
		if err != nil {
			panic(err)
		}
		adapters[name] = adapter
	}

	lg := &Logger{
		adapters:   adapters,
		loggerName: loggerName,
	}
	lg.SetLevel(int(level))
	return lg
}

func (l *Logger) SetLevel(level int) {
	l.IsDebug = level <= 0
	l.IsInfo = level <= 1
	l.IsWarn = level <= 2
	l.IsErr = level <= 3
	l.IsFatal = level <= 4
}

// 设置服务器信息
// 不设置时，使用默认配置信息
func (l *Logger) SetContext(c *proto.Context) {
	l.context = c
}

// Debug
// debug level is use to log anything, so it will make a lot of log.
//
// Param
// msg -- inteface of json object
func (l *Logger) Debug(msg ...interface{}) {
	if !l.IsDebug {
		return
	}

	l.put(&proto.Data{
		time.Now(),
		proto.LevelDebug,
		l.loggerName,
		proto.ToMsg(msg...),
	})
}

func (l *Logger) Debugf(f string, msg ...interface{}) {
	l.Debug(fmt.Sprintf(f, msg...))
}

// Info
// info level is used to log something is changed with program envirement.
//
// Param
// msg -- inteface of json object
func (l *Logger) Info(msg ...interface{}) {
	if !l.IsInfo {
		return
	}
	l.put(&proto.Data{
		time.Now(),
		proto.LevelInfo,
		l.loggerName,
		proto.ToMsg(msg...),
	})
}

func (l *Logger) Infof(f string, msg ...interface{}) {
	l.Info(fmt.Sprintf(f, msg...))
}

// Warn
// warn level is used to log program has some error in controling,
// but if this warning happen again, program will go to error level, or panic.
//
// Param
// msg -- inteface of json object
func (l *Logger) Warn(msg ...interface{}) {
	if !l.IsWarn {
		return
	}

	l.put(&proto.Data{
		time.Now(),
		proto.LevelWarn,
		l.loggerName,
		proto.ToMsg(msg...),
	})
}

func (l *Logger) Warnf(f string, msg ...interface{}) {
	l.Warn(fmt.Sprintf(f, msg...))
}

// Error
// error level is used to log program need someone do something immediately.
//
// Param
// msg -- inteface of json object
func (l *Logger) Error(msg ...interface{}) {
	if !l.IsErr {
		return
	}

	l.put(&proto.Data{
		time.Now(),
		proto.LevelError,
		l.loggerName,
		proto.ToMsg(msg...),
	})
}

func (l *Logger) Errorf(f string, msg ...interface{}) {
	l.Error(fmt.Sprintf(f, msg...))
}

// Fatal
// fatal level is used to log program panic.
//
// Param
// msg -- inteface of json object
func (l *Logger) Fatal(msg ...interface{}) {
	if !l.IsFatal {
		return
	}
	data := proto.Data{
		time.Now(),
		proto.LevelFatal,
		l.loggerName,
		proto.ToMsg(msg...),
	}
	l.put(&data)
	l.Close()
	panic(data)
}

func (l *Logger) Fatalf(f string, msg ...interface{}) {
	l.Fatal(fmt.Sprintf(f, msg...))
}

// Exit
// log an info level message, and close log, then call os.Exit(code)
//
// Param
// code -- code of os exit
// msg -- exit message
func (l *Logger) Exit(code int, msg ...interface{}) {
	if !l.IsFatal {
		return
	}
	data := proto.Data{
		time.Now(),
		proto.LevelInfo,
		l.loggerName,
		proto.ToMsg(msg...),
	}
	l.put(&data)
	l.Close()
	os.Exit(code)
}

func (l *Logger) put(data *proto.Data) {
	if data == nil {
		return
	}

	l.Put([]*proto.Data{data})
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

	log := &proto.LogProto{
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
