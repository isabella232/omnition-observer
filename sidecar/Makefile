.PHONY: build publish


all: build


build:
	docker build -t omnition/sidecar:v1 .

publish: build
	docker tag omnition-sidecar:latest $(DOCKER_REGISTRY)/omnition-sidecar:v1
	docker push $(DOCKER_REGISTRY).dkr.ecr.us-west-2.amazonaws.com/omnition-sidecar:v1
