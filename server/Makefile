BUILD_DOCKER_ARGS=-v `pwd`:/home/go/src/weevil/server \
		  -e CGO_ENABLED=0 \
		  --workdir=/home/go/src/weevil/server \
		  -e GOPATH=/home/go
BUILD_FLAGS=-ldflags "-s -extldflags \"-static\"" -tags netgo -a

STATIC_FILES:=../js

.PHONY: all clean run

all: server.uptodate

build.uptodate:
	docker run --name=weevil-build $(BUILD_DOCKER_ARGS) \
		google/golang \
			sh -c "go clean -i net && \
			go install -tags netgo std && \
			go get -tags netgo ."
	docker commit weevil-build weevil/build
	docker rm -f weevil-build
	touch build.uptodate

server.uptodate: weevil Dockerfile
	mkdir -p res
	cp -R $(STATIC_FILES) res/
	cp ../index.html ./
	docker build -t weevil/server .
	touch server.uptodate

weevil: build.uptodate *.go
	docker run --rm $(BUILD_DOCKER_ARGS) \
		weevil/build go build $(BUILD_FLAGS) -o $@ .

clean:
	rm -f build.uptodate server.uptodate
	rm -f weevil
	docker rm -f weevil-build || true

run: server.uptodate
	docker run --rm -p 7070:7070 weevil/server
