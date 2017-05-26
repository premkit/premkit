.PHONY: clean install test build run swagger-spec docker shell build_docker package_docker all

clean:
	rm -rf ./bin ./deploy/bin

install:
	govendor install +std +local +vendor,^program

test:
	govendor test +local

build:
	mkdir -p ./bin
	go build -o ./bin/premkit .

run:
	./bin/premkit daemon

swagger-spec:
	mkdir -p ./spec/v1
	swagger generate spec -b github.com/premkit/premkit/handlers/v1 -o ./spec/v1/swagger.json
	swagger validate ./spec/v1/swagger.json

docker:
	docker build -t premkit/premkit:dev .

shell:
	docker run --rm -it -P --name premkit \
	  -p 80:80 \
	  -p 443:443 \
	  -v `pwd`:/go/src/github.com/premkit/premkit \
	  -v `pwd`/data:/data \
	  premkit/premkit:dev

build_docker:
	mkdir -p ./package/bin
	go build -tags "netgo" --ldflags '-extldflags "-static"' -o ./deploy/bin/premkit .

package_docker:
	docker build -t premkit/premkit:$(PREMKIT_TAG) -f ./deploy/Dockerfile .

all: build test
