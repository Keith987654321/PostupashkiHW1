package rwmutex

import (
	"homework1/internal/futex"
	"sync/atomic"
)

const (
	writer        = uint32(1 << 31)
	writerWaiting = uint32(1 << 30)
	readerMask    = writerWaiting - 1
)

type RWMutex struct {
	state          uint32
	waitingWriters uint32
}

func (m *RWMutex) RLock() {
	for {
		state := atomic.LoadUint32(&m.state)

		if state&(writer|writerWaiting) != 0 {
			futex.Wait(&m.state, state)
			continue
		}

		if atomic.CompareAndSwapUint32(&m.state, state, state+1) {
			return
		}
	}
}

func (m *RWMutex) RUnlock() {
	for {
		state := atomic.LoadUint32(&m.state)

		readers := state & readerMask

		if readers == 0 {
			panic("can't RUnlock the unlocked mutex")
		}

		newState := state - 1

		if atomic.CompareAndSwapUint32(&m.state, state, newState) {
			if readers == 1 {
				futex.WakeAll(&m.state)
			}

			return
		}
	}
}

func (m *RWMutex) Lock() {
	atomic.AddUint32(&m.waitingWriters, 1)

	for {
		state := atomic.LoadUint32(&m.state)

		if state&writerWaiting != 0 {
			break
		}

		if atomic.CompareAndSwapUint32(
			&m.state,
			state,
			state|writerWaiting,
		) {
			break
		}
	}

	for {
		state := atomic.LoadUint32(&m.state)

		if state&(readerMask|writer) == 0 {
			newState := state | writer

			if atomic.CompareAndSwapUint32(&m.state, state, newState) {
				atomic.AddUint32(&m.waitingWriters, ^uint32(0))

				for {
					waiting := atomic.LoadUint32(&m.waitingWriters)

					if waiting != 0 {
						break
					}

					current := atomic.LoadUint32(&m.state)

					if current&writerWaiting == 0 {
						break
					}

					if atomic.CompareAndSwapUint32(
						&m.state,
						current,
						current&^writerWaiting,
					) {
						break
					}
				}

				return
			}

			continue
		}

		futex.Wait(&m.state, state)
	}
}

func (m *RWMutex) Unlock() {
	for {
		state := atomic.LoadUint32(&m.state)

		if state&writer == 0 {
			panic("Unlock of unlocked mutex")
		}

		newState := state &^ writer

		if atomic.CompareAndSwapUint32(&m.state, state, newState) {
			futex.WakeAll(&m.state)
			return
		}
	}
}
