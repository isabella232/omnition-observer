.PHONY: build publish build-envoy

all: build

build: build-envoy build-k8s-init build-container-proxy

build-envoy:
	@$(MAKE) -C envoy_filter build
	@mkdir -p containers/proxy/bin
	@cp envoy_filter/build/envoy containers/proxy/bin/

build-k8s-init:
	@$(MAKE) -C containers/init build

build-container-proxy:
	@$(MAKE) -C containers/proxy build

publish: build
	@$(MAKE) -C containers/init publish
	@$(MAKE) -C containers/proxy publish

