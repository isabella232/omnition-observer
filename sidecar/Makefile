.PHONY: build publish


all: build


build:
	docker build -t omnition/sidecar:v1 .

publish:
	docker push omnition/sidecar:v1
