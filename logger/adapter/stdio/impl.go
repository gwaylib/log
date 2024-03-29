package stdio

import (
	"fmt"
	"io"

	"github.com/gwaylib/log/proto"
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
		if a.stderr != nil {
			a.stderr.Print(output)
		} else if a.stdout != nil {
			a.stdout.Print(output)
		}
	}
}

func (a *Adapter) Close() {
}

func New(stdout, stderr io.Writer) proto.Adapter {
	stdoutC := new(color.Color)
	stdoutC.SetOutput(stdout)
	stderrC := new(color.Color)
	stderrC.SetOutput(stderr)
	return &Adapter{
		stdout: stdoutC,
		stderr: stderrC,
	}
}
