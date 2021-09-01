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
