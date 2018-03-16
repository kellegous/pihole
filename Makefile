export PATH := $(shell pwd)/.build/bin:$(PATH)
export GOPATH := $(shell pwd)

export GOOS ?= $(shell go env GOOS)
export GOARCH ?= $(shell go env GOARCH)

DESC := $(GOOS)-$(GOARCH)

SHA := $(shell git rev-parse HEAD)
REF := $(shell git rev-parse --abbrev-ref HEAD)

BUILD_FLAGS := -ldflags "-X pihole/build.SHA=$(SHA) -X pihole/build.Ref=$(REF)"

ALL: bin/$(DESC)/deploy bin/$(DESC)/client bin/$(DESC)/server

.build/bin/protoc-gen-go:
	GOPATH=$(shell pwd)/.build go get github.com/golang/protobuf/...

src/pihole/api/api.pb.go: src/pihole/api/api.proto .build/bin/protoc-gen-go
	protoc -Isrc/pihole $< --go_out=plugins=grpc:src/pihole

bin/$(DESC)/client: src/pihole/api/api.pb.go $(shell find . -name '*.go')
	go build $(BUILD_FLAGS) -o $@ pihole/client

bin/$(DESC)/server: src/pihole/api/api.pb.go $(shell find . -name '*.go')
	go build $(BUILD_FLAGS) -o $@ pihole/server

bin/$(DESC)/deploy: $(shell find src/pihole/deploy -name '*.go')
	go build $(BUILD_FLAGS) -o $@ pihole/deploy

deploy: bin/$(DESC)/deploy
	bin/$(DESC)/deploy

clean:
	rm -rf bin