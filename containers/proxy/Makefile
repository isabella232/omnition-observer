.PHONY: build publish

all: build

build:
	docker build -t omnition/proxy .

publish: build
	docker tag omnition/proxy:latest $(DOCKER_REGISTRY)/omnition/proxy:$(PROXY_VERSION)
	docker push $(DOCKER_REGISTRY)/omnition/proxy:$(PROXY_VERSION)

run: build
	docker run -it --rm --name omnition-proxy -e 'ZIPKIN_HOST="zipkin_host"' omnition/proxy:$(PROXY_VERSION) "show-config"

run-show-config: build
	docker run -it --rm --name omnition-proxy -e 'ZIPKIN_HOST="hello"' omnition/proxy:$(PROXY_VERSION) "show-config"
