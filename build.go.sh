#!/bin/sh
mkdir gopath
export GOPATH=$PWD/gopath
mkdir -p $GOPATH/github.com/emembrives/dispotrains
ln -s $PWD/dispotrains.webapp $GOPATH/github.com/emembrives/dispotrains/dispotrains.webapp
pushd $GOPATH/github.com/emembrives/dispotrains/dispotrains.webapp
go get ./...
go build -v ./...
