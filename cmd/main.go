package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	signalch chan os.Signal
)

type Pubsub struct {
	mu     sync.RWMutex
	closed bool
	subs   map[string][]chan string
}

func NewPubsub() *Pubsub {
	ps := &Pubsub{}
	ps.subs = make(map[string][]chan string)
	return ps
}

func (ps *Pubsub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *Pubsub) Publish(topic string, msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if ps.closed {
		return
	}
	for _, ch := range ps.subs[topic] {
		go func(ch chan string) {
			ch <- msg
		}(ch)

	}
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}

func qualityagent(input interface{}, agenttype int) interface{} {
	time.Sleep(2 * time.Second)
	return input
}

func worker(id int) {
	content := fmt.Sprintf("obj: %d", id)
	r := qualityagent(content, 1)
	fmt.Printf("id: %d\n", id)
	fmt.Printf("%s\n", r)
}

func main() {

	signalch = make(chan os.Signal, 1)
	NodeNumber := 5

	ps := NewPubsub()

	listener := func(name string, ch <-chan string) {
		for i := range ch {
			fmt.Printf("[%s] got %s\n", name, i)
		}
		fmt.Printf("[%s] done\n", name)
	}

	for i := 1; i < NodeNumber+1; i++ {
		//topicname := fmt.Sprintf("node-%d", i)
		topicname := fmt.Sprintf("node-%d", 0)
		ch := ps.Subscribe(topicname)
		go listener(topicname, ch)
	}
	for i := 1; i < NodeNumber+1; i++ {
		//topicname := fmt.Sprintf("node-%d", i)
		topicname := fmt.Sprintf("node-%d", 0)
		msg := fmt.Sprintf("msg-%d", i)
		ps.Publish(topicname, msg)

	}

	signal.Notify(signalch, os.Interrupt, os.Kill, syscall.SIGTERM)
	signalType := <-signalch
	signal.Stop(signalch)
	fmt.Printf("On Signal <%s>\n", signalType)
}
