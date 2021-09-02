package pubsub

import (
	"sync"
)

type Subscriber struct {
	Chan chan interface{}
	Id   string
}

func NewSubscriber(id string) *Subscriber {
	ch := make(chan interface{}, 1)
	return &Subscriber{Chan: ch, Id: id}
}

type Pubsub struct {
	mu     sync.RWMutex
	closed bool
	subs   map[string][]*Subscriber
}

func NewPubsub() *Pubsub {
	ps := &Pubsub{}
	ps.subs = make(map[string][]*Subscriber)
	return ps
}

//func (ps *Pubsub) Subscribe(id string, topic string) <-chan interface{} {
func (ps *Pubsub) Subscribe(id string, topic string) *Subscriber {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub := NewSubscriber(id)
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
		go func(sub *Subscriber) {
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
