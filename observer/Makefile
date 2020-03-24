.PHONY: build-image build run sh

DOCKER_CMD ?= sh

docker_run_args =--rm -it
# run_cmd = bash
ifeq ($(IS_CI), true)
	docker_run_args =--rm -i
	# run_cmd = $(DOCKER_CMD)
endif

# Docker commands
build-image:
	docker build -t observer .

build:
	$(MAKE) build-image
	-docker rm observer-build
	docker create --name observer-build observer
	mkdir -p build
	docker cp observer-build:/app/observer build/observer
	docker rm observer-build

sh:
	$(MAKE) docker_cmd DOCKER_CMD="sh"

docker_cmd: build-image
	docker run $(docker_run_args) -e USER_ID=`id -u` -v "$(PWD)":"/src" observer $(DOCKER_CMD)

