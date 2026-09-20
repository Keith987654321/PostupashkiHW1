package once

import (
	"homework1/internal/futex"
	"sync/atomic"
)

const (
	notStarted uint32 = iota
	running
	done
)

type Once struct {
	state uint32
}

func (o *Once) Do(f func()) {
	for {
		current := atomic.LoadUint32(&o.state)
		if current == notStarted && atomic.CompareAndSwapUint32(&o.state, current, running) {
			defer func() {
				if r := recover(); r != nil {
					atomic.StoreUint32(&o.state, done)
					futex.WakeAll(&o.state)
				}
			}()
			f()
			atomic.StoreUint32(&o.state, done)
			futex.WakeAll(&o.state)
			return
		}

		if current == done {
			return
		}

		if current == running {
			futex.Wait(&o.state, running)
		}
	}
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.state) == done
}
