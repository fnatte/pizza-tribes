GO_MODULE := github.com/fnatte/pizza-mouse

PROTOS := $(wildcard proto/*.proto)
PBGO := $(patsubst proto/%.proto,internal/%.pb.go,$(PROTOS))

internal/%.pb.go: proto/%.proto
	protoc -I=proto/ --go_out=./ --go_opt=module=$(GO_MODULE) $?

.PHONY: build-api
build-api: $(PBGO)
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
