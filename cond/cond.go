//go:build !solution

package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

type Barrier chan struct{}

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *sync.Mutex or *sync.RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
type Cond struct {
	locker   Locker
	barriers chan Barrier
}

func New(l Locker) Cond {
	c := Cond{
		locker:   l,
		barriers: make(chan Barrier, 1),
	}

	c.barriers <- make(Barrier)
	return c
}

func (c Cond) Broadcast() {
	b := <-c.barriers
	close(b)

	c.barriers <- make(Barrier)
}

func (c Cond) Signal() {
	b := <-c.barriers

	select {
	case b <- struct{}{}:
	default:
	}

	c.barriers <- b
}

func (c Cond) Wait() {
	b := <-c.barriers
	c.locker.Unlock()
	c.barriers <- b
	<-b
	c.locker.Lock()
}
