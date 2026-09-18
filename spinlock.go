package spinlock

import "sync/atomic"

type Spinlock struct {
	locked atomic.Bool
}

func (s *Spinlock) Lock() {
	for !s.locked.CompareAndSwap(false, true) {
	}
}

func (s *Spinlock) TryLock() bool {
	return s.locked.CompareAndSwap(false, true)
}

func (s *Spinlock) Unlock() {
	if !s.locked.Load() {
		panic("already unlocked")
	}
	s.locked.Store(false)
}

type TTAS struct {
	locked atomic.Bool
}

func (s *TTAS) Lock() {
	for {
		for s.locked.Load() {
		}

		if s.locked.CompareAndSwap(false, true) {
			break
		}
	}
}

func (s *TTAS) TryLock() bool {
	if s.locked.Load() {
		return false
	}
	return s.locked.CompareAndSwap(false, true)
}

func (s *TTAS) Unlock() {
	if !s.locked.Load() {
		panic("already unlocked")
	}

	s.locked.Store(false)
}
