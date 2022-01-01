package lserver

import (
	"testing"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/redis"
)

func TestPut(t *testing.T) {
	rs, err := redis.NewRediStore(1, "tcp", "127.0.0.1:6379", "")
	if err != nil {
		t.Fatal(err)
	}
	log := logger.NewDefaultLogger("testing", New(rs, "log.gway.cc"))
	log.Debug("debug")
	log.Debug([]byte{0xff})
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	//log.Fatal("fatal")
	//log.Exit(0, "exit")
	log.Close()
}
