package mutex

import (
	"homework1/internal/futex"
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
		}

		if atomic.CompareAndSwapUint32(&m.state, held, contended) {
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
		atomic.StoreUint32(&m.state, free)
	case contended:
		atomic.StoreUint32(&m.state, free)
		futex.Wake(&m.state)
	}
}
