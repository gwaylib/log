package behavior

import (
	"crypto/md5"
	"fmt"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log/proto"
	"github.com/gwaylib/redis"
	rmsq "github.com/gwaylib/redis/msq"
)

type rmsqAdapter struct {
	p rmsq.MsqProducer
}

// put a log protocol to log queue
func (a *rmsqAdapter) Put(p *Event) {
	if p == nil {
		panic("argument is nil")
	}
	data := p.ToJson()
	// 16*1024*1024(16M)
	if len(data) > 16777216 {
		proto.FailLog(errors.New("data too big").As(*p))
		return
	}
	if err := a.p.Put(fmt.Sprintf("%x", md5.Sum(data)), data); err != nil {
		proto.FailLog(errors.As(err, *p))
		return
	}
}

func (a *rmsqAdapter) Close() {
	if err := a.p.Close(); err != nil {
		proto.FailLog(err)
	}
}

func NewRMSQClient(rs *redis.RediStore, streamName string) Client {
	return &rmsqAdapter{
		p: rmsq.NewMsqProducer(rs, streamName),
	}
}
