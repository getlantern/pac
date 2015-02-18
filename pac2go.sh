#!/usr/bin/env bash

function die() {
  echo $*
  exit 1
}
which 2goarray > /dev/null || go get github.com/cratonica/2goarray/...
which 2goarray > /dev/null || die "Please install 2goarray manually, then try again"

if [ -r "pac-cmd/pac" ]
then
  2goarray pacBytes pac < pac-cmd/pac > pac_bytes_darwin.go
fi

if [ -r "pac-cmd/pac-linux" ]
then
  2goarray pacBytes pac < pac-cmd/pac-linux > pac_bytes_linux.go
fi

if [ -r "pac-cmd/pac" ]
then
  2goarray pacBytes pac < pac-cmd/pac.exe > pac_bytes_windows.go
fi
