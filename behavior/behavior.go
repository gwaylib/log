package behavior

import (
	"encoding/json"
	"time"
)

type Event struct {
	IndexKey string // need build index key your self

	ReqHeader  string
	ReqParams  string
	RespStatus string
	RespParams string

	UseTime   time.Duration
	EventTime time.Time
}

func (b *Event) ToJson() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		panic(*b)
	}
	return data
}

func Parse(src []byte) (*Event, error) {
	b := &Event{}
	if err := json.Unmarshal(src, b); err != nil {
		return nil, err
	}
	return b, nil
}

type Client interface {
	Close()
	Put(*Event)
}
