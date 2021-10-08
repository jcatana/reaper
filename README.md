# Reaper

## Description

Reaper is a daemon that automatically deletes items in a Kubernetes namespace that has been running longer than a specified time.

## Caveats


- Reaper is intended to be used on Kubernetes DEV clusters where the developers aren't using autoscaler-scale-to-zero or don't clean up after themselves.
- Reaper will only remove resource consuming objects.
- Reaper does not remove things like configMaps, secrets, or services (currently).
- Reaper currently only supports native resource types and hasn't been tested on CRDs.

## How Reaper Works

- Reaper works through setting labels on namespaces.
- When a namespace is labelled with `reaper.io/enabled=True`, the reaper daemon will begin montoring the objects in that namespace to check if their creation timestamp, is passed the allocated time scale.
- Many configuration parameters can be overridden on a per namespace level.


<!--
## Docker Images

TODO
-->

## Dependency Setup

1. If you don't have the latest version of `go` installed, use https://golang.org/.
2. If you don't have `kubectl` installed, use https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/. (Do not install kubectl via snap, as it has some defects.)
3. If you don't have `kind` installed, use https://kind.sigs.k8s.io/docs/user/quick-start/.
4. If you don't have `helm` installed, use:

```shell
curl https://get.helm.sh/helm-v3.7.0-linux-amd64.tar.gz -O
tar -zxvf helm-v3.7.0-linux-amd64.tar.gz
rm -rf helm-v3.7.0-linux-amd64.tar.gz

#mv linux-amd64/helm /usr/local/bin/helm
mv linux-amd64/helm ~/bin/
#chmod +x /usr/local/bin/helm
chmod +x ~/bin/helm

rm -rf linux-amd64
```

5. Create a `kind` cluster:

```shell
# add ~/bin to your ~/.bashrc in the PATH environment variable
which kind || echo 'export PATH=$PATH:~/bin' >> ~/.bashrc && . ~/.bashrc

# create the new cluster in kind
kind create cluster
```

6. Setup the environment for using Kubernetes via cli:

```shell
echo "alias k=kubectl" >> ~/.bashrc
# TODO Find a better solution tot eh following line, as it's very invasive in 
# how it sprays thousands of lines into ~/.bashrc.
# echo "complete -F __start_kubectl k" >> ~/.bashrc
. ~/.bashrc

# install kubectl is you don't already have it installed, just not via snap, 
# or you may encounter issues with 
# "error: write /dev/stdout: permission denied"

# TODO Find a better solution tot eh following line, as it's very invasive in 
# how it sprays thousands of lines into ~/.bashrc.
#kubectl completion bash >> ~/.bashrc
```

## Installation

Deploy reaper.

```shell
git clone https://github.com/jcatana/reaper.git
cd reaper/helm/reaper
k create namespace reaper

helm --namespace reaper install reaper .
```

## Development Environment

Create a bunch of namespaces and label them to be monitored:

```shell
# TODO Figure out why the following doesn't work.
for i in `seq 1 10`; do
    k create namespace test${i}
    k label namespace test${i} reaper.io/enabled=True
    k annotate namespace reaper.io/killTime=${i}m
done
```

Create a bunch of sleep deployments to be killed:

```shell
# TODO Figure out why the following doesn't work.
for i in `seq 1 10`; do
    k -n create test-sleep-deployment.yaml
done
```

## TODO

The TODO list has been moved [HERE](https://github.com/jcatana/reaper/projects/1).

