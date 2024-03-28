//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type bc = chan bool

type KeyLock struct {
	mutex *sync.Mutex
	m     map[string]bc
}

func New() *KeyLock {
	return &KeyLock{mutex: &sync.Mutex{}, m: make(map[string]bc)}
}

func Prepare(keys []string) ([]string, bool) {
	key := make([]string, len(keys))
	copy(key, keys)
	sort.Strings(key)
	r := false

	return key, r
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	key, r := Prepare(keys)

	unlocked := func() {
		l.mutex.Lock()

		for _, k := range key {
			select {
			case l.m[k] <- true:
			default:
			}
		}

		l.mutex.Unlock()
	}

	for _, k := range key {
		l.mutex.Lock()
		value, ok := l.m[k]

		if !ok {
			c := make(bc, 1)
			c <- true
			l.m[k] = c
			value = c
		}

		l.mutex.Unlock()
	L:
		for {
			select {
			case <-cancel:
				{
					unlocked()
					r = true

					break L
				}
			case <-value:
				{
					break L
				}
			}
		}
	}

	return r, unlocked
}
