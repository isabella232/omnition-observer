.PHONY: build publish

all: build

build:
	docker build -t omnition/proxy-init .

publish: build
	@echo "publishing proxy-init version: " ${PROXY_VERSION}
	docker tag omnition/proxy-init:latest $(DOCKER_REGISTRY)omnition/proxy-init:${PROXY_VERSION}
	docker push $(DOCKER_REGISTRY)omnition/proxy-init:${PROXY_VERSION}
