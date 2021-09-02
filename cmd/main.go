package main

import (
	"fmt"
	"github.com/huo-ju/quercus/pkg/pubsub"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	signalch chan os.Signal
)

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

	ps := pubsub.NewPubsub()

	//ch <-chan interface{}
	listener := func(topicname string, sub *pubsub.Subscriber) {
		for i := range sub.Chan {
			fmt.Printf("Node:(%s) [%s] got %s\n", sub.Id, topicname, i)
		}
		fmt.Printf("Node:(%s) [%s] done\n", sub.Id, topicname)
	}

	for i := 1; i < NodeNumber+1; i++ {
		nodename := fmt.Sprintf("node-%d", i)
		topicname := fmt.Sprintf("topic-%d", 0)
		sub := ps.Subscribe(nodename, topicname)
		go listener(topicname, sub)
	}
	for i := 1; i < NodeNumber+1; i++ {
		//topicname := fmt.Sprintf("node-%d", i)
		topicname := fmt.Sprintf("topic-%d", 0)
		msg := fmt.Sprintf("msg-%d", i)
		ps.Publish(topicname, msg)
	}

	signal.Notify(signalch, os.Interrupt, os.Kill, syscall.SIGTERM)
	signalType := <-signalch
	signal.Stop(signalch)
	fmt.Printf("On Signal <%s>\n", signalType)
}
