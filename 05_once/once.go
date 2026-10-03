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

func (o *Once) DoSlow(f func()) {
	current := atomic.LoadUint32(&o.state)
	if current == notStarted && atomic.CompareAndSwapUint32(&o.state, current, running) {
		defer func() {
			atomic.StoreUint32(&o.state, done)
			futex.WakeAll(&o.state)
		}()
		f()
		atomic.StoreUint32(&o.state, done)
		futex.WakeAll(&o.state)
		return
	}

	if current == running {
		futex.Wait(&o.state, running)
	}
}

func (o *Once) Do(f func()) {
	for {
		if atomic.LoadUint32(&o.state) == done {
			return
		}

		o.DoSlow(f)

	}
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.state) == done
}
