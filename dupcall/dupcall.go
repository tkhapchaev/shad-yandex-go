//go:build !solution

package dupcall

import (
	"context"
	"errors"
	"sync"
)

type Call struct {
	mutex sync.Mutex
	ctx   context.Context

	count    int
	executed bool
	err      error

	store func()
	obj   interface{}
}

func (o *Call) Prepare() func() {
	o.count = 0
	o.executed = true
	s, c := context.WithCancel(context.Background())

	o.store = c
	o.ctx = s
	o.mutex.Unlock()

	return c
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.mutex.Lock()

	if o.executed {
		o.count++
		o.mutex.Unlock()

		select {
		case <-o.ctx.Done():
			defer func() {
				o.mutex.Unlock()
			}()

			o.mutex.Lock()
			o.count--

			return o.obj, o.err
		case <-ctx.Done():
			defer func() {
				o.mutex.Unlock()
			}()

			o.mutex.Lock()
			o.count--

			if o.count <= 0 {
				o.store()

				return o.obj, errors.New("task failed")
			}

			return o.obj, errors.New("task failed")
		}
	}

	cancel := o.Prepare()

	go func() {
		defer func() {
			o.executed = false
			o.count = 0
			o.mutex.Unlock()

			cancel()
		}()

		value, e := cb(o.ctx)
		o.mutex.Lock()
		o.err = e
		o.obj = value
	}()

	return o.Do(ctx, cb)
}
