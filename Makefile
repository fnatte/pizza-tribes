GO_MODULE := github.com/fnatte/pizza-tribes
GO_FLAGS := -ldflags '-s -w -extldflags "-static"'

PROTOS := $(wildcard protos/*.proto)
PBGO := $(patsubst protos/%.proto,internal/%.pb.go,$(PROTOS))

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

internal/%.pb.go: protos/%.proto
	protoc -I=protos/ --go_out=./ --go_opt=module=$(GO_MODULE) --experimental_allow_proto3_optional $?

.PHONY: generate
generate: $(PBGO)

.PHONY: build-api
build-api: $(PBGO)
	go build $(GO_FLAGS) -o out/pizza-tribes-api github.com/fnatte/pizza-tribes/cmd/api

.PHONY: build-worker
build-worker:
	go build $(GO_FLAGS) -o out/pizza-tribes-worker github.com/fnatte/pizza-tribes/cmd/worker

.PHONY: build-updater
build-updater: $(PBGO)
	go build $(GO_FLAGS) -o out/pizza-tribes-updater github.com/fnatte/pizza-tribes/cmd/updater

.PHONY: build-migrator
build-migrator: $(PBGO)
	go build $(GO_FLAGS) -o out/pizza-tribes-migrator github.com/fnatte/pizza-tribes/cmd/migrator

.PHONY: build-admin
build-admin: $(PBGO)
	go build $(GO_FLAGS) -o out/pizza-tribes-admin github.com/fnatte/pizza-tribes/cmd/admin

build: build-api build-worker build-updater build-migrator build-admin

start-migrator: build-migrator
	out/pizza-tribes-migrator

start-api: build-api
	out/pizza-tribes-api

start-worker: build-worker
	out/pizza-tribes-worker

start-updater: build-updater
	out/pizza-tribes-updater

start-admin: build-admin
	out/pizza-tribes-admin serve

start: start-api start-worker start-updater start-migrator start-admin
