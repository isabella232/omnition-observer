.PHONY: echo-mode docker-run docker-exec docker-stop sh cmd build test

DOCKER_CMD=bash

CONTAINER_NAME?=envoy-build
BUILD_IMAGE_TAG?=55606ec3f0810683e9cab31b2f9db2f2cdb43ac0
BUILD_IMAGE?=envoyproxy/envoy-build:$(BUILD_IMAGE_TAG)

IS_RUNNING = $(shell docker inspect -f {{.State.Running}} $(CONTAINER_NAME))

ifeq ($(IS_RUNNING), true)
	MODE = docker-exec
else
	MODE = docker-run
endif

all: build

echo-mode:
	@echo $(MODE)

pull-build-image:
	docker pull $(BUILD_IMAGE)

docker-run:
	docker run -d --name="$(CONTAINER_NAME)" -v "$(CURDIR)":/source -v "$(CURDIR)/caches/home/.cache":"/root/.cache" -e HOST_USERID=`id -u` -e HOST_GROUPID=`id -g` $(BUILD_IMAGE) bash
	$(MAKE) docker-exec DOCKER_CMD=$(DOCKER_CMD)

docker-exec:
	docker exec $(CONTAINER_NAME) "$(DOCKER_CMD)"

docker-stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

sh:
	@$(MAKE) $(MODE) COMMAND="bash"

cmd:
	@$(MAKE) $(MODE) COMMAND="$(DOCKER_CMD)"

build:
	@$(MAKE) $(MODE) DOCKER_CMD=/source/scripts/build_envoy.sh
	@chmod 775 build/envoy

test:
	@$(MAKE) $(MODE) DOCKER_CMD=/source/scripts/test_envoy.sh
