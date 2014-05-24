#!/bin/bash

set -e

go build
version=$(./timeat --version)

filename=timeat-${version}-amd64.bz2
echo "building ${filename}"
bzip2 -c timeat > ${filename}

filename=timeat-${version}-386.bz2
echo "building ${filename}"
GOARCH=386 go build
bzip2 -c timeat > ${filename}
