package backup

import (
	"fmt"
	"github.com/jcatana/reaper/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

func DoBackup(cfg *config.Config, pObj metav1.Object) {
	var obj runtime.Object
	var scope conversion.Scope // While not actually used within the function, need to pass in
	err := runtime.Convert_runtime_RawExtension_To_runtime_Object(&pObj, &obj, scope)
	if err != nil {
		fmt.Printf("\nerror\n%v\n\n", err)
		//return nil, err
	}

	innerObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		fmt.Printf("\nerror\n%v\n\n", err)
		//return nil, err
	}
	u := unstructured.Unstructured{Object: innerObj}
	labels := u.GetLabels()
	kind := u.GetKind()
	fmt.Printf("lables: %v\n\n", labels)
	fmt.Printf("kind: %v\n\n", kind)
	return
}
