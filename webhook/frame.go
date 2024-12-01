package webhook

import (
	"runtime/debug"
	"strconv"
	"sync/atomic"

	"github.com/2mf8/Go-QQ-Client/openapi"
	log "github.com/sirupsen/logrus"
)

type Frame struct {
	BotId   uint64
	Echo    string
	Ok      bool
	Openapi openapi.OpenAPI
}

var GlobalId int64 = 0

func SafeGo(fn func()) {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				log.Errorf("err recovered: %+v", e)
				log.Errorf("%s", debug.Stack())
			}
		}()
		fn()
	}()
}

func GenerateId() int64 {
	return atomic.AddInt64(&GlobalId, 1)
}

func GenerateIdStr() string {
	return strconv.FormatInt(GenerateId(), 10)
}
