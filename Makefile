default: build

all: build docker

build: build-linux

build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

build-windows:
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

build-macOS:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build

clean:
	rm -rf ./bin

docker:
	docker build --no-cache -t suslmk/tksinfo -f Dockerfile .