.PHONY: build publish

all: build

build:
	docker build -t omnition/proxy-init .

publish: build
	docker tag omnition/proxy-init:latest $(DOCKER_REGISTRY)/omnition/proxy-init:latest
	docker push $(DOCKER_REGISTRY)/omnition/proxy-init:latest
