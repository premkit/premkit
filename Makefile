.PHONY: clean test build run swagger-spec docker shell package_docker all

ifeq ($(BUILD_VERSION),)
BUILD_VERSION := 0.0.1
endif

.PHONY: clean
clean:
	rm -rf ./bin ./deploy/bin

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	mkdir -p ./bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
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

grype-install:
	curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh | sh -s -- -b .

scan: IMAGE=registry.replicated.com/library/premkit:local
scan: build grype-install
	docker build --pull -t ${IMAGE} -f ./deploy/Dockerfile .
	./grype --fail-on=medium --only-fixed -vv ${IMAGE}

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
