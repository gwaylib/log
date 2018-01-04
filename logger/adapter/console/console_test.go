package console

import (
	"testing"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/adapter"
	"github.com/gwaylib/log/logger/proto"
)

func TestPut(t *testing.T) {
	adapter.Register(AdapterName, New())
	log := logger.NewLogger([]string{AdapterName}, "console_test", proto.LevelDebug)
	log.Debug("debug")
	log.Debug([]byte{0xff})
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	//log.Fatal("fatal")
	log.Exit(0, "exit")
	log.Close()
}
