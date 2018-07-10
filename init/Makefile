.PHONY: build publish

all: build

build:
	docker build -t omnition/sidecar-init .

publish: build
	docker tag omnition/sidecar-init:latest $(DOCKER_REGISTRY)/omnition/sidecar-init:latest
	docker push $(DOCKER_REGISTRY).dkr.ecr.us-west-2.amazonaws.com/omnition/sidecar-init:latest
