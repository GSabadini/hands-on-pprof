memory:
	go run memory_leak.go

memory-trace-gc:
	GODEBUG=gctrace=1 go run memory_leak.go

memory-init-trace:
	GODEBUG=init=1 go run memory_leak.go

run:
	docker-compose up -d

memory-docker:
	docker run --rm -it -w /app -v ${PWD}:/app -v ${GOPATH}/pkg/mod/cache:/go/pkg/mod/cache -p 6060:6060 golang:1.16-stretch go run memory_leak.go

goroutine:
	go run goroutine_leak.go

trace:
	$(shell wget -P ./tmp/ http://localhost:6061/debug/pprof/trace\?seconds\=30s)

profile:
	$(shell wget -P ./tmp/ http://localhost:6060/debug/pprof/profile\?seconds\=10s)

request:
	n=20; \
	while [ $${n} -gt 0 ] ; do \
		curl -v --header "Connection: keep-alive" "http://localhost:6060/leak-query"; \
		n=`expr $$n - 1`; \
	done; \
	true

pprof:
	pprof -alloc_space -http "localhost:6062" 'http://localhost:6061/debug/pprof/heap'