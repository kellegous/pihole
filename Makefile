export PATH := $(shell pwd)/.build/bin:$(PATH)
export GOPATH := $(shell pwd)

export GOOS ?= $(shell go env GOOS)
export GOARCH ?= $(shell go env GOARCH)

DESC := $(GOOS)-$(GOARCH)

ALL: bin/$(DESC)/deploy bin/$(DESC)/client bin/$(DESC)/server

.build/bin/protoc-gen-go:
	GOPATH=$(shell pwd)/.build go get github.com/golang/protobuf/...

src/pihole/api/api.pb.go: src/pihole/api/api.proto .build/bin/protoc-gen-go
	protoc -Isrc/pihole $< --go_out=plugins=grpc:src/pihole

bin/$(DESC)/client: src/pihole/api/api.pb.go $(shell find . -name *.go)
	go build -o $@ pihole/client

bin/$(DESC)/server: src/pihole/api/api.pb.go $(shell find . -name *.go)
	go build -o $@ pihole/server

bin/$(DESC)/deploy: $(shell find src/pihole/deploy -name '*.go')
	go build -o $@ pihole/deploy

deploy: bin/$(DESC)/deploy
	bin/$(DESC)/deploy

clean:
	rm -rf bin