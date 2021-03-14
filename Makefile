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

all:lint build

fmt:
	$(GOCMD) fmt ./...
lint:
	./scripts/lint.sh
tidy:
	$(GOMOD) tidy -v
build:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v
	$(MOVESANDBOX)
.PHONY: all build install test
