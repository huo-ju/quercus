package main

import (
	"fmt"
	"github.com/huo-ju/quercus/pkg/pubsub"
	"github.com/huo-ju/quercus/pkg/quality"
	"os"
	"os/signal"
	"syscall"
)

var (
	signalch chan os.Signal
)

func main() {

	signalch = make(chan os.Signal, 1)
	NodeNumber := 5

	ps := pubsub.NewPubsub()

	listener := func(topicname string, sub *pubsub.Subscriber) {
		dqa := quality.NewDelayQualityAgent(0, 10)
		for i := range sub.Chan {
			out := dqa.Pass(i)
			if out != nil {
				fmt.Printf("Node:(%s) [%s] got %s\n", sub.Id, topicname, out)
			}
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
		topicname := fmt.Sprintf("topic-%d", 0)
		msg := fmt.Sprintf("msg-%d", i)
		ps.Publish(topicname, msg)
	}

	signal.Notify(signalch, os.Interrupt, os.Kill, syscall.SIGTERM)
	signalType := <-signalch
	signal.Stop(signalch)
	fmt.Printf("On Signal <%s>\n", signalType)
}
