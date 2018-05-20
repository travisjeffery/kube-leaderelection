package main

import (
	"context"
	"flag"
	"log"
	"os"

	election "github.com/travisjeffery/kube-leaderelection"
)

func main() {
	flag.Parse()
	elector, err := election.NewLeaderElector()
	if err != nil {
		panic(err)
	}
	elector.Register(&listener{})
	elector.Run(context.Background())
}

type listener struct {
}

func (l *listener) StartedLeading() {
	log.Printf("[INFO] %s: started leading", hostname())
}

// invoked when this node stops being the leader
func (l *listener) StoppedLeading() {
	log.Printf("[INFO] %s: stopped leading", hostname())
}

// invoked when a new leader is elected
func (l *listener) NewLeader(id string) {
	log.Printf("[INFO] %s: new leader: %s", hostname(), id)
}

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname
}
