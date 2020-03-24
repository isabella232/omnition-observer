.PHONY: build publish

all: build

build:
	docker build -t omnition/omnition-observer-init .

publish:
	@echo "publicly publishing omnition-observer-init version: " ${PROXY_VERSION}
	docker tag omnition/omnition-observer-init:latest omnition/omnition-observer-init:${PROXY_VERSION}
	docker push omnition/omnition-observer-init:${PROXY_VERSION}
