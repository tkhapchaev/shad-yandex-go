//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

var ErrStopped = errors.New("limiter stopped")

type Limiter struct {
	mutex    *Mutex
	max      int
	times    []time.Time
	interval time.Duration
	stop     chan struct{}
}

type Mutex struct {
	ch chan struct{}
}

func (l *Limiter) Acquire(ctx context.Context) error {
	now := time.Now()

	select {
	case <-l.stop:
		return ErrStopped
	default:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if len(l.times) < l.max {
		l.times = append(l.times, now)

		return nil
	}

	for l.times[0].Add(l.interval).After(now) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			now = time.Now()
		}
	}

	current := l.times[1:]
	l.times = append(current, now)

	return nil
}

func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	return &Limiter{
		mutex:    NewMutex(),
		max:      maxCount,
		interval: interval,
		stop:     make(chan struct{}),
		times:    make([]time.Time, 0, maxCount),
	}
}

func (l *Limiter) Stop() {
	close(l.stop)
}

func NewMutex() *Mutex {
	return &Mutex{ch: make(chan struct{}, 1)}
}

func (m *Mutex) Lock() {
	m.ch <- struct{}{}
}

func (m *Mutex) Unlock() {
	<-m.ch
}
