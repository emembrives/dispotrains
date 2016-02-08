#!/bin/sh
BASE_DIR=$PWD
mkdir gopath
export GOPATH=$PWD/gopath
mkdir -p $GOPATH/github.com/emembrives/dispotrains
ln -s $PWD/dispotrains.webapp $GOPATH/github.com/emembrives/dispotrains/dispotrains.webapp

cd $GOPATH/github.com/emembrives/dispotrains/dispotrains.webapp
go get ./...
go build -v ./...
cd $BASE_DIR
