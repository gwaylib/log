// Deprecated, use rmsq-server instead
package behavior

import (
	"github.com/gwaylib/beanmsq"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/logger"
)

type lserverAdapter struct {
	p beanmsq.Producer
}

// put a log protocol to log queue
func (a *lserverAdapter) Put(e *Event) {
	data := e.ToJson()
	// 16*1024*1024(16M)
	if len(data) > 16777216 {
		logger.FailLog(errors.New("data too big").As(*e))
		return
	}
	if err := a.p.Put(data); err != nil {
		logger.FailLog(errors.As(err, *e))
		return
	}
}

func (a *lserverAdapter) Close() {
	if err := a.p.Close(); err != nil {
		logger.FailLog(err)
	}
}

func NewLServerClient(addr, tube string, pool int) Client {
	return &lserverAdapter{
		p: beanmsq.NewProducer(pool, addr, tube),
	}
}
