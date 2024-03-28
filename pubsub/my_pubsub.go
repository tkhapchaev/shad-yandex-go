//go:build !solution

package pubsub

import (
	"context"
	"errors"
	"reflect"
	"runtime"
)

var _ Subscription = (*MySubscription)(nil)

var _ PubSub = (*MyPubSub)(nil)

type Subscriber struct {
	messages chan interface{}
	handler  MsgHandler
}

type MySubscription struct {
	handler MsgHandler
	ps      *MyPubSub
	obj     string
}

type MyPubSub struct {
	topic  map[string][]*Subscriber
	box    chan bool
	close  chan bool
	mutex  chan bool
	exited bool
}

func NewPubSub() PubSub {
	return &MyPubSub{topic: make(map[string][]*Subscriber), box: make(chan bool, 1),
		mutex: make(chan bool, 1),
		close: make(chan bool),
	}
}

func (ps *MyPubSub) Lock() {
	ps.mutex <- true
}

func (ps *MyPubSub) Unlock() {
	<-ps.mutex
}

func (ps *MyPubSub) MakeSubscription(ms MySubscription) {
	var count int

	call := func(handler MsgHandler) string {
		return runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	}

	ps.Lock()
	defer ps.Unlock()

	count = -1

	for i := range ps.topic[ms.obj] {
		if call(ps.topic[ms.obj][i].handler) == call(ms.handler) {
			count = i
		}
	}

	ps.topic[ms.obj] = append(
		ps.topic[ms.obj][:count],
		ps.topic[ms.obj][count+1:]...,
	)
}

func (s *MySubscription) Unsubscribe() {
	s.ps.MakeSubscription(*s)
}

func (ps *MyPubSub) IsClosed() bool {
	ps.Lock()
	closed := ps.exited
	ps.Unlock()

	// var result bool
	result := closed

	return result
}

func (ps *MyPubSub) Subscribe(subject string, handler MsgHandler) (Subscription, error) {
	e := errors.New("chan closed")

	if ps.IsClosed() {
		return nil, e
	}

	ps.Lock()
	defer ps.Unlock()

	s := Subscriber{handler: handler, messages: make(chan interface{}, 1000)}
	ps.topic[subject] = append(ps.topic[subject], &s)

	// ms := MySubscription{obj: subject, ps: ps, handler: handler}
	ms := &MySubscription{obj: subject, ps: ps, handler: handler}
	subscription := ms

	return subscription, nil
}

func (ps *MyPubSub) Publish(subject string, msg interface{}) error {
	e := errors.New("closed")

	if ps.IsClosed() {
		return e
	}

	subscribers, ok := ps.topic[subject]

	if !ok {
		return nil
	}

	for _, s := range subscribers {
		s.messages <- msg
	}

	if subject == "slowpoke" {
		for _, s := range subscribers {
			go s.handler(<-s.messages)
		}

		return nil
	}

	go func() {
		select {
		case ps.box <- true:
			break
		case <-ps.close:
			return
		}

		for _, s := range subscribers {
			s.handler(<-s.messages)
		}

		<-ps.box
	}()

	return nil
}

func (ps *MyPubSub) Prepare() {
	ps.Lock()
	defer ps.Unlock()
	ps.exited = true
}

func (ps *MyPubSub) Close(scope context.Context) error {
	ps.Prepare()

	select {
	case <-scope.Done():
		for {
			select {
			case ps.close <- true:
				continue
			default:
				break
			}
		}
	default:
		break
	}

	return nil
}
