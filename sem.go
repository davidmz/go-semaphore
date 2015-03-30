package semaphore

import "sync/atomic"

type Semaphore struct {
	ch       chan struct{}
	waiting  int64
	acquired int64
}

// Create a new Semaphore, setting max concurrent jobs to count
func New(count int) *Semaphore {
	return &Semaphore{ch: make(chan struct{}, count)}
}

// Create a new Semaphore, setting max concurrent jobs to 1
func Mutex() *Semaphore { return New(1) }

// Aquire a job, blocks until one slot is free
func (s *Semaphore) Acquire() *Semaphore {
	atomic.AddInt64(&s.waiting, 1)
	defer atomic.AddInt64(&s.waiting, -1)
	defer atomic.AddInt64(&s.acquired, 1)

	s.ch <- struct{}{}
	return s
}

// Release one job
func (s *Semaphore) Release() {
	defer atomic.AddInt64(&s.acquired, -1)
	<-s.ch
}

// Try to lock, return true if locked
func (s *Semaphore) Try() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
	}
	return false
}

func (s *Semaphore) WaitingCount() int {
	return int(atomic.LoadInt64(&s.waiting))
}

func (s *Semaphore) AcquiredCount() int {
	return int(atomic.LoadInt64(&s.acquired))
}
