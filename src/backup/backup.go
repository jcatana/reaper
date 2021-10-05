package backup

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jcatana/reaper/backup/targets"
	"github.com/jcatana/reaper/config"
	"github.com/sirupsen/logrus"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	//"k8s.io/apimachinery/pkg/runtime/serializer/json"
	//"k8s.io/apimachinery/pkg/runtime"
	"encoding/json"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
)

func init() {
}

func DoBackup(cfg *config.Config, log *logrus.Logger, gvk string) error {
	kconfig := cfg.GetKconf()
	clientset, err := kubernetes.NewForConfig(kconfig)
	if err != nil {
		fmt.Printf("\nerror\n%v\n\n", err)
	}
	data, err := clientset.RESTClient().Get().AbsPath(gvk).DoRaw(context.TODO())
	if err != nil {
		fmt.Printf("\nerror\n%v\n\n", err)
	}
	out := unstructured.Unstructured{}
	err = json.Unmarshal(data, &out)
	if err != nil {
		fmt.Printf("\nerror\n%v\n\n", err)
	}
	rObj := out.DeepCopyObject()

	var yamlBool bool
	switch format := cfg.GetBackupFormat(); format {
	case "json":
		yamlBool = false
	case "yaml":
		yamlBool = true
	}

	serializer := k8sjson.NewSerializerWithOptions(
		k8sjson.DefaultMetaFactory, nil, nil,
		k8sjson.SerializerOptions{
			Yaml:   yamlBool,
			Pretty: false,
			Strict: true,
		},
	)

	buffer := new(bytes.Buffer)
	serializer.Encode(rObj, buffer)
	for _, method := range cfg.GetEnabledTargets() {
		switch method {
		case "stdout":
			{
				err := stdout.Backup(buffer.String())
				if err != nil {
					log.WithFields(logrus.Fields{"selflink": gvk}).Error("Failed backup to %s", method)
				}
			}
		}
	}
	return nil
}
