# Reaper

This project is for larger shared dev clusters where users may not use autoscalers or clean up after themselves.

Reaper will monitor the targeted namespaces in your cluster and delete the resources that have been running for over a specified amount of time. It will only remove resource consuming objects and does not remove things like configMaps, secrets, or services (currently).

## Instalation
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


## TODO
- documentation
- fix the build scripts and helm releases
- metrics - resources freed per iteration?
- add single namespace only operation
- ignore self in single namespace mode
- metrics. # of kills
- Make a global config helm configMap
- Backup data to PVC (or maybe config-maps?) before deleting
- SMTP support
- Heirarchical killtime overrides system. xGlobal -> xNamespace -> ?Entity
- Helm chart
- logo
- license

