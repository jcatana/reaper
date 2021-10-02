package reaper

import (
	"context"
	"github.com/jcatana/reaper/backup"
	"github.com/jcatana/reaper/config"
	"github.com/jcatana/reaper/watcher"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"time"
)

//constants
const refTime = "2006-01-02 15:04:05 -0700 MST"

func reapObject(log *logrus.Logger, clientset kubernetes.Interface, namespace string, ownkind string, resource string) error {
	deletePolicy := metav1.DeletePropagationBackground
	var err error
	switch ownkind {
	case "ReplicaSet":
		if err = clientset.AppsV1().ReplicaSets(namespace).Delete(context.TODO(), resource, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
			log.Error("Error", err)
		}
	case "DaemonSet":
		if err = clientset.AppsV1().DaemonSets(namespace).Delete(context.TODO(), resource, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
			log.Error("Error", err)
		}
	case "StatefulSet":
		if err = clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), resource, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
			log.Error("Error", err)
		}
	case "Deployment":
		if err = clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), resource, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
			log.Error("Error", err)
		}
	default:
		log.WithFields(logrus.Fields{"namespace": namespace, "ownkind": ownkind, "resource": resource}).Debug("This might be CRD cannot identify")
	}
	if err == nil {
		return nil
	} else {
		return err
	}
}

func Reap(stopper <-chan struct{}, log *logrus.Logger, cfg *config.Config, reap watcher.Watch) {
	for {
		log.Trace(reap)
		for namespace, resources := range reap {
			for resource, params := range resources {
				//figure out time
				currentTime := time.Now().Truncate(time.Second)
				creationTime, _ := time.Parse(refTime, params.GetCreationTimestamp())
				killTime := creationTime.Add(params.GetKillTime())
				log.WithFields(logrus.Fields{"namespace": namespace, "resource": resource, "KTduration": params.GetKillTime()}).Trace("Checking")
				log.WithFields(logrus.Fields{"creationTimestamp": params.GetCreationTimestamp(), "killTime": params.GetKillTime(), "currentTime": currentTime}).Trace("Times")

				if currentTime.After(killTime) {
					if len(cfg.GetBackup()) > 0 {
						err := backup.DoBackup(cfg, params.GetGvkPath())
						if err != nil {
							log.WithFields(logrus.Fields{"namespace": namespace, "kind": params.GetOwnkind(), "resource": resource}).Info("Cannot backup")
							break
						}
					}

					err := reapObject(log, cfg.GetClientset(), namespace, params.GetOwnkind(), resource)
					if err == nil {
						delete(reap[namespace], resource)
						break
					}
				}
			}
		}
		loopSeconds, _ := strconv.Atoi(cfg.GetLoopSeconds())
		time.Sleep(time.Duration(loopSeconds) * time.Second)
	}
}
