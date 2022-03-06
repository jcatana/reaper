# Reaper

## Description

Reaper is a daemon that automatically deletes items in a Kubernetes namespace that has been running longer than a specified time.

## How Reaper Works

- Reaper works through setting labels on namespaces.
- When a namespace is labeled with `reaper.io/enabled=True`, the reaper daemon will begin monitoring the objects in that namespace to check if their creation timestamp, is passed the allocated time interval.
- Many configuration parameters can be overridden on a per namespace level.

## Caveats

- Currently requires featureGate RemoveSelfLink to be set to false as it relies on selflink to derive the object yaml for backups.
- Reaper is intended to be used on Kubernetes DEV clusters where the developers aren't using autoscaler-scale-to-zero or don't clean up after themselves.
- Reaper will only remove resource consuming objects.
- Reaper does not remove things like configMaps, secrets, or services (currently).
- Reaper currently only supports native resource types and hasn't been tested on CRDs.

## Development Environment Setup

### Dependency Setup

NOTE: Make sure you're not using a `snap` installed version of the below dependencies, or you'll likely run into permissions or missing file path issues.

1. If you don't have the latest, or second latest, minor version of `go` installed, use https://golang.org/.
2. If you don't have `kubectl` installed, use https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/.
3. Setup the environment for using Kubernetes via cli using an alias of `k`:

```shell
# add 'k' as an alias for 'kubectl'
alias k || echo "alias k=kubectl" >> ~/.bashrc && . ~/.bashrc
 
# prepare autocomplete for kubectl commands
kubectl completion bash > ~/.bashrc.kubectl
grep '^\. ~/\.bashrc\.kubectl$' ~/.bashrc || echo '. ~/.bashrc.kubectl' >> ~/.bashrc && . ~/.bashrc

# enable autocomplete for 'k' the alias for 'kubectl', if not already setup
grep "^complete -F __start_kubectl k$" ~/.bashrc || echo "complete -F __start_kubectl k" >> ~/.bashrc && . ~/.bashrc
```
4. If you don't have `kind` installed, use https://kind.sigs.k8s.io/docs/user/quick-start/.

NOTE: If you installed `kind` v0.10.x using the default kubernetes ('kindest/node') v1.20.x image, or a version of kind `kind` where you explicitly installed the kubernetes ('kindest/mode') v1.20.x image, then you'll need to set RemoveSelfLink to false. This workaround will not work with the kubernetes ('kindest/node') v1.21.x or later images. To set RemoveSelfLink to false, run the following:
   
```shell
cat << EOF > kind.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
featureGates:
  "RemoveSelfLink": false
EOF
```

5. If you don't have `helm` installed, use:

```shell
curl https://get.helm.sh/helm-v3.7.0-linux-amd64.tar.gz -O
tar -zxvf helm-v3.7.0-linux-amd64.tar.gz
rm -rf helm-v3.7.0-linux-amd64.tar.gz

mkdir -p ~/bin
#mv linux-amd64/helm /usr/local/bin/helm
mv linux-amd64/helm ~/bin/
#chmod +x /usr/local/bin/helm
chmod +x ~/bin/helm

rm -rf linux-amd64
```

6. Create a `kind` cluster:

```shell
# add ~/bin to your ~/.bashrc in the PATH environment variable
which kind || echo 'export PATH=$PATH:~/bin' >> ~/.bashrc && . ~/.bashrc

# create the new cluster in kind
if [ -f kind.yaml ]; then
    kind create cluster --config kind.yaml
else
    kind create cluster
fi
```

7. Clone Reaper & Create Reaper Namespace in Kubernetes

```shell
git clone https://github.com/jcatana/reaper.git
cd reaper
baseDir=$(pwd)

k create namespace reaper
```

8. Create Reaper Namespace

```shell
kubectl create namespace reaper
cd helm/reaper 
helm --namespace reaper install reaper .
cd $baseDir
```
9. Compile Reaper, Build Reaper Image, & Load in kubernetes

```shell
./build.sh
```

10. Create Testing Namespaces to Monitor

You can run `./test.sh` it will do the same thing as these two commands

```shell
for i in `seq 1 9`; do
    k create namespace test${i}
    k label namespace test${i} reaper.io/enabled=True
    k annotate namespace test${i} reaper.io/killTime=${i}m
done
```

```shell
# cd to the repository root where test-sleep-deployment.yaml is, and then run this:
for i in `seq 1 9`; do
    k --namespace="test${i}" create --filename="test-sleep-deployment.yaml"
done
```


12. Watch the Reaper and Test Namespaces
```shell
watch -n 1 'kubectl get pods -A'
```

13. What to watch for?
- ?
