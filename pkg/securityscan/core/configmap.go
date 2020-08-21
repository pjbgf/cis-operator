package core

import (
	"bytes"
	"encoding/json"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"

	"github.com/rancher/wrangler/pkg/name"

	cisoperatorapiv1 "github.com/rancher/clusterscan-operator/pkg/apis/securityscan.cattle.io/v1"
)

type OverrideSkipInfoData struct {
	Skip map[string][]string `json:"skip"`
}

var (
	CurrentBenchmarkKey = "current"
	ConfigFileName      = "config.json"
)

func NewConfigMaps(clusterscan *cisoperatorapiv1.ClusterScan, clusterscanprofile *cisoperatorapiv1.ClusterScanProfile, controllerName string, imageConfig *cisoperatorapiv1.ScanImageConfig) (configmaps []*corev1.ConfigMap, err error) {

	configdata := map[string]interface{}{
		"namespace":        cisoperatorapiv1.ClusterScanNS,
		"name":             name.SafeConcatName(cisoperatorapiv1.ClusterScanConfigMap, clusterscan.Name),
		"runName":          name.SafeConcatName("security-scan-runner", clusterscan.Name),
		"appName":          "rancher-cis-benchmark",
		"advertiseAddress": cisoperatorapiv1.ClusterScanService,
		"sonobuoyImage":    imageConfig.SonobuoyImage + ":" + imageConfig.SonobuoyImageTag,
		"sonobuoyVersion":  imageConfig.SonobuoyImageTag,
	}
	configcm, err := generateConfigMap(clusterscan, "cisscanConfig.template", "./pkg/securityscan/core/templates/cisscanConfig.template", configdata)
	if err != nil {
		return configmaps, err
	}
	plugindata := map[string]interface{}{
		"namespace":         cisoperatorapiv1.ClusterScanNS,
		"name":              name.SafeConcatName(cisoperatorapiv1.ClusterScanPluginsConfigMap, clusterscan.Name),
		"runName":           name.SafeConcatName("security-scan-runner", clusterscan.Name),
		"appName":           "rancher-cis-benchmark",
		"serviceaccount":    cisoperatorapiv1.ClusterScanSA,
		"securityScanImage": imageConfig.SecurityScanImage + ":" + imageConfig.SecurityScanImageTag,
		"benchmarkVersion":  clusterscanprofile.Spec.BenchmarkVersion,
	}
	plugincm, err := generateConfigMap(clusterscan, "pluginConfig.template", "./pkg/securityscan/core/templates/pluginConfig.template", plugindata)
	if err != nil {
		return configmaps, err
	}

	var skipConfigcm *corev1.ConfigMap
	if clusterscanprofile.Spec.SkipTests != nil && len(clusterscanprofile.Spec.SkipTests) > 0 {
		//create user skip config map as well
		// create the cm
		skipDataBytes, err := getOverrideSkipInfoData(clusterscanprofile.Spec.SkipTests)
		if err != nil {
			return configmaps, err
		}
		skipConfigcm = getConfigMapObject(getOverrideConfigMapName(clusterscan), string(skipDataBytes))
	}

	configmaps = append(configmaps, configcm, plugincm, skipConfigcm)
	return configmaps, nil
}

func generateConfigMap(clusterscan *cisoperatorapiv1.ClusterScan, templateName string, templateFile string, data map[string]interface{}) (*corev1.ConfigMap, error) {
	configcm := &corev1.ConfigMap{}

	obj, err := parseTemplate(clusterscan, templateName, templateFile, data)
	if err != nil {
		return nil, err
	}

	if err := obj.Decode(&configcm); err != nil {
		return nil, err
	}
	return configcm, nil
}

func parseTemplate(clusterscan *cisoperatorapiv1.ClusterScan, templateName string, templateFile string, data map[string]interface{}) (*k8Yaml.YAMLOrJSONDecoder, error) {
	cmTemplate, err := template.New(templateName).ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = cmTemplate.Execute(&b, data)
	if err != nil {
		return nil, err
	}

	return k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(b.String())), 1000), nil
}

func getOverrideConfigMapName(cs *cisoperatorapiv1.ClusterScan) string {
	return name.SafeConcatName(cisoperatorapiv1.ClusterScanUserSkipConfigMap, cs.Name)
}

func getOverrideSkipInfoData(skip []string) ([]byte, error) {
	s := OverrideSkipInfoData{Skip: map[string][]string{CurrentBenchmarkKey: skip}}
	return json.Marshal(s)
}

func getConfigMapObject(cmName, data string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: cisoperatorapiv1.ClusterScanNS,
		},
		Data: map[string]string{
			ConfigFileName: data,
		},
	}
}