.PHONY: build publish

all: build

build:
	docker build -t omnition/omnition-observer .

publish: build
	docker tag omnition/omnition-observer:latest omnition/omnition-observer:$(PROXY_VERSION)
	docker push $(DOCKER_REGISTRY)/omnition/omnition-observer:$(PROXY_VERSION)

publish-public: build
	docker tag omnition/omnition-observer:latest omnition/omnition-observer:$(PROXY_VERSION)
	docker push omnition/omnition-observer:$(PROXY_VERSION)

run: build
	docker run -it --rm --name omnition-observer -e 'ZIPKIN_HOST="zipkin_host"' omnition/omnition-observer:$(PROXY_VERSION) "show-config"

run-show-config: build
	docker run -it --rm --name omnition-observer -e 'ZIPKIN_HOST="hello"' omnition/omnition-observer:$(PROXY_VERSION) "show-config"
