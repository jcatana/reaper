for i in `seq 1 9`; do 
	kubectl get namespace reap-test${i} > /dev/null 2>&1
	if [ $? == 1 ]; then
		kubectl create namespace reap-test${i}
		kubectl label namespace reap-test${i} reaper.io/enabled=True
		kubectl annotate namespace reap-test${i} reaper.io/killTime=${i}m
	fi
	kubectl -n reap-test${i} create -f test-sleep-deployment.yaml
done
