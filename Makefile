GOPATH=`pwd`

all: deps build

build:
	@GOPATH=$(GOPATH) go build -o resizefoldericon src/github.com/ScottESanDiego/resizefoldericon/resizefoldericon.go

deps:
	@GOPATH=$(GOPATH) go get github.com/nfnt/resize
