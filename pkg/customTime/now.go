package customTime

import (
	"sync/atomic"
	"time"
)

var (
	Now        = atomic.Int64{}
	WasStopped = atomic.Bool{}
)

func Start() {
	Now.Store(time.Now().Unix())
	WasStopped.Store(false)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if WasStopped.Load() {
				return
			}
			Now.Store(time.Now().Unix())
		}
	}()
}
