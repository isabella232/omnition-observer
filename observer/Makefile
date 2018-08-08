.PHONY: build-image build run sh go_build

DOCKER_CMD ?= sh

docker_run_args =--rm -it
# run_cmd = bash
ifeq ($(IS_CI), true)
	docker_run_args =--rm -i
	# run_cmd = $(DOCKER_CMD)
endif

# Docker commands
build-image:
	docker build -t observer-build .

build:
	$(MAKE) docker_cmd DOCKER_CMD="make go_build"

sh:
	$(MAKE) docker_cmd DOCKER_CMD="sh"

docker_cmd: build-image
	docker run $(docker_run_args) -e USER_ID=`id -u` -v "$(PWD)":"/src" observer-build $(DOCKER_CMD)

# Native commands
go_clean:
	rm -f ./build/observer

go_build: go_clean
	go build -o build/observer ./cmd/observer

go_run: go_build
	./build/observer

go_test:
	go test -v ./cmd/observer

