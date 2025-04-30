package stdio

import (
	"os"
	"testing"

	"github.com/gwaylib/log/logger"
)

func TestPut(t *testing.T) {
	log := logger.NewDefaultLogger("testing", New(os.Stdout, os.Stderr))
	log.Debug("debug")
	log.Debug([]byte{0xff})
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	// log.Fatal("fatal")
	//log.Exit(0, "exit")
	log.Close()
}
