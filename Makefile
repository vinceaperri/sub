all: i386 amd64

i386:
	GOARCH=386 go build -o sub-i386

amd64:
	GOARCH=amd64 go build -o sub-amd64
