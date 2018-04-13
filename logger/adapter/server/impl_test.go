package stdio

import (
	"testing"

	"github.com/gwaylib/log/logger"
)

func TestPut(t *testing.T) {
	log := logger.New("testing", New("127.0.0.1:11300", "log.gway.cc", 100))
	log.Debug("debug")
	log.Debug([]byte{0xff})
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	//log.Fatal("fatal")
	log.Exit(0, "exit")
	log.Close()
}
