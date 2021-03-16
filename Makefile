#!/usr/bin/make -f
.PHONY: help bin test clean

include $(CURDIR)/Config.mk

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

bin: ## build a binary
ifdef CI
	$(MAKE) -f Build.mk bin
else
	$(DOCKER_RUN) $(CENTOS_CONTAINER) make -f Build.mk bin
	sudo chown -R $(shell id -u):$(shell id -g) bin
endif

test: ## run tests
ifdef CI
	$(MAKE) -f Build.mk test
else
	$(DOCKER_RUN) $(CENTOS_CONTAINER) make -f Build.mk test
	sudo chown -R $(shell id -u):$(shell id -g) c.out
endif

clean: ## clean all artifacts
	-rm -rf bin/
	-rm -rf $(RPMBUILD)
	-rm -f cc-test-reporter
	-rm -f c.out
