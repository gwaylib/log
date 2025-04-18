// Deprecated, use redismsq instead
package lserver

import (
	"github.com/gwaylib/beanmsq"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/behavior"
	"github.com/gwaylib/log/proto"
)

type lserverAdapter struct {
	p beanmsq.Producer
}

// put a log protocol to log queue
func (a *lserverAdapter) Put(e *behavior.Event) {
	data := e.ToJson()
	// 16*1024*1024(16M)
	if len(data) > 16777216 {
		proto.FailLog(errors.New("data too big").As(*e))
		return
	}
	if err := a.p.Put(data); err != nil {
		proto.FailLog(errors.As(err, *e))
		return
	}
}

func (a *lserverAdapter) Close() {
	if err := a.p.Close(); err != nil {
		proto.FailLog(err)
	}
}

func NewLServerClient(addr, tube string, pool int) behavior.Client {
	return &lserverAdapter{
		p: beanmsq.NewProducer(pool, addr, tube),
	}
}
