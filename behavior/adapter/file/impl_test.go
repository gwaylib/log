package file

import (
	"strconv"
	"testing"
	"time"

	"github.com/gwaylib/log/behavior"
)

func TestPut(t *testing.T) {
	a := NewFile("./behavior.log", 10, 1024)
	a.Put(&behavior.Event{
		IndexKey:   strconv.FormatInt(time.Now().UnixNano(), 10),
		ReqHeader:  "testing",
		ReqParams:  "p=testing",
		RespStatus: "200",
		RespParams: "p=done",
		UseTime:    1,
		EventTime:  time.Now(),
	})
	a.Close()
}
