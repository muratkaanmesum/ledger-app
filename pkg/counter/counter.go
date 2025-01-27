package counter

import (
	"sync/atomic"
)

type Stats struct {
	totalRequests int64
	success       int64
	fail          int64
}

var instance Stats

func InitStats() {
	instance = Stats{}
}

func AddSuccess() {
	atomic.AddInt64(&instance.success, 1)
}

func AddFail() {
	atomic.AddInt64(&instance.fail, 1)
}

func AddTotalRequests() {
	atomic.AddInt64(&instance.totalRequests, 1)
}
