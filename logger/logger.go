package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gwaylib/log/proto"
)

var (
	DefaultContext = proto.Context{Platform: "default", Version: "0.0.0", Ip: HostName}
)

// logger for private
type Logger struct {
	mu          sync.Mutex
	context     *proto.Context
	adapters    []proto.Adapter
	loggerName  string
	level       proto.Level
	callerDepth int
}

func New(ctx *proto.Context, loggerName string, callerDepth int, level proto.Level, adapter ...proto.Adapter) *Logger {
	lg := &Logger{
		context:     ctx,
		adapters:    adapter,
		loggerName:  loggerName,
		level:       level,
		callerDepth: callerDepth,
	}
	return lg
}

func NewDefaultLogger(loggerName string, adapter ...proto.Adapter) *Logger {
	return New(&DefaultContext, loggerName, 4, proto.LevelDebug, adapter...)
}

// when depth<=0, close the caller path
func (l *Logger) SetCallerDepth(depth int) {
	l.callerDepth = depth
}

func (l *Logger) SetOutputLevel(level proto.Level) {
	l.level = level
}
func (l *Logger) GetOutputLevel() proto.Level {
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
	if l.levelClosed(proto.LevelDebug) {
		return
	}
	l.put(proto.LevelDebug, proto.ToMsg(msg...))
}
func (l *Logger) Debugf(f string, msg ...interface{}) {
	if l.levelClosed(proto.LevelDebug) {
		return
	}
	l.put(proto.LevelDebug, proto.ToMsgf(f, msg...))
}

func (l *Logger) Info(msg ...interface{}) {
	if l.levelClosed(proto.LevelInfo) {
		return
	}
	l.put(proto.LevelInfo, proto.ToMsg(msg...))
}
func (l *Logger) Infof(f string, msg ...interface{}) {
	if l.levelClosed(proto.LevelInfo) {
		return
	}
	l.put(proto.LevelInfo, proto.ToMsgf(f, msg...))
}
func (l *Logger) Print(msg ...interface{}) {
	if l.levelClosed(proto.LevelInfo) {
		return
	}
	l.put(proto.LevelInfo, proto.ToMsg(msg...))
}
func (l *Logger) Printf(f string, msg ...interface{}) {
	if l.levelClosed(proto.LevelInfo) {
		return
	}
	l.put(proto.LevelInfo, proto.ToMsgf(f, msg...))
}
func (l *Logger) Warn(msg ...interface{}) {
	if l.levelClosed(proto.LevelWarn) {
		return
	}
	l.put(proto.LevelWarn, proto.ToMsg(msg...))
}
func (l *Logger) Warnf(f string, msg ...interface{}) {
	if l.levelClosed(proto.LevelWarn) {
		return
	}
	l.put(proto.LevelWarn, proto.ToMsgf(f, msg...))
}
func (l *Logger) Error(msg ...interface{}) {
	if l.levelClosed(proto.LevelError) {
		return
	}
	l.put(proto.LevelError, proto.ToMsg(msg...))
}
func (l *Logger) Errorf(f string, msg ...interface{}) {
	if l.levelClosed(proto.LevelError) {
		return
	}
	l.put(proto.LevelError, proto.ToMsgf(f, msg...))
}
func (l *Logger) Fatal(msg ...interface{}) {
	if l.levelClosed(proto.LevelFatal) {
		return
	}
	m := proto.ToMsg(msg...)
	l.put(proto.LevelFatal, m)
	l.Close()
	panic(m)
}
func (l *Logger) Fatalf(f string, msg ...interface{}) {
	if l.levelClosed(proto.LevelFatal) {
		return
	}
	m := proto.ToMsgf(f, msg...)
	l.put(proto.LevelFatal, m)
	l.Close()
	panic(m)
}

func (l *Logger) Panic(msg ...interface{}) {
	if l.levelClosed(proto.LevelFatal) {
		return
	}
	m := proto.ToMsg(msg...)
	l.put(proto.LevelFatal, m)
	l.Close()
	panic(m)
}
func (l *Logger) Panicf(f string, msg ...interface{}) {
	if l.levelClosed(proto.LevelFatal) {
		return
	}
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
	if l.levelClosed(level) {
		return
	}

	m := proto.ToMsg(msg...)
	l.put(level, m)
	l.Close()
	os.Exit(code)
}
func (l *Logger) Exit(code int, msg ...interface{}) {
	if l.levelClosed(proto.LevelInfo) {
		return
	}

	m := proto.ToMsg(msg...)
	l.put(proto.LevelInfo, m)
	l.Close()
	os.Exit(code)
}

func caller(depth int) string {
	at := ""
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		at = "unknown"
	}
	fileFields := strings.Split(file, "/")
	if len(fileFields) < 1 {
		at = "no-fields"
		return at
	}

	fileName := strings.Join(fileFields[len(fileFields)-1:], "/")
	return fmt.Sprintf("%s:%d", fileName, line)
}
func (l *Logger) levelClosed(level proto.Level) bool {
	// TODO: performance of lock?
	return level > l.level
}
func (l *Logger) put(level proto.Level, msg []byte) {
	if l.levelClosed(level) {
		return
	}

	loggerName := l.loggerName
	if l.callerDepth > 0 {
		if len(loggerName) > 0 {
			loggerName = loggerName + ":" + caller(l.callerDepth)
		} else {
			loggerName = caller(l.callerDepth)
		}
	}
	l.Put([]*proto.Data{&proto.Data{
		Date:   time.Now(),
		Level:  level,
		Logger: loggerName,
		Msg:    msg,
	}})
}

func (l *Logger) Put(data []*proto.Data) {
	if len(data) == 0 {
		panic("len(data) == 0")
	}

	log := &proto.Proto{
		Context: l.context,
		Data:    data,
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
