SHELL := /bin/bash

GOCMD=go
MOVESANDBOX=mv kube-knark ~/vagrant_file/.
GOPACKR=$(GOCMD) get -u github.com/gobuffalo/packr/packr && packr
GOMOD=$(GOCMD) mod
GOMOCKS=$(GOCMD) generate ./...
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=kube-knark
GOCOPY=cp kube-knark ~/vagrant_file/.

all:test lint build

fmt:
	$(GOCMD) fmt ./...
lint:
	./scripts/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	@go get github.com/golang/mock/mockgen@latest
	@go install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	@go generate ./...
	$(GOTEST) ./... -coverprofile cp.out
build:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v

build_debug:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 go build -v -gcflags='-N -l'
dlv:
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./kube-knark

.PHONY: all build install test
