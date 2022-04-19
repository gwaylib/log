package logger

import (
	"os"
	"sync"
	"time"

	"github.com/gwaylib/log/proto"
)

var (
	DefaultContext = proto.Context{"default", "0.0.0", HostName}
)

// logger for private
type Logger struct {
	mu         sync.Mutex
	context    *proto.Context
	adapters   []proto.Adapter
	loggerName string
	level      proto.Level
}

func New(ctx *proto.Context, loggerName string, level proto.Level, adapter ...proto.Adapter) *Logger {
	lg := &Logger{
		context:    ctx,
		adapters:   adapter,
		loggerName: loggerName,
		level:      level,
	}
	return lg
}

func NewDefaultLogger(loggerName string, adapter ...proto.Adapter) *Logger {
	return New(&DefaultContext, loggerName, proto.LevelDebug, adapter...)
}

func (l *Logger) SetOutputLevel(level proto.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}
func (l *Logger) GetOutputLevel() proto.Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *Logger) SetAdapter(adapters []proto.Adapter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.adapters = adapters
}
func (l *Logger) GetAdapter() []proto.Adapter {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.adapters
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
func (l *Logger) ExitWithLevel(code int, level proto.Level, msg ...interface{}) {
	m := proto.ToMsg(msg...)
	l.put(level, m)
	l.Close()
	os.Exit(code)
}
func (l *Logger) Exit(code int, msg ...interface{}) {
	l.ExitWithLevel(code, proto.LevelInfo, msg...)
}
func (l *Logger) put(level proto.Level, msg []byte) {
	// TODO: performance of lock?
	l.mu.Lock()
	if level < l.level {
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

	l.Put([]*proto.Data{&proto.Data{time.Now(), level, l.loggerName, msg}})
}

func (l *Logger) Put(data []*proto.Data) {
	if len(data) == 0 {
		panic("len(data) == 0")
	}

	log := &proto.Proto{
		l.context,
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
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, a := range l.adapters {
		a.Close()
	}
}

func (l *Logger) Print(msg ...interface{}) {
	l.Info(msg...)
}
func (l *Logger) Printf(f string, msg ...interface{}) {
	l.Infof(f, msg...)
}
func (l *Logger) Panic(msg ...interface{}) {
	l.Fatal(msg...)
}
func (l *Logger) Panicf(f string, msg ...interface{}) {
	l.Fatalf(f, msg...)
}
