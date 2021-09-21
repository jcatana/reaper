package config

import (
    "flag"
    "os"
    "path/filepath"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
)

type Config struct {
    kconf       *rest.Config
    clientset   kubernetes.Interface
    global      Global
}

type Global struct {
    killTime    string
    vendor      string
    loopSeconds string
    logLevel    string
}

func NewConfig() *Config {
    cfg := &Config {}
    cfg.Populate()
    return cfg
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func (c *Config) GetClientset() kubernetes.Interface {
    return c.clientset
}
func (c *Config) GetLogLevel() string {
    return c.global.logLevel
}
func (c *Config) GetVendor() string {
    return c.global.vendor
}
func (c *Config) GetKillTime() string {
    return c.global.killTime
}
func (c *Config) GetLoopSeconds() string {
    return c.global.loopSeconds
}

func (c *Config) Populate() {
    //var err error
    kconf, clientset, err := k8sConfig()
    if err == nil {
        c.kconf = kconf
        c.clientset = clientset
    }
    c.global.killTime = getEnv("killTime", "48h")
    c.global.vendor = getEnv("vendor", "reaper.io")
    c.global.loopSeconds = getEnv("loopSeconds", "10")
    c.global.logLevel = getEnv("logLevel", "trace")
}

func k8sConfig() (*rest.Config, kubernetes.Interface, error) {
    var config *rest.Config
    if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount"); err == nil {
        config, err = rest.InClusterConfig()
        if err != nil {
            panic(err.Error())
        }
    } else {
        var kubeconfig *string
        if home := homedir.HomeDir(); home != "" {
            kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
        } else {
            kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
        }
        flag.Parse()

        config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
        if err != nil {
            panic(err.Error())
        }
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    return config, clientset, nil
}
