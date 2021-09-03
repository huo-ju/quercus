package pubsub

import (
	"context"
	"errors"
	"sync"
)

type Subscription struct {
	Chan  chan interface{}
	topic string
	ctx   context.Context
	Id    string
}

func NewSubscription(id string, topic string) *Subscription {
	ch := make(chan interface{}, 1)
	return &Subscription{Chan: ch, topic: topic, Id: id}
}

func (sub *Subscription) Topic() string {
	return sub.topic
}

func (sub *Subscription) Next(ctx context.Context) (interface{}, error) {
	select {
	case msg, ok := <-sub.Chan:
		if !ok {
			return msg, errors.New("some err")
		}
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type Pubsub struct {
	mu     sync.RWMutex
	closed bool
	subs   map[string][]*Subscription
}

func NewPubsub() *Pubsub {
	ps := &Pubsub{}
	ps.subs = make(map[string][]*Subscription)
	return ps
}

func (ps *Pubsub) Subscribe(id string, topic string) *Subscription {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub := NewSubscription(id, topic)
	ps.subs[topic] = append(ps.subs[topic], sub)
	//return sub.Chan
	return sub
}

func (ps *Pubsub) Publish(topic string, msg interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if ps.closed {
		return
	}
	for _, sub := range ps.subs[topic] {
		go func(sub *Subscription) {
			sub.Chan <- msg
		}(sub)

	}
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, sub := range subs {
				close(sub.Chan)
			}
		}
	}
}
