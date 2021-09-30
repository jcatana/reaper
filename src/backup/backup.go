package backup

import (
    "context"
    "fmt"
    //"os"
    "github.com/jcatana/reaper/config"
    //metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    //"k8s.io/apimachinery/pkg/runtime/serializer/json"
    "encoding/json"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/client-go/kubernetes"
    //"k8s.io/client-go/rest"
)

//func DoBackup(cfg *config.Config, mLst metav1.List) {
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
    oJson := unstructured.Unstructured{}
    err = json.Unmarshal(data, &oJson)
    if err != nil {
        fmt.Printf("\nerror\n%v\n\n", err)
    }
    fmt.Printf("\njson\n%v\n\n", oJson)
    /*oLst, err := clientset.AppsV1().Deployments(pObj.GetNamespace()).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("\nerror\n%v\n\n", err)
    }

    fmt.Printf("\n\ncfg\n%v\n", cfg)
    rObj := oLst.DeepCopyObject()
    fmt.Printf("\n\nrobj\n%v\n", rObj)
    s := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, 
        json.SerializerOptions{
            Yaml: true,
            Pretty: true,
            Strict: true,
        },
    )
    err = s.Encode(rObj, os.Stdout)
    */
    return nil
}
