package stdio

import (
	"fmt"
	"io"

	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/proto"
	"github.com/labstack/gommon/color"
)

type Adapter struct {
	stdout *color.Color
	stderr *color.Color
}

// put a log protocol to log queue
func (a *Adapter) Put(p *proto.Proto) {
	if p == nil {
		panic("argument is nil")
	}
	for _, val := range p.Data {
		date := val.Date.Format("2006-01-02T15:04:05.000Z07:00")
		output := fmt.Sprintf("%s %-5s [%s] %s",
			date,
			val.Level.ColorString(),
			color.Cyan(val.Logger),
			string(val.Msg),
		)
		if val.Level == proto.LevelDebug {
			a.stdout.Print(output)
		} else {
			a.stderr.Print(output)
		}
	}
}

func (a *Adapter) Close() {
}

func New(stdout, stderr io.Writer) logger.Adapter {
	stdoutC := new(color.Color)
	stdoutC.SetOutput(stdout)
	stderrC := new(color.Color)
	stderrC.SetOutput(stderr)
	return &Adapter{
		stdout: stdoutC,
		stderr: stderrC,
	}
}
