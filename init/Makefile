.PHONY: build publish

all: build

build:
	docker build -t omnition/omnition-sidecar-init:v1 .

publish: build
	docker tag omnition/sidecar-init:latest $(DOCKER_REGISTRY)/omnition/sidecar-init:v1
	docker push $(DOCKER_REGISTRY).dkr.ecr.us-west-2.amazonaws.com/omnition/sidecar-init:v1
