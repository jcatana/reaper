package watcher

import (
	"context"
	"github.com/jcatana/reaper/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

// Watch is a table of: [namespace][resource] values.. where the resource is the v1/app deployment, stateful set, etc
type Watch map[string]map[string]WatchResource

type WatchResource struct {
	creationTimestamp string
	ownkind           string
	killTime          time.Duration
	gvkPath           string
}

func NewWatcher() Watch {
	reap := make(Watch)
	return reap
}

func (w WatchResource) GetCreationTimestamp() string {
	return w.creationTimestamp
}
func (w WatchResource) GetGvkPath() string {
	return w.gvkPath
}
func (w WatchResource) GetOwnkind() string {
	return w.ownkind
}
func (w WatchResource) GetKillTime() time.Duration {
	return w.killTime
}

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
					panic("Oh shit") // Need to remove this with an actual error
				}
				//set global kill time then see if there is an override on the namespace annotation
				killTime, _ := time.ParseDuration(cfg.GetKillTime())
				if key, ok := ns.Annotations[cfg.GetVendor()+"/killTime"]; ok {
					killTime, _ = time.ParseDuration(key)
				}

				//search for parent
				pObj, ownkind := findParent(mObj, log, clientset, "Pod")

				//add returned parent to the store
				reap[pObj.GetNamespace()][pObj.GetName()] = WatchResource{
					creationTimestamp: pObj.GetCreationTimestamp().String(),
					ownkind:           ownkind,
					killTime:          killTime,
					gvkPath:           pObj.GetSelfLink(),
				}
				log.WithFields(logrus.Fields{"namespace": pObj.GetNamespace(), "kind": ownkind, "name": pObj.GetName()}).Info("Adding Object to store")
			} else {
				reap[mObj.GetName()] = make(map[string]WatchResource)
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
