.PHONY: clean install test build run swagger-spec container shell build_docker package_docker all

clean:
	rm -f ./bin/premkit ./deploy/bin/premkit

install:
	govendor install

test:
	govendor test +local

build:
	mkdir -p ./bin
	govendor build -o ./bin/premkit .

run:
	./bin/premkit daemon

swagger-spec:
	mkdir -p ./spec/v1
	swagger generate spec -b github.com/premkit/premkit/handlers/v1 -o ./spec/v1/swagger.json
	swagger validate ./spec/v1/swagger.json

container:
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
	CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags '-extldflags "-static"' -o ./deploy/bin/premkit .

package_docker:
	docker build -t premkit/premkit:$(PREMKIT_TAG) -f ./deploy/Dockerfile .

all: build test
