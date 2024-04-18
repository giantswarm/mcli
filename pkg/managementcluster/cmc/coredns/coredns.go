package coredns

import (
	"fmt"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func GetCoreDNSFile(config string) (string, error) {
	log.Debug().Msg("Creating CoreDNS configmap")
	configmap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "coredns",
			Namespace: "kube-system",
			Annotations: map[string]string{
				"meta.helm.sh/release-name":      "coredns",
				"meta.helm.sh/release-namespace": "kube-system",
			},
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "Helm",
				"k8s-app":                      "coredns",
			},
		},
		Data: map[string]string{
			"Corefile": config,
		},
	}
	data, err := yaml.Marshal(configmap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CoreDNS configmap.\n%w", err)
	}
	return string(data), nil
}

func GetCoreDNSValues(file string) (string, error) {
	log.Debug().Msg("Creating CoreDNS values")

	var configmap v1.ConfigMap
	err := yaml.Unmarshal([]byte(file), &configmap)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal CoreDNS configmap.\n%w", err)
	}

	return configmap.Data["Corefile"], nil
}
