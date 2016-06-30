premkit
=======

## Setup
We use Docker for the official build environment.  Docker is not required to run the binary, but to simplify the development and build enviroment, there are 
Dockerfiles to use if you want to build the binary.

## Build the development container
```shell
$ docker build -t premkit/premkit:dev .
```

## Run the development environment

### Start the container, build the execuatable, and run the service.
```
$ make shell
$ make build run
```

## Run tests
Tests are automatically run in CircleCI after pushing.  Tests can be run manually with
```shell
$ make shell
$ make test
```
