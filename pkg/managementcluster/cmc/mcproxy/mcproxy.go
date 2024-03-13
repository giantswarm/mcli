package mcproxy

import (
	"fmt"

	"github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Hostname string
	Port     string
}

func GetHTTPSProxy(kustomizationFile string) (Config, error) {
	log.Debug().Msg("Getting HTTPS proxy configuration")

	var kustomization v1beta1.Kustomization

	if err := yaml.Unmarshal([]byte(kustomizationFile), &kustomization); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal kustomization file: %w", err)
	}

	return Config{
		Hostname: kustomization.Spec.PostBuild.Substitute["proxy_hostname"],
		Port:     kustomization.Spec.PostBuild.Substitute["proxy_port"],
	}, nil
}

func GetAllowNetPolFile(c Config) string {
	log.Debug().Msg("Creating CiliumNetworkPolicy for proxy")

	giantswarm := getCiliumNetworkPolicy("giantswarm", c.Hostname, c.Port)
	kubeSystem := getCiliumNetworkPolicy("kube-system", c.Hostname, c.Port)
	return fmt.Sprintf("%s\n---\n%s", giantswarm, kubeSystem)
}

func getCiliumNetworkPolicy(namespace string, hostname string, port string) string {
	return fmt.Sprintf(`apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: allow-egress-to-proxy
  namespace: %s
spec:
  endpointSelector: {}
  egress:
    - toCIDRSet:
        - cidr: "%s/32"
      toPorts:
        - ports:
            - port: "%s"`, namespace, hostname, port)
}

func GetKustomization(c Config) string {
	return fmt.Sprintf(`apiVersion: kustomize.toolkit.fluxcd.io/v1
kind: Kustomization
metadata:
  name: flux
  namespace: flux-giantswarm
spec:
  postBuild:
	substitute:
	  proxy_hostname: %s
	  proxy_port: %s`, c.Hostname, c.Port)
}
