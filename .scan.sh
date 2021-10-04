#!/bin/sh
#set -x # uncomment this to have the shell output the commands before they are ran
#set +e # do not exit on error, but permit script to control when to exit

# This file gets ran inside of the docker container. It is required for proper container operation

echo
echo "# INFO: Attempting to build the binary"
echo

cd src && CGO_ENABLED=0 go build -v -o goreaper.gobin && curl -H 'Cache-Control: no-cache' https://raw.githubusercontent.com/fossas/fossa-cli/master/install.sh | sh && fossa init && fossa analyze


