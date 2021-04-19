.PHONY: build-api
build-api:
	go build -o out/pizza-mouse-api github.com/fnatte/pizza-mouse/cmd/api

.PHONY: build-worker
build-worker:
	go build -o out/pizza-mouse-worker github.com/fnatte/pizza-mouse/cmd/worker

build: build-api build-worker

start-api: build-api
	out/pizza-mouse-api

start-worker: build-worker
	out/pizza-mouse-worker

start: start-api start-worker
