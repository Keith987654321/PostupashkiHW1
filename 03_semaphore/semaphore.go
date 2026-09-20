package semaphore

type Semaphore struct {
	ch chan struct{}
}

func New(n int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, n),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}
