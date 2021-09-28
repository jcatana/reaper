#!/bin/sh
#set -x

mkdir -p "tmp.dockerWorkArea"
if [ "$?" -ne 0 ]; then
    echo "ERROR: Unable to create 'tmp.dockerWorkArea' directory." >&2
fi

( cd src && CGO_ENABLED=0 go build -o ../tmp.dockerWorkArea/reaper )
if [ "$?" -ne 0 ]; then
    echo "ERROR: Unable to cd into 'tmp.dockerWorkArea' and build the binary." >&2
fi

which go
if [ "$?" -ne 0 ]; then
    echo "ERROR: Unable to find the go compiler." >&2
fi

( cd tmp.dockerWorkArea && cp -a ../Dockerfile ./ && docker build -t reaper:1 . )

kind load docker-image reaper:1

