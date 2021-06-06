.PHONY: build build-amd64 build-arm64

build:
	GOOS=$(os) GOARCH=$(arch) go build -o bot .  

build-amd64:
	make build os=linux arch=amd64

build-arm64:
	make build os=linux arch=arm64
