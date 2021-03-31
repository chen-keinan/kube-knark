SHELL := /bin/bash

GOCMD=go
GOPACKR=$(GOCMD) get -u github.com/gobuffalo/packr/packr && packr
GOMOD=$(GOCMD) mod
GOMOCKS=$(GOCMD) generate ./...
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=kube-knark

all:test lint build

fmt:
	$(GOCMD) fmt ./...
lint:
	./scripts/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOCMD) get github.com/golang/mock/mockgen@latest
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./cmd... ./internal... ./pkg... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	$(GOCMD) tool cover  -func coverage.md
build:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v
install:build
	cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
build_debug:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 go build -v -gcflags='-N -l'
dlv:
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./kube-knark

.PHONY: all build install test
