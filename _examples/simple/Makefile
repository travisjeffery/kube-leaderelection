server: deps
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o server main.go

container: server
	@docker build -t travisjeffery/kube-leaderelection-example:0.0.1 .

run: container
	@kubectl run kube-leaderelection-example --image travisjeffery/kube-leaderelection-example:0.0.1 --replicas=3

clean:
	@rm -f server
	@kubectl delete deploy kube-leaderelection-example

deps:
	@dep ensure
