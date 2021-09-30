# Reaper

This project is for larger shared dev clusters where users may not use autoscalers or clean up after themselves.

Reaper will monitor the targeted namespaces in your cluster and delete the resources that have been running for over a specified amount of time. It will only remove resource consuming objects and does not remove things like configMaps, secrets, or services (currently).


## Installation

This won't actually work right now without releasing the Docker image

```
git clone https://github.com/jcatana/reaper.git
cd reaper/helm/reaper
k create namespace reaper
helm -n reaper install . 
```


## Configuration

Reaper works through setting labels on namespace. When a namespace is labelled with `reaper.io/enabled=True`, the reaper daemon will begin montoring the objects in that namespace to check if their creation timestamp, is passed the allocated time scale.

Many configuration parameters can be overridden on a per namespace level.


## Docker images

?

## Devlopment environment
A couple things to start you out with cli kubernetes:
```
echo "alias k=kubectl" >> ~/.bashrc
echo "complete -F __start_kubectl k" >> ~/.bashrc
kubectl completion bash >> ~/.bashrc
```
Install helm to deploy the chart
```
curl https://get.helm.sh/helm-v3.7.0-linux-amd64.tar.gz -O
tar -zxvf helm-v3.7.0-linux-amd64.tar.gz
rm -rf helm-v3.7.0-linux-amd64.tar.gz
mv linux-amd64/helm /usr/local/bin/helm
chmod +x /usr/local/bin/helm
rm -rf linux-amd64
```
I use kind to stand up a dev instance of k8s and do my testing.
```
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.11.1/kind-linux-amd64
chmod +x ./kind
mv ./kind /usr/local/bin/kind
```
after that you create a kind cluster
```
kind create cluster
```
Create a bunch of namespaces and label them to be monitored
```
for i in `seq 1 10`; do k create namespace test${i}; k label namespace test${i} reaper.io/enabled=True; k annotate namespace reaper.io/killTime=${i}m; done
```
Create a bunch of sleep deployments to be killed
```
for i in `seq 1 10`; do k -n create test-sleep-deployment.yaml; done
```

## TODO
Todo list has been moved [here](https://github.com/jcatana/reaper/projects/1).
