#!/usr/bin/env bash

function die() {
  echo $*
  exit 1
}

function gen() {
  bin=$1
  dst=$2
  if [ -f "$bin" ]
  then
    echo Generating $dst from $bin...
    2goarray pacBytes pac < $bin > $dst
  else
    echo $bin does not exist, skipping...
  fi
}

which 2goarray > /dev/null || go get github.com/cratonica/2goarray/...
which 2goarray > /dev/null || die "Please install 2goarray manually, then try again"
cmd_path=$GOPATH/src/github.com/getlantern/pac-cmd
gen $cmd_path/pac pac_bytes_darwin.go
gen $cmd_path/pac-linux pac_bytes_linux.go
gen $cmd_path/pac.exe pac_bytes_windows.go
