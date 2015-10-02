BUILD_DOCKER_ARGS= -v `pwd`:/home/go/src/weevil/server \
	  	  --workdir=/home/go/src/weevil/server \
		  -e GOPATH=/home/go
BUILD_FLAGS=-ldflags "-extldflags \"-static\"" -tags netgo

.PHONY: all clean

all: server.uptodate

build.uptodate:
	docker run --name=weevil-build $(BUILD_DOCKER_ARGS) \
		google/golang \
		go get -tags netgo .
	docker commit weevil-build weevil/build
	docker rm -f weevil-build
	touch build.uptodate

server.uptodate: weevil Dockerfile
	docker build -t weevil/server .
	touch server.uptodate

weevil: build.uptodate *.go
	docker run --rm $(BUILD_DOCKER_ARGS) \
		weevil/build go build $(BUILD_FLAGS) -o $@ .

clean:
	rm -f build.uptodate server.uptodate
	rm -f weevil
	docker rm -f weevil-build || true
