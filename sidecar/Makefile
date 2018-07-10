.PHONY: build publish


all: build


build:
	docker build -t omnition/sidecar .

publish: build
	docker tag omnition/sidecar:latest $(DOCKER_REGISTRY)/omnition/sidecar:latest
	docker push $(DOCKER_REGISTRY)/omnition/sidecar:latest
