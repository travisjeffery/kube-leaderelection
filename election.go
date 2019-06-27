package election

import (
	"context"
	"os"
	"sync"
	"time"

	api "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

// LeaderElector gives you leader election built on Kubernetes/etcd's APIs.
type LeaderElector struct {
	sync.Mutex
	leaderID  string
	leader    bool
	config    Config
	elector   *leaderelection.LeaderElector
	listeners map[Listener]struct{}
}

// Config is used to configure a LeaderElector.
type Config struct {
	// NodeID is the identifier for this node.
	NodeID string
	// Namespace is the namespace this elector should run in/store its lock.
	Namespace string
	// RenewDeadline is the duration that the acting master will retry
	// refreshing leadership before giving up.
	RenewDeadline time.Duration
	// LeaseDuration is the duration that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack.
	LeaseDuration time.Duration
	// RetryPeriod is the duration the LeaderElector clients should wait
	// between tries of actions.
	RetryPeriod time.Duration
	// ComponentName is used as the group name for this leader election group. If you run
	// multiple leader elector instances within a service you probably want to differentiate
	// them by name.
	ComponentName string
	// LockName is name of the resource lock. You probably don't ever need to set this.
	LockName string
	// Path to kubeconfig
	KubeConfig string
}

// NewLeaderElector creates a new leader elector instance. Takes an optional Config, by default the
// hostname is for this node's identifier.
func NewLeaderElector(args ...interface{}) (*LeaderElector, error) {
	var config Config
	if args == nil {
		config = Config{}
	} else {
		config = args[0].(Config)
	}
	var err error
	if config.NodeID == "" {
		config.NodeID, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}
	if config.Namespace == "" {
		config.Namespace = api.NamespaceDefault
	}
	if config.LockName == "" {
		config.LockName = "leader-election"
	}
	if config.ComponentName == "" {
		config.ComponentName = "leader-elector"
	}
	if config.LeaseDuration == 0 {
		config.LeaseDuration = time.Second * 10
	}
	if config.RenewDeadline == 0 {
		config.RenewDeadline = time.Second * 5
	}
	if config.RetryPeriod == 0 {
		config.RetryPeriod = time.Second * 3
	}

	// This function uses kubeconfig if specified. If kubeconfig is not passed in it fallbacks to inClusterConfig.
	clientConfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(clientset.CoreV1().RESTClient()).Events("")})
	recorder := broadcaster.NewRecorder(legacyscheme.Scheme, api.EventSource{Component: config.ComponentName})

	rl, err := resourcelock.New(
		resourcelock.ConfigMapsResourceLock,
		config.Namespace,
		config.LockName,
		clientset.CoreV1(),
		resourcelock.ResourceLockConfig{
			Identity:      config.NodeID,
			EventRecorder: recorder,
		})
	if err != nil {
		return nil, err
	}

	e := &LeaderElector{listeners: make(map[Listener]struct{})}

	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: e.startedLeading,
		OnStoppedLeading: e.stoppedLeading,
		OnNewLeader:      e.newLeader,
	}

	e.elector, err = leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: config.LeaseDuration,
		RenewDeadline: config.RenewDeadline,
		RetryPeriod:   config.RetryPeriod,
		Callbacks:     callbacks,
	})
	if err != nil {
		return nil, err
	}

	return e, nil
}

// Listener is an interface for the methods you need to implement to listen and handle leader
// election events.
type Listener interface {
	// StartedLeading is invoked when this node becomes the leader.
	StartedLeading()
	// StoppedLeading is invoked when this node stops being the leader.
	StoppedLeading()
	// NewLeader is invoked when a new leader is elected.
	NewLeader(id string)
}

// IsLeader returns true if this node is the leader.
func (e *LeaderElector) IsLeader() bool {
	return e.elector.IsLeader()
}

// GetLeader returns the ID of the leader.
func (e *LeaderElector) GetLeader() string {
	return e.elector.GetLeader()
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
	e.elector.Run(ctx)
}

func (e *LeaderElector) startedLeading(ctx context.Context) {
	e.Lock()
	defer e.Unlock()
	e.leader = true
	e.leaderID = e.config.NodeID
	for listener := range e.listeners {
		listener.StartedLeading()
	}
	<-ctx.Done()
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
