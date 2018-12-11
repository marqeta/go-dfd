#!/usr/bin/env bash

curl -L -s https://github.com/golang/dep/releases/download/v$DEP_VERSION/dep-linux-amd64 -o $GOPATH/bin/dep
chmod +x $GOPATH/bin/dep
PACKAGE_PATH=$GOPATH/src/go-dfd
ln -s $DRONE_WORKSPACE $PACKAGE_PATH
cd $PACKAGE_PATH
dep ensure
