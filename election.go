package election

import (
	"context"
	"os"
	"sync"
	"time"

	"golang.org/x/build/kubernetes/api"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

type Config struct {
	NodeID        string
	Namespace     string
	RenewDeadline time.Duration
	LeaseDeadline time.Duration
	RetryPeriod   time.Duration
	LockName      string
	ComponentName string
}

// LeaderElector gives you leader election built on Kubernetes/etcd's APIs.
type LeaderElector struct {
	sync.Mutex
	leaderID  string
	leader    bool
	config    Config
	elector   *leaderelection.LeaderElector
	listeners map[Listener]struct{}
}

// NewLeaderElector creates a new leader elector instance. Takes an optional Config, by default the
// hostname is for this node's identifier.
func NewLeaderElector(args ...interface{}) (*LeaderElector, error) {
	var c Config
	if args == nil {
		c = Config{}
	} else {
		c = args[0].(Config)
	}
	var err error
	if c.NodeID == "" {
		c.NodeID, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}
	if c.Namespace == "" {
		c.Namespace = api.NamespaceDefault
	}
	if c.LockName == "" {
		c.LockName = "leader-election"
	}
	if c.ComponentName == "" {
		c.ComponentName = "leader-elector"
	}
	if c.RenewDuration == 0 {
		c.RenewDuration = time.Second * 5
	}
	if c.LeaseDuration == 0 {
		c.LeaseDuration = time.Second * 5
	}
	if c.RetryPeriod == 0 {
		c.RetryPeriod = time.Second * 2
	}

	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(c.CoreV1().RESTClient()).Events("")})
	recorder := broadcaster.NewRecorder(legacyscheme.Scheme, api.EventSource{Component: c.ComponentName})

	rl, err := resourcelock.New(
		resourcelock.ConfigMapsResourceLock,
		c.Namespace,
		c.LockName,
		c.CoreV1(),
		resourcelock.ResourceLockConfig{
			Identity:      c.NodeID,
			EventRecorder: recorder,
		})
	if err != nil {
		return nil, err
	}

	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: e.startedLeading,
		OnStoppedLeading: e.stoppedLeading,
		OnNewLeader:      e.newLeader,
	}

	config := leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: c.LeaseDuration,
		RenewDeadline: c.RenewDeadline,
		RetryPeriod:   c.RetryPeriod,
		Callbacks:     callbacks,
	}

	elector, err := leaderelection.NewLeaderElector(config)
	if err != nil {
		return nil, err
	}

	return &LeaderElector{elector: elector, listeners: make(map[Listener]struct{})}, nil
}

// Listener is an interface for the methods you need to implement to listen and handle leader
// election events.
type Listener interface {
	StartedLeading()
	StoppedLeading()
	NewLeader(id string)
}

// Register registers a listener to be called on leader election events.
func (e *LeaderElector) Register(l Listener) {
	e.Lock()
	defer e.Unlock()
	e.listeners[l] = struct{}{}
}

// Register deregisters a listener from being called on leader election events.
func (e *LeaderElector) Deregister(l Listener) {
	e.Lock()
	defer e.Unlock()
	delete(e.listeners, l)
}

// Run adds this node in the leader election group, starting the leader election process for this
// node.
func (e *LeaderElector) Run(ctx context.Context) {
	wait.Until(e.elector.Run, 0, ctx.Done())
}

func (e *LeaderElector) startedLeading(stop <-chan struct{}) {
	e.Lock()
	defer e.Unlock()
	e.leader = true
	e.leaderID = e.config.NodeID
	for listener := range e.listeners {
		listener.StartedLeading()
	}
	close(stop)
}

func (e *LeaderElector) stoppedLeading() {
	e.Lock()
	defer e.Unlock()
	e.leader = false
	e.leaderID = ""
	for listener := range e.listeners {
		listener.StoppedLeading()
	}
}

func (e *LeaderElector) newLeader(identity string) {
	e.Lock()
	defer e.Unlock()
	e.leaderID = identity
	for listener := range e.listeners {
		listener.NewLeader(identity)
	}
}
