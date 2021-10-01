package backup

import (
	"context"
	"fmt"
	"os"
	"github.com/jcatana/reaper/config"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"encoding/json"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
)

func DoBackup(cfg *config.Config, gvk string) error {
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
		Pretty: true,
		Strict: true,
		},
	)
	serializer.Encode(rObj, os.Stdout)

	return nil
}
