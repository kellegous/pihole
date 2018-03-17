export PATH := $(shell pwd)/.build/bin:$(PATH)
export GOPATH ?= $(shell cd ../../../.. && pwd)

export GOOS ?= $(shell go env GOOS)
export GOARCH ?= $(shell go env GOARCH)

DESC := $(GOOS)-$(GOARCH)

SHA := $(shell git rev-parse HEAD)
REF := $(shell git rev-parse --abbrev-ref HEAD)

BUILD_FLAGS := -ldflags "-X github.com/kellegous/pihole/build.SHA=$(SHA) -X github.com/kellegous/pihole/build.Ref=$(REF)"

ALL: bin/$(DESC)/deploy bin/$(DESC)/client bin/$(DESC)/server

.build/bin/protoc-gen-go:
	GOPATH=$(shell pwd)/.build go get github.com/golang/protobuf/...

api/api.pb.go: api/api.proto .build/bin/protoc-gen-go
	protoc -I. $< --go_out=plugins=grpc:.

bin/$(DESC)/client: api/api.pb.go $(shell find . -name '*.go')
	go build $(BUILD_FLAGS) -o $@ github.com/kellegous/pihole/client

bin/$(DESC)/server: api/api.pb.go $(shell find . -name '*.go')
	go build $(BUILD_FLAGS) -o $@ github.com/kellegous/pihole/server

bin/$(DESC)/deploy: $(shell find deploy -name '*.go')
	go build $(BUILD_FLAGS) -o $@ github.com/kellegous/pihole/deploy

deploy: bin/$(DESC)/deploy
	bin/$(DESC)/deploy

clean:
	rm -rf bin