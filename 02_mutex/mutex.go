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
	if atomic.CompareAndSwapUint32(&m.state, free, held) {
		return
	}

	for i := 0; i < 200; i++ {
		if atomic.CompareAndSwapUint32(&m.state, free, held) {
			return
		}

		runtime.Gosched()
	}

	for {
		old := atomic.SwapUint32(&m.state, contended)

		if old == free {
			return
		} else {
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
	old := atomic.SwapUint32(&m.state, free)

	switch old {
	case free:
		panic("already unlocked")

	case held:
		return

	case contended:
		futex.Wake(&m.state)
	}
}
