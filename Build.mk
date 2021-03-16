#!/usr/bin/make -f
.PHONY: bin test link vet goimports 

include $(CURDIR)/Config.mk

bin: test ## build just binary
	mkdir -p bin
	$(GO) build \
		-o bin/$(BINNAME) *.go

test: goimports lint vet ## run unit tests
	$(GO) test -coverprofile c.out ./...

lint: ## run golint
	golint -set_exit_status ./...

vet: ## run go vet
	$(GO) vet ./...

goimports: ## run goimports
	goimports -l ./ | xargs -r false

