GO_MODULE := github.com/fnatte/pizza-mouse

PROTOS := $(wildcard protos/*.proto)
PBGO := $(patsubst protos/%.proto,internal/%.pb.go,$(PROTOS))

internal/%.pb.go: protos/%.proto
	protoc -I=protos/ --go_out=./ --go_opt=module=$(GO_MODULE) $?

.PHONY: build-api
build-api: $(PBGO)
	go build -o out/pizza-mouse-api github.com/fnatte/pizza-mouse/cmd/api

.PHONY: build-worker
build-worker:
	go build -o out/pizza-mouse-worker github.com/fnatte/pizza-mouse/cmd/worker

.PHONY: build-updater
build-updater: $(PBGO)
	go build -o out/pizza-mouse-updater github.com/fnatte/pizza-mouse/cmd/updater

build: build-api build-worker build-updater

start-api: build-api
	out/pizza-mouse-api

start-worker: build-worker
	out/pizza-mouse-worker

start-updater: build-updater
	out/pizza-mouse-updater

start: start-api start-worker start-updater
