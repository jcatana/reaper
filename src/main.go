package main

//imports
import (
    "os"
    "os/signal"
    "github.com/jcatana/reaper/config"
    "github.com/jcatana/reaper/logging"
    "github.com/jcatana/reaper/watcher"
    "github.com/jcatana/reaper/reaper"
    "github.com/sirupsen/logrus"

    "k8s.io/client-go/informers"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

//global variables
var cfg *config.Config
var log *logrus.Logger
//var reap watcher.Watch

//init: get the configuration
func init() {
    cfg = config.NewConfig()
    log = logging.NewLogger(cfg)
}

//main
func main () {
    reap := watcher.NewWatcher()

    factory := informers.NewFilteredSharedInformerFactory(cfg.GetClientset(), 0, corev1.NamespaceAll, func(options *metav1.ListOptions) {options.LabelSelector = cfg.GetVendor() + "/enabled=True"})
    informer := factory.Core().V1().Namespaces()

    stopper := make(chan struct{})

    //launch go routines
    go watcher.StartWatching(stopper, informer.Informer(), log, cfg, reap)
    go reaper.Reap(stopper, log, cfg, reap)

    signalChannel := make(chan os.Signal, 0)
    signal.Notify(signalChannel, os.Kill, os.Interrupt)

    <-signalChannel
    close(stopper)
}
