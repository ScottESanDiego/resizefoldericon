GOPATH=`pwd`

all: build

build:
	go build -o resizefoldericon src/github.com/ScottESanDiego/resizefoldericon/resizefoldericon.go

deps:
	go get github.com/nfnt/resize
