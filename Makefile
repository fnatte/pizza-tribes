GO_MODULE := github.com/fnatte/pizza-tribes

PROTOS := $(wildcard protos/*.proto)
PBGO := $(patsubst protos/%.proto,internal/%.pb.go,$(PROTOS))

internal/%.pb.go: protos/%.proto
	protoc -I=protos/ --go_out=./ --go_opt=module=$(GO_MODULE) $?

.PHONY: build-api
build-api: $(PBGO)
	go build -o out/pizza-tribes-api github.com/fnatte/pizza-tribes/cmd/api

.PHONY: build-worker
build-worker:
	go build -o out/pizza-tribes-worker github.com/fnatte/pizza-tribes/cmd/worker

.PHONY: build-updater
build-updater: $(PBGO)
	go build -o out/pizza-tribes-updater github.com/fnatte/pizza-tribes/cmd/updater

build: build-api build-worker build-updater

start-api: build-api
	out/pizza-tribes-api

start-worker: build-worker
	out/pizza-tribes-worker

start-updater: build-updater
	out/pizza-tribes-updater

start: start-api start-worker start-updater
