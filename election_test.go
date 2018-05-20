package election_test

import (
	"context"
	"log"

	election "github.com/travisjeffery/kube-leaderelection"
)

func ExampleElection() {
	// create new elector
	elector, err := election.NewLeaderElector()
	if err != nil {
		panic(err)
	}
	l := &listener{}
	// register a listener to be called on election events
	elector.Register(l)
	// add this group to the nodes taking part in the election
	elector.Run(context.Background())
}

type listener struct {
}

func (l *listener) StartedLeading() {
	log.Print("started leading")
}

func (l *listener) StoppedLeading() {
	log.Print("stopped leading")
}

func (l *listener) NewLeader(id string) {
	logq.Printf("new leader: %s", id)
}
