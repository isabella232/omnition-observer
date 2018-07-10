.PHONY: build publish


all: build

build:
	$(MAKE) -C init build
	$(MAKE) -C sidecar build

publish: build
	$(MAKE) -C init publish
	$(MAKE) -C sidecar publish
