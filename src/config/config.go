package config

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

// This must be updated when a new valid target is added
func realTargets() []string {
    return []string{"smtp", "git", "file", "stdout"}
}

type Config struct {
	kconf     *rest.Config
	clientset kubernetes.Interface
	global    Global
}

type Global struct {
	killTime     string
	vendor       string
	loopSeconds  string
	logLevel     string
	backup       []string
	backupFormat string
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.Populate()
	return cfg
}

func (c *Config) GetKconf() *rest.Config {
	return c.kconf
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
func (c *Config) GetBackup() []string {
	return c.global.backup
}
func (c *Config) GetBackupFormat() string {
	return c.global.backupFormat
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func cleanString(targets []string) []string {
    for i, v := range targets {
        targets[i] = strings.ToLower(strings.TrimSpace(v))
    }
    return targets
}
func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
func checkTargets(targets []string) []string {
    var newSlice []string
    for _, v := range targets {
        for _, c := range realTargets() {
            if v == c {
                if contains(newSlice, v) {
                    break
                } else {
                    newSlice = append(newSlice, c)
                }
            }
        }
    }
    return newSlice
}
func ParseSlice(target string) []string {
    targets := strings.Split(target, ",")
    targets = cleanString(targets)
    targets = checkTargets(targets)
    return targets
}


func (c *Config) Populate() {
	//var err error
	kconf, clientset, err := k8sConfig()
	if err == nil {
		c.kconf = kconf
		c.clientset = clientset
	}
	c.global.killTime = getEnv("KILL_TIME", "48h")
	c.global.vendor = getEnv("VENDOR", "reaper.io")
	c.global.loopSeconds = getEnv("LOOP_SECONDS", "10")
	c.global.logLevel = getEnv("LOG_LEVEL", "trace")
	c.global.backup = ParseSlice(getEnv("BACKUP", ""))
	c.global.backupFormat = getEnv("BACKUP_FORMAT", "yaml")
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
