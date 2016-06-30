.PHONY: clean godep deps run test build all

clean:
	rm -f ./bin/premkit

run:
	./bin/premkit daemon

test:
	govendor test +local

build:
	mkdir -p ./bin
	govendor build -o ./bin/premkit .

shell:
	docker run --rm -it -P --name premkit \
                -p 80:80 \
                -p 443:443 \
		-v `pwd`:/go/src/github.com/premkit/premkit \
                -v `pwd`/data:/data \
		premkit/premkit:dev

all: build test
