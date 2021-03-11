SHELL := /bin/bash

GOCMD=go
MOVESANDBOX=mv kube-beacon ~/vagrant_file/.
GOPACKR=$(GOCMD) get -u github.com/gobuffalo/packr/packr && packr
GOMOD=$(GOCMD) mod
GOMOCKS=$(GOCMD) generate ./...
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=kube-beacon
GOCOPY=cp kube-beacon ~/vagrant_file/.

all:lint build

fmt:
	$(GOCMD) fmt ./...

build:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v
	$(MOVESANDBOX)
.PHONY: all build install test
