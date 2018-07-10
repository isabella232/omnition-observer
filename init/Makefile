.PHONY: build publish


all: build


build:
	docker build -t omnition/sidecar-init:v1 .

publish:
	docker push omnition/sidecar-init:v1
