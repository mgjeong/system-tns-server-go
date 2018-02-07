#!/bin/bash

export GOPATH=$PWD

rm -rf $GOPATH/src/github.com/docker/distribution/vendor/github.com/opencontainers

pkg_list=()

count=0
for pkg in "${pkg_list[@]}"; do
 if [ $? -ne 0 ]; then
    echo "Unittest is failed."
    func_cleanup
    exit 1
 fi
 count=$count.0
done


