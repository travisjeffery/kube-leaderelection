# kube-leaderelection

Leader election for your services, using Kubernetes' APIs/etcd.

## Example

``` go
// create new elector
elector, err := election.NewLeaderElector()
if err != nil {
	panic(err)
}
// register a listener to be called on election events
elector.Register(listener)
// add this group to the nodes taking part in the election
elector.Run(ctx)
```

Your listener(s) looks something like this. You can register as many as you want.

``` go
// invoked when this node becomes the leader
func (l *listener) StartedLeading() {
	log.Print("started leading")
}

// invoked when this node stops being the leader
func (l *listener) StoppedLeading() {
	log.Print("stopped leading")
}

// invoked when a new leader is elected
func (l *listener) NewLeader(id string) {
	logq.Printf("new leader: %s", id)
}

```

## License

MIT

---

- Twitter [@travisjeffery](https://twitter.com/travisjeffery)
- Medium [@travisjeffery](https://medium.com/@travisjeffery)
- GitHub [@travisjeffery](https://github.com/travisjeffery)
- [travisjeffery.com](http://travisjeffery.com)
