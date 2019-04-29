package lserver

import (
	"github.com/gwaylib/beanmsq"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/proto"
)

type Adapter struct {
	p beanmsq.Producer
}

// put a log protocol to log queue
func (a *Adapter) Put(p *proto.Proto) {
	if p == nil {
		panic("argument is nil")
	}
	data, err := proto.Marshal(p)
	if err != nil {
		logger.FailLog(errors.As(err, *p))
		return
	}
	// 16*1024*1024(16M)
	if len(data) > 16777216 {
		logger.FailLog(errors.New("data too big").As(err, *p))
		return
	}
	if err := a.p.Put(data); err != nil {
		logger.FailLog(errors.As(err, *p))
		return
	}
}

func (a *Adapter) Close() {
	if err := a.p.Close(); err != nil {
		logger.FailLog(err)
	}
}

func New(addr, tube string, pool int) logger.Adapter {
	return &Adapter{
		p: beanmsq.NewProducer(pool, addr, tube),
	}
}
