build-i386:
	GOARCH=386 go build -o sub-i386

build-amd64:
	GOARCH=amd64 go build -o sub-amd64

build: build-i386 build-amd64
