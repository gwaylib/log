package lserver

import (
	"crypto/md5"
	"fmt"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/logger"
	"github.com/gwaylib/log/logger/proto"
	"github.com/gwaylib/redis"
	rmsq "github.com/gwaylib/redis/msq"
)

type Adapter struct {
	rs *redis.RediStore
	p  *rmsq.RedisMsqProducer
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
	if err := a.p.Put(fmt.Sprintf("%x", md5.Sum(data)), data); err != nil {
		logger.FailLog(errors.As(err, *p))
		return
	}
}

func (a *Adapter) Close() {
	if err := a.rs.Close(); err != nil {
		logger.FailLog(err)
	}
}

func New(rs *redis.RediStore, streamName string) logger.Adapter {
	return &Adapter{
		rs: rs,
		p:  rmsq.NewRedisMsqProducer(rs, streamName),
	}
}
