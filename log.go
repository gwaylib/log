package log

import (
	"os"
	"strconv"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/adapter/stdio"
	"github.com/gwaylib/log/proto"
)

const (
	// Setting the log level
	// 0 debug or above
	// 1 info or above
	// 2 warn or above
	// 3 error or above
	// 4 fatal
	GWAYLIB_LOG_LEVEL = "GWAYLIB_LOG_LEVEL"
)

var (
	Log = newWithCaller("", 4)
)

func newWithCaller(loggerName string, callerDepth int) proto.Logger {
	level, _ := strconv.Atoi(os.Getenv(GWAYLIB_LOG_LEVEL))
	ctx := &logger.DefaultContext
	adapter := []proto.Adapter{stdio.New(os.Stdout, os.Stderr)} // implement the adapter what you need
	return logger.New(ctx, loggerName, callerDepth, proto.Level(level), adapter...)
}

func New(loggerName string) proto.Logger {
	return newWithCaller(loggerName, 3)
}

// Debug level, for developer println the debug message.
//
// 值 0
// 调试信息,不提交服务器.
// 由平台控制,可控制台输出,相当于fmt.Print输出.
func Debug(msg ...interface{}) {
	Log.Debug(msg...)
}
func Debugf(f string, msg ...interface{}) {
	Log.Debugf(f, msg...)
}

// Same as Infof, golang log interface.
func Printf(format string, v ...interface{}) {
	Log.Printf(format, v...)
}

// Same as Info, golang log interface.
func Print(v ...interface{}) {
	Log.Print(v...)
}

// Same as Info, golang log interface.
func Println(v ...interface{}) {
	Info(v...)
}

// Information level, for status changed
// Example: boot, stop, reconnect, which changed the program enviroment.
// 值 1
// 程序运行状态信息,提交服务器.
// 如启动、停止、重连等信息，体现了程序环境的变更状态
func Info(msg ...interface{}) {
	Log.Info(msg...)
}
func Infof(f string, msg ...interface{}) {
	Log.Infof(f, msg...)
}

// Warning level, for status warning, maybe the status should be broken soon, but also maybe auto recover soon.
// Expect send a warning mail to maintaner.
//
// 值 2
// 程序异常信息，提交服务器. 本类别不影响程序继续使用,但其结果可能会引出潜在的重大问题.
// 例如：请求的数据格式错误；网络连接错误(但重新连接后可自动修复), 连接超时等行为。
// 此类异常在一段时间如果出现过多，那么应该分析其中的原因，例如可能存在的问题：
// 被攻击、硬件老化、硬件达到了承载上限、对方服务出现异常等问题。
// 日志系统将发送一封邮件到相关人员。
func Warn(msg ...interface{}) {
	Log.Warn(msg...)
}
func Warnf(f string, msg ...interface{}) {
	Log.Warnf(f, msg...)
}

// Error level, for status is error, need someone process the erorr, system refuse the service, but it still can be run.
// The Log system should call the mainter immediately.
//
// 值 3
// 程序致命的错误信息， 提交服务器。此错误将影响到正常逻辑, 甚至平台因此而恐慌、停止服务的行为.
// 例如：数据库不可用、充值不可用、短信不可用、vos不可用等需要及时处理的行为都可定义为此类别。
// 日志系统将发送一封邮件、短信(或者其他实时联系方式)至相关人员
func Error(msg ...interface{}) {
	Log.Error(msg...)
}
func Errorf(f string, msg ...interface{}) {
	Log.Errorf(f, msg...)
}

// Fatal level, the status is fatal, if continue run the program, it will cause irreparable things
// The Log system should call the mainter immediately.
//
// 值 4
// 检测到程序非正常结束。
// 日志系统应调用所有实时联系方式联系相关人员处理。
func Fatal(msg ...interface{}) {
	Log.Fatal(msg...)
}
func Fatalf(f string, msg ...interface{}) {
	Log.Fatalf(f, msg...)
}

// Same as Fatal
func Panic(msg ...interface{}) {
	Log.Fatal(msg...)
}
func Panicf(f string, msg ...interface{}) {
	Log.Fatalf(f, msg...)
}

// Exit the program with message log.
//
// 退出操作，执行关闭日志操作并调用os.Exit
// 对于可执行程序来说，日志的退出意味着程序的退出。
//
// Params
//
//	code -- code to call os.Exit
//	msg -- message log in information level.
func Exit(code int, msg ...interface{}) {
	Log.Exit(code, msg...)
}
