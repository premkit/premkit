.PHONY: clean install test build run swagger-spec docker shell package_docker all

ifeq ($(BUILD_VERSION),)
BUILD_VERSION := 0.0.1
endif

.PHONY: clean
clean:
	rm -rf ./bin ./deploy/bin

.PHONY: install
install:
	govendor install +std +local +vendor,^program

.PHONY: test
test:
	govendor test +local

.PHONY: build
build:
	mkdir -p ./bin

.PHONY: build_ci
build_ci:
	mkdir -p ./bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -i \
	--ldflags=" \
	-X github.com/premkit/premkit/version.version=$(BUILD_VERSION) \
	-X github.com/premkit/premkit/version.gitSHA=${BUILD_SHA} \
	-X \"github.com/premkit/premkit/version.buildTime=$(shell date)\" \
	" \
	-o ./bin/premkit .

.PHONY: run
run:
	./bin/premkit daemon

.PHONY: swagger_spec
swagger_spec:
	mkdir -p ./spec/v1
	swagger generate spec -w ./handlers/v1 -o ./spec/v1/swagger.json  ./...
	swagger validate ./spec/v1/swagger.json

.PHONY: docker
docker:
	docker build -t premkit/premkit:dev .

.PHONY: shell
shell:
	docker run --rm -it -P --name premkit \
		-p 80:80 \
		-p 443:443 \
		-v `pwd`:/go/src/github.com/premkit/premkit \
		-v `pwd`/data:/data \
		premkit/premkit:dev

.PHONY: all
all: build swagger_spec test
