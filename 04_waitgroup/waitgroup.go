package waitgroup

import (
	"homework1/internal/futex"
	"sync/atomic"
)

type WaitGroup struct {
	count uint32
}

func (wg *WaitGroup) Add(delta int) {
	if delta == 0 {
		return
	}

	if delta > 0 {
		atomic.AddUint32(&wg.count, uint32(delta))
		return
	}

	decrement := uint32(-delta)

	for {
		current := atomic.LoadUint32(&wg.count)

		after := int(current) - int(decrement)
		if after < 0 {
			panic("counter can't be negative")
		}

		if atomic.CompareAndSwapUint32(&wg.count, current, uint32(after)) {
			if after == 0 {
				futex.WakeAll(&wg.count)
			}

			return
		}
	}

}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	for {
		current := atomic.LoadUint32(&wg.count)

		if current == 0 {
			return
		}

		futex.Wait(&wg.count, current)
	}
}
