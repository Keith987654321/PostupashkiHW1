package rwmutex

import (
	"homework1/internal/futex"
	"sync/atomic"
)

const writer = uint32(1 << 31)
const readerMask = writer - 1

type RWMutex struct {
	state uint32
}

func (m *RWMutex) RLock() {
	for {
		state := atomic.LoadUint32(&m.state)

		if state&writer != 0 {
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

		if state&writer != 0 {
			panic("can't RUnlock while Locked")
		}

		readers := state & readerMask

		if readers == 0 {
			panic("can't RUnlock the unlocked mutex")
		}

		if atomic.CompareAndSwapUint32(&m.state, state, state-1) {
			if readers == 1 {
				futex.WakeAll(&m.state)
			}

			return
		}
	}
}

func (m *RWMutex) Lock() {
	for {
		state := atomic.LoadUint32(&m.state)

		if state == 0 {
			if atomic.CompareAndSwapUint32(&m.state, 0, writer) {
				return
			}

			continue
		}

		futex.Wait(&m.state, state)
	}
}

func (m *RWMutex) Unlock() {
	if !atomic.CompareAndSwapUint32(&m.state, writer, 0) {
		panic("Unlock of unlocked mutex")
	}

	futex.WakeAll(&m.state)
}
