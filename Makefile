PROXY_VERSION?=$$(cat VERSION)

.PHONY: build publish

all: build

build: build-observer build-k8s-init build-container-proxy

build-observer:
	@$(MAKE) -C observer build
	@mkdir -p containers/proxy/bin
	@cp observer/build/observer containers/proxy/bin/

build-k8s-init:
	@$(MAKE) -C containers/init build

build-container-proxy:
	@$(MAKE) -C containers/proxy build

publish: build
	@echo "Publishing version: " $(PROXY_VERSION)
	@$(MAKE) -C containers/init publish PROXY_VERSION=${PROXY_VERSION}
	@$(MAKE) -C containers/proxy publish PROXY_VERSION=${PROXY_VERSION}
