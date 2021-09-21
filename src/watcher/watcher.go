package watcher

import (
    "time"
    "context"
    "github.com/jcatana/reaper/config"
    "github.com/sirupsen/logrus"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/informers"
    
)

type Watch map[string][]WatchResource

type WatchResource struct {
    resourceName string `json:"resourceName"`
    creationTimestamp string `json:"creationTimestamp"`
    ownkind string `json:"ownkind"`
    killTime time.Duration `json:"killTime"`
}

func NewWatcher() Watch {
    reap := make(Watch)
    //watch := &Watch {}
    //reap := make(map[string][]WatchResources)
    return reap
}
func (w Watch) GetResource(namespace string, idx int) WatchResource {
    return w[namespace][idx]
}
func (w Watch) GetResourceName(namespace string, idx int) string {
    return w[namespace][idx].resourceName
}
func (w Watch) GetCreationTimestamp(namespace string, idx int) string {
    return w[namespace][idx].creationTimestamp
}
func (w Watch) GetOwnkind(namespace string, idx int) string {
    return w[namespace][idx].ownkind
}
func (w Watch) GetKillTime(namespace string, idx int) time.Duration {
    return w[namespace][idx].killTime
}
/*
func (r Watch) removeObject(namespace string, log *logrus.Logger, idx int) []WatchResource {
    rw := r[namespace]
    if idx == 0 {
        rw = rw[idx+1:]
    }
    if idx == len(rw) {
        rw = rw[:idx]
    }
    if (idx > len(rw) && idx > 0) {
        rw = append(rw[:idx], rw[idx+1:]...)
    }
    log.WithFields(logrus.Fields{"idx": idx}).Trace("Removing id from slice")
    return rw
}
*/

func findParent(mObj metav1.Object, log *logrus.Logger, clientset kubernetes.Interface, ownkind string) (metav1.Object, string) {
    if len(mObj.GetOwnerReferences()) == 0 {
        //fmt.Printf("Namespace: %s Name: %s Timestamp: %s", mObj.GetNamespace(), mObj.GetName(), mObj.GetCreationTimestamp())
        log.WithFields(logrus.Fields{"namespace": mObj.GetNamespace(), "name": mObj.GetName(), "creationTimestamp": mObj.GetCreationTimestamp()}).Debug("Checking for parent")
        return mObj, ownkind
    } else {
        //fmt.Printf("Namespace: %s Name: %s Timestamp: %s", mObj.GetNamespace(), mObj.GetName(), mObj.GetCreationTimestamp())
        log.WithFields(logrus.Fields{"namespace": mObj.GetNamespace(), "name": mObj.GetName(), "creationTimestamp": mObj.GetCreationTimestamp()}).Debug("Checking for parent")
        owner := mObj.GetOwnerReferences()
        //fmt.Printf("%s", owner[0].Kind)
        log.WithFields(logrus.Fields{"owner": owner[0].Kind}).Debug("Owner")
        var child metav1.Object
        var err interface{}
        switch owner[0].Kind {
        case "ReplicaSet":
            ownkind = owner[0].Kind
            child, err = clientset.AppsV1().ReplicaSets(mObj.GetNamespace()).Get(context.TODO(), owner[0].Name, metav1.GetOptions{})
            if err != nil {
                //fmt.Printf("Error")
                log.Error("Error", err)
            }   
        case "DaemonSet":
            ownkind = owner[0].Kind
            child, err = clientset.AppsV1().DaemonSets(mObj.GetNamespace()).Get(context.TODO(), owner[0].Name, metav1.GetOptions{})
            if err != nil {
                //fmt.Printf("Error")
                log.Error("Error", err)
            }   
        case "StatefulSet":
            ownkind = owner[0].Kind
            child, err = clientset.AppsV1().StatefulSets(mObj.GetNamespace()).Get(context.TODO(), owner[0].Name, metav1.GetOptions{})
            if err != nil {
                //fmt.Printf("Error")
                log.Error("Error", err)
            }   
        case "Deployment":
            ownkind = owner[0].Kind
            child, err = clientset.AppsV1().Deployments(mObj.GetNamespace()).Get(context.TODO(), owner[0].Name, metav1.GetOptions{})
            if err != nil {
                //fmt.Printf("Error")
                log.Error("Error", err)
            }   
        default:
            ownkind = owner[0].Kind
            //fmt.Printf("Dunno, might be CRD, defaulting")
            log.WithFields(logrus.Fields{"namespace": mObj.GetNamespace, "ownkind": owner[0].Name, "resource": mObj.GetName()}).Debug("This might be CRD cannot identify")
        }
        mObj, ownkind := findParent(child, log, clientset, ownkind)
        return mObj, ownkind
    }
    return mObj, ownkind
}

func StartWatching(stopper <-chan struct{}, s cache.SharedIndexInformer, log *logrus.Logger, cfg *config.Config, reap Watch) {
    clientset := cfg.GetClientset()
    handlers := cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            mObj := obj.(metav1.Object)
            //This means the object is a namespace, or "not a namespaced object" like an admissions controller or cluster-rbac
            if len(mObj.GetNamespace()) > 0 {
                //retrieve namespace annotations to look for overrides
                ns, err := clientset.CoreV1().Namespaces().Get(context.TODO(), mObj.GetNamespace(), metav1.GetOptions{})
                if err != nil {
                    panic("Oh shit")
                }
                //set global kill time then see if there is an override on the namespace annotation
                killTime, _ := time.ParseDuration(cfg.GetKillTime())
                if key, ok := ns.Annotations[cfg.GetVendor()+"/killTime"]; ok {
                    killTime, _ = time.ParseDuration(key)
                }

                //search for parent
                pObj, ownkind := findParent(mObj, log, clientset, "Pod")

                //add returned parent to the store
                reap[pObj.GetNamespace()] = append(
                    reap[pObj.GetNamespace()],
                    WatchResource{
                        resourceName: pObj.GetName(),
                        creationTimestamp: pObj.GetCreationTimestamp().String(),
                        ownkind: ownkind,
                        killTime: killTime,
                    },
                )
                log.WithFields(logrus.Fields{"namespace": pObj.GetNamespace(), "kind": ownkind, "name": pObj.GetName()}).Info("Adding Object to store")
            } else {
                log.WithFields(logrus.Fields{"namespace": mObj.GetName()}).Info("Watching namespace")
            }
            //Start watchers for the objects in the namespace
            factory := informers.NewFilteredSharedInformerFactory(clientset, 0, mObj.GetName(), nil)
            informer := factory.Core().V1().Pods()
            go StartWatching(stopper, informer.Informer(), log, cfg, reap)

        },
        //Is this needed? Will probably need this to update killTime.
        UpdateFunc: func(oldObj, obj interface{}) {
            mObj := obj.(metav1.Object)
            //This means the object is a namespace, or "not a namespaced object" like an admissions controller or cluster-rbac
            /*
            if len(mObj.GetNamespace()) > 0 {
            } else {
            }
            */
            //fmt.Printf("Updated object in store. Namespace: %s Name: %s Timestamp: %s", mObj.GetName(), mObj.GetCreationTimestamp())
            log.WithFields(logrus.Fields{"namespace": mObj.GetName(), "name": mObj.GetCreationTimestamp()}).Info("Updated object in store")
        },
        DeleteFunc: func(obj interface{}) {
            mObj := obj.(metav1.Object)
            //If the namespace label is removed, remove the key to remove all the objects from being watched.
            delete(reap, mObj.GetName())
            //fmt.Printf("Deleted object from store. Namespace: %s Name: %s Timestamp: %s", mObj.GetName(), mObj.GetCreationTimestamp())
            log.WithFields(logrus.Fields{"namespace": mObj.GetName(), "name": mObj.GetCreationTimestamp()}).Info("Deleted object from store")
        },
    }
    s.AddEventHandler(handlers)
    s.Run(stopper)
}

