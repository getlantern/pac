#!/usr/bin/env bash

function die() {
  echo $*
  exit 1
}
echo "check and install 2goarray..."
which 2goarray > /dev/null || go get github.com/cratonica/2goarray
which 2goarray > /dev/null || die "install failed,  please install 2goarray manually"
gcc -x objective-c -framework Cocoa -framework SystemConfiguration -framework Security helper.m -o helper || die
2goarray helper pacon < helper > helper.go
