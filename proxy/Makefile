.PHONY: build publish


all: build


build:
	docker build -t omnition/proxy .

publish: build
	docker tag omnition/proxy:latest $(DOCKER_REGISTRY)/omnition/proxy:latest
	docker push $(DOCKER_REGISTRY)/omnition/proxy:latest
