#!/usr/bin/make -f

GO               := go
DOCKER_WORKDIR   := /go/src/github.com/moxspec/moxspec-metrics-agent
DOCKER_RUN       := sudo docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR) --workdir=$(DOCKER_WORKDIR)
CENTOS_CONTAINER := actapio/moxspec-centos:7
BINNAME          := mox-metrics-agent
