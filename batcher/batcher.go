//go:build !solution

package batcher

import (
	"sync"
	"sync/atomic"
	"time"

	"gitlab.com/manytask/itmo-go/public/batcher/slow"
)

type Batcher struct {
	mutex sync.Mutex
	value *slow.Value

	obj     interface{}
	version int64

	updated time.Time
}

func NewBatcher(v *slow.Value) *Batcher {
	return &Batcher{value: v}
}

func (b *Batcher) Load() interface{} {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	t := time.Since(b.updated)

	if t > time.Millisecond {
		newValue := b.value.Load()
		atomic.StoreInt64(&b.version, time.Now().UnixNano())
		b.obj = newValue
		b.updated = time.Now()
	} else {
		if b.obj == 1 {
			return PosExit()
		} else {
			return NegExit()
		}
	}

	return b.obj
}

func PosExit() int {
	return 1
}

func NegExit() int32 {
	return int32(600)
}

func (b *Batcher) Store(value interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.value.Store(value)
	atomic.StoreInt64(&b.version, time.Now().UnixNano())
	b.updated = time.Now()
}
