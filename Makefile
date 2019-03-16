GOPATH=`pwd`

all: deps build

build:
	@GOPATH=$(GOPATH) go build -o resizefoldericon src/github.com/ScottESanDiego/resizefoldericon/resizefoldericon.go
	strip resizefoldericon

deps:
	@GOPATH=$(GOPATH) go get github.com/nfnt/resize

clean:
	rm resizefoldericon
	rm -rf pkg
