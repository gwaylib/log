package logger

import (
	"fmt"
	"time"

	"github.com/gwaylib/log/logger/proto"
)

type Adapter interface {
	Put(log *proto.Proto)
	Close()
}

// 处理适配器记录日志时产的错误，以println的方式输出至控制台。
// TODO: 将此错误写到系统日志当中。
func FailLog(log *proto.Proto, err error) {
	// log error
	failMsg := fmt.Sprintf("err:%s,len:%d", err.Error(), len(log.Data))
	data := &proto.Data{
		time.Now(),
		proto.LevelFatal,
		"faillog",
		proto.ToMsg(failMsg),
	}
	println(data.String())

	// log to console with error of console code.
	for _, data := range log.Data {
		println(data.String())
	}
}
