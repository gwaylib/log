package console

import (
	"github.com/gwaylib/log/logger/adapter"
	"github.com/gwaylib/log/logger/proto"
	"github.com/labstack/gommon/color"
)

const (
	AdapterName = "console"
)

type consoleAdapter struct {
}

// put a log protocol to log queue
func (conA *consoleAdapter) Put(log *proto.LogProto) {
	if log == nil {
		panic("argument is nil")
	}
	for _, val := range log.Data {
		date := val.Date.Format("2006-01-02 15:04:05.000")
		color.Printf("%s %s [%s] %s\n",
			date,
			val.Level.ColorString(),
			color.Cyan(val.Logger),
			string(val.Msg),
		)
	}
}

func (conA *consoleAdapter) Close() {
	// nothing to close
}

func New() adapter.Adapter {
	return &consoleAdapter{}
}

func init() {
	adapter.Register(AdapterName, New())
}
