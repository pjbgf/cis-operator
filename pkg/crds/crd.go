package crds

import (
	"fmt"
	"io/ioutil"
	"strings"

	cisoperator "github.com/rancher/clusterscan-operator/pkg/apis/securityscan.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/crd"
	_ "github.com/rancher/wrangler/pkg/generated/controllers/apiextensions.k8s.io" //using init
	"github.com/rancher/wrangler/pkg/yaml"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func WriteCRD() error {
	for _, crdDef := range List() {
		bCrd, err := crdDef.ToCustomResourceDefinition()
		if err != nil {
			return err
		}
		yamlBytes, err := yaml.Export(&bCrd)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("./crds/%s.yaml", strings.ToLower(bCrd.Spec.Names.Kind))
		err = ioutil.WriteFile(filename, yamlBytes, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func List() []crd.CRD {
	return []crd.CRD{
		newCRD(&cisoperator.ClusterScan{}, func(c crd.CRD) crd.CRD {
			return c
		}),
		newCRD(&cisoperator.ClusterScanProfile{}, func(c crd.CRD) crd.CRD {
			return c
		}),
		newCRD(&cisoperator.ClusterScanReport{}, func(c crd.CRD) crd.CRD {
			return c
		}),
		newCRD(&cisoperator.ScheduledScan{}, func(c crd.CRD) crd.CRD {
			return c
		}),
	}
}

func newCRD(obj interface{}, customize func(crd.CRD) crd.CRD) crd.CRD {
	crd := crd.CRD{
		GVK: schema.GroupVersionKind{
			Group:   "securityscan.cattle.io",
			Version: "v1",
		},
		NonNamespace: true,
		Status:       true,
		SchemaObject: obj,
	}
	if customize != nil {
		crd = customize(crd)
	}
	return crd
}