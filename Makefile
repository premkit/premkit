.PHONY: clean install test build run swagger-spec docker shell build_docker package_docker all

PREMKIT_TAG ?= 0.0.0

clean:
	rm -rf ./bin ./deploy/bin

install:
	govendor install +std +local +vendor,^program

test:
	govendor test +local

build:
	mkdir -p ./bin
	go build -o ./bin/premkit .

build_ci:
	mkdir -p ./bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -i \
	--ldflags=" \
	-X github.com/premkit/premkit/version.version=${PREMKIT_TAG} \
	-X github.com/premkit/premkit/version.gitSHA=${BUILD_SHA} \
	-X \"github.com/premkit/premkit/version.buildTime=$(shell date)\" \
	" \
	-o ./bin/premkit .

run:
	./bin/premkit daemon

swagger-spec:
	mkdir -p ./spec/v1
	swagger generate spec -w ./handlers/v1 -o ./spec/v1/swagger.json  ./...
	swagger validate ./spec/v1/swagger.json

docker:
	docker build -t premkit/premkit:dev .

docker_build:
	docker build --rm=false -t premkit/premkit:build -f deploy/Dockerfile-build .

shell:
	docker run --rm -it -P --name premkit \
		-p 80:80 \
		-p 443:443 \
		-v `pwd`:/go/src/github.com/premkit/premkit \
		-v `pwd`/data:/data \
		premkit/premkit:dev

build_docker:
	mkdir -p ./deploy/bin
	go build \
		--ldflags '-extldflags "-static"' \
		-o ./deploy/bin/premkit .

build_docker_local:
	mkdir -p ./deploy/bin
	docker run --rm -it \
		-v `pwd`:/go/src/github.com/premkit/premkit \
		premkit/premkit:build go build \
			--ldflags '-extldflags "-static"' \
			-o ./deploy/bin/premkit .

package_docker:
	docker build --rm=false -t premkit/premkit:$(PREMKIT_TAG) -f ./deploy/Dockerfile .

all: build test
