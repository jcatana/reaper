docker build -f docker/Dockerfile -t reaper:1 .
kind load docker-image reaper:1
kubectl -n reaper rollout restart deployment reaper
