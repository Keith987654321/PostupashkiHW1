package mutex

import (
	"homework1/internal/futex"
	"runtime"
	"sync/atomic"
)

const (
	free = iota
	held
	contended
)

type Mutex struct {
	state uint32
}

func (m *Mutex) Lock() {
	for {
		if atomic.CompareAndSwapUint32(&m.state, free, held) {
			return
		}

		for i := 0; i < 200; i++ {
			if atomic.CompareAndSwapUint32(&m.state, free, held) {
				return
			}

			runtime.Gosched()
		}

		atomic.CompareAndSwapUint32(&m.state, held, contended)

		if atomic.LoadUint32(&m.state) == contended {
			futex.Wait(&m.state, contended)
		}
	}
}

func (m *Mutex) TryLock() bool {
	if atomic.LoadUint32(&m.state) != free {
		return false
	}

	return atomic.CompareAndSwapUint32(&m.state, free, held)
}

func (m *Mutex) Unlock() {
	currentState := atomic.LoadUint32(&m.state)

	switch currentState {
	case free:
		panic("already unlocked")
	case held:
		if atomic.CompareAndSwapUint32(&m.state, held, free) {
			return
		}
	case contended:
		if atomic.CompareAndSwapUint32(&m.state, contended, free) {
			futex.WakeAll(&m.state)
			return
		}
	}
}
