package main

import (
	"github.com/jcatana/reaper/config"
	"github.com/jcatana/reaper/logging"
	"github.com/jcatana/reaper/reaper"
	"github.com/jcatana/reaper/watcher"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
)

//global variables
var cfg *config.Config
var log *logrus.Logger

func init() {
	cfg = config.NewConfig()
	log = logging.NewLogger(cfg)
}

func main() {
	reap := watcher.NewWatcher()

	factory := informers.NewFilteredSharedInformerFactory(cfg.GetClientset(), 0, corev1.NamespaceAll, func(options *metav1.ListOptions) {
		options.LabelSelector = cfg.GetVendor() + "/enabled=True"
	})
	informer := factory.Core().V1().Namespaces()

	stopper := make(chan struct{})

	// launch watcher go routine
	go watcher.StartWatching(stopper, informer.Informer(), log, cfg, reap)
	// launch reaper go routine
	go reaper.Reap(stopper, log, cfg, reap)

	signalChannel := make(chan os.Signal, 0)
	signal.Notify(signalChannel, os.Kill, os.Interrupt)

	<-signalChannel
	close(stopper)
}
