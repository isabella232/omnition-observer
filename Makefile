.PHONY: build publish

all: build

build:
	$(MAKE) -C init build
	$(MAKE) -C proxy build

publish: build
	$(MAKE) -C init publish
	$(MAKE) -C proxy publish
