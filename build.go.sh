#!/bin/sh
mkdir gopath
export GOPATH=$PWD/gopath
mkdir $GOPATH/github.com/emembrives/dispotrains
ln -s $PWD/dispotrains.webapp $GOPATH/github.com/emembrives/dispotrains/dispotrains.webapp

go get ./dispotrains.webapp/...
go build -v ./dispotrains.webapp/...
