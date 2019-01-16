.PHONY: build publish

all: build

build:
	$(MAKE) -C init build
	$(MAKE) -C proxy build

publish: build
	@echo "Publishing version: " ${PROXY_VERSION_TO_USE}
	$(MAKE) -C init publish
	$(MAKE) -C proxy publish
