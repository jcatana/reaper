#!/bin/sh
#set -x # uncomment this to have the shell output the commands before they are ran
set +e # do not exit on error, but permit script to control when to exit

if [ "$(echo "$FOSSA_API_KEY" | grep . | wc -l)" -ne 1 ]; then
    echo "ERROR: Please set your FOSSA_API_KEY and try again." >&2
    exit 1
fi

baseDir=$(pwd)

if [ -d "tmp.dockerWorkArea" ]; then
    echo "INFO: Deleting tmp.dockerWorkArea/ because we care."
    echo
    rm -rf tmp.dockerWorkArea/
fi

if ! cd src; then
    echo "ERROR: Unable to cd into 'src'." >&2
    exit 1
fi

goPath=$(which go)
goPathCount=$(echo "$goPath" | grep . | wc -l)
if [ "$goPathCount" -eq 1 ]; then
    echo "# Found 'go' at '$goPath'"
    echo
else
    echo "ERROR: Unable to find the 'go' compiler." >&2
    exit 1
fi

cd $baseDir

echo "Currently In: '$(pwd)'"
echo 

if ! docker build -f Dockerfile.scan -t reaper:dev .; then
    echo "ERROR: Unable to build the docker image." >&2
    exit 1
fi

#if [ "$DEBUG" == "1" ]; then
#    docker run -it -e FOSSA_API_KEY reaper:dev sh
#else
    docker run --rm -e FOSSA_API_KEY reaper:dev ./.scan.sh
#fi

