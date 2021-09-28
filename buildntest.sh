#!/bin/sh
#set -x # uncomment this to have the shell output the commands before they are ran

if ! mkdir -p "tmp.dockerWorkArea"; then
    echo "ERROR: Unable to create 'tmp.dockerWorkArea' directory." >&2
    exit 1
fi

baseDir=$(pwd)

if ! cd src; then
    echo "ERROR: Unable to cd into 'tmp.dockerWorkArea' and build the binary." >&2
    exit 1
fi

if ! CGO_ENABLED=0 go build -v -o ../tmp.dockerWorkArea/reaper; then
    echo "ERROR: Unable to cd into 'tmp.dockerWorkArea' and build the binary." >&2
    exit 1
fi

if ! which go; then
    echo "ERROR: Unable to find the go compiler." >&2
    exit 1
fi

cd $baseDir

if ! cp -a Dockerfile tmp.dockerWorkArea/; then
    echo "ERROR: Unable to copy the Dockerfile." >&2
    exit 1
fi

if ! cd tmp.dockerWorkArea; then
    echo "ERROR: Unable to build the docker image." >&2
    exit 1
fi

if ! docker build -f Dockerfile -t reaper:1 .; then
    echo "ERROR: Unable to build the docker image." >&2
    exit 1
fi

if ! docker image save reaper:1 -o reaper-1.dockerImage.tar; then
    echo "ERROR: Unable save out the docker image." >&2
    exit 1
fi

cd $baseDir

# microk8s is micro kubernetes: https://microk8s.io/
if which microk8s.kubectl; then
    if ! microk8s.ctr image import ./tmp.dockerWorkArea/reaper-1.dockerImage.tar; then
        echo "ERROR: Unable to load the docker image into microk8s." >&2
        exit 1
    fi

    if ! microk8s.ctr images ls|grep "SIZE\|reaper"; then
        echo "ERROR: Unable to load the docker image into microk8s." >&2
        exit 1
    fi

    echo "INFO: The image for reaper was successfully imported into microk8s to create containers from."

    echo "INFO: Now attempting to start reaper in microk8s."

    # start/replace reaper in microk8s
    microk8s kubectl replace --force -f manifest.yaml

    # stop reaper in microk8s
    #microk8s kubectl delete -f manifest.yaml

elif which kind; then
    # kind is kubernetes in docker: https://kind.sigs.k8s.io/
    if ! kind load docker-image reaper:1; then
        echo "ERROR: Unable to load the docker image into kind." >&2
        exit 1
    fi
    
    echo "INFO: The image for reaper was successfully imported into kind to create containers from."
else
    echo "ERROR: Unable to find microk8s or kind to load the image into." >&2
    exit 1
fi

echo "INFO: Deleting tmp.dockerWorkArea/ because we care."
rm -rf tmp.dockerWorkArea/

exit 0

