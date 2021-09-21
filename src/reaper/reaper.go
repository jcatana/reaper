package reaper

import (
    "context"
    "fmt"
    "strconv"
    "time"
    "github.com/jcatana/reaper/config"
    "github.com/jcatana/reaper/watcher"
    "github.com/sirupsen/logrus"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

//constants
const refTime = "2006-01-02 15:04:05 -0700 MST"

func reapObject(log *logrus.Logger, clientset kubernetes.Interface, namespace string, ownkind string, resource string) error {
    deletePolicy := metav1.DeletePropagationBackground
    var err error
    switch ownkind {
    case "ReplicaSet":
        if err = clientset.AppsV1().ReplicaSets(namespace).Delete(context.TODO(), resource, metav1.DeleteOptions{PropagationPolicy: &deletePolicy,}); err != nil {
            //fmt.Printf("Error %v\n", err)
            log.Error("Error", err)
        }
    case "DaemonSet":
        if err = clientset.AppsV1().DaemonSets(namespace).Delete(context.TODO(), resource, metav1.DeleteOptions{PropagationPolicy: &deletePolicy,}); err != nil {
            //fmt.Printf("Error %v\n", err)
            log.Error("Error", err)
        }
    case "StatefulSet":
        if err = clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), resource, metav1.DeleteOptions{PropagationPolicy: &deletePolicy,}); err != nil {
            //fmt.Printf("Error %v\n", err)
            log.Error("Error", err)
        }
    case "Deployment":
        if err = clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), resource, metav1.DeleteOptions{PropagationPolicy: &deletePolicy,}); err != nil {
            //fmt.Printf("Error %v\n", err)
            log.Error("Error", err)
        }
    default:
        //fmt.Printf("Dunno, might be CRD, defaulting")
        log.WithFields(logrus.Fields{"namespace": namespace, "ownkind": ownkind, "resource": resource}).Debug("This might be CRD cannot identify")
    }
    if err == nil {
        return nil
    } else {
        return err
    }
}

func removeObject(r []watcher.WatchResource, log *logrus.Logger, idx int) []watcher.WatchResource {
    fmt.Printf("%v\n", r)
    if idx == 0 {
        r = r[idx+1:]
    }
    if idx == len(r) {
        r = r[:idx]
    }
    if (idx > len(r) && idx > 0) {
        r = append(r[:idx], r[idx+1:]...)
    }
    log.WithFields(logrus.Fields{"idx": idx}).Trace("Removing id from slice")
    return r
}
/*
func (r []watcher.WatchResource) removeObject(log *logrus.Logger, idx int) []watcher.WatchResource {
    if idx == 0 {
        r = r[idx+1:]
    }
    if idx == len(r) {
        r = r[:idx]
    }
    if (idx > len(r) && idx > 0) {
        r = append(r[:idx], r[idx+1:]...)
    }
    log.WithFields(logrus.Fields{"idx": idx}).Trace("Removing id from slice")
    return r
}
*/

func Reap(stopper <-chan struct{}, log *logrus.Logger, cfg *config.Config, reap watcher.Watch) {
    for {
        log.Trace(reap)
        for namespace, object := range reap{
            for idx, _ := range object{
                //define current resource set
                resName := reap.GetResourceName(namespace, idx)
                resCreationTime := reap.GetCreationTimestamp(namespace, idx)
                resKillTime := reap.GetKillTime(namespace, idx)
                resOwnkind := reap.GetOwnkind(namespace, idx)
                //figure out time
                currentTime := time.Now().Truncate(time.Second)
                creationTime, _ := time.Parse(refTime, resCreationTime)
                killTime := creationTime.Add(resKillTime)
                log.WithFields(logrus.Fields{"namespace": namespace, "idx": idx, "resource": resName, "KTduration": resKillTime}).Trace("Checking")
                log.WithFields(logrus.Fields{"creationTimestamp": resCreationTime, "killTime": killTime, "currentTime": currentTime}).Trace("Times")

                if currentTime.After(killTime) {
                    err := reapObject(log, cfg.GetClientset(), namespace, resOwnkind, resName)
                    if err == nil {
                        reap[namespace] = removeObject(reap[namespace], log, idx)
                        break
                        //reap[namespace] = reap[namespace].removeObject(log, idx)
                    }
                }
            }
        }
        loopSeconds, _ := strconv.Atoi(cfg.GetLoopSeconds())
        time.Sleep(time.Duration(loopSeconds) * time.Second)
    }
}
