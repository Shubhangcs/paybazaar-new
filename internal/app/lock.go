package app

import "sync/atomic"

var apiLocked int32 = 0

func LockAPI() {
	atomic.StoreInt32(&apiLocked, 1)
}

func UnlockAPI() {
	atomic.StoreInt32(&apiLocked, 0)
}

func IsLocked() bool {
	return atomic.LoadInt32(&apiLocked) == 1
}
