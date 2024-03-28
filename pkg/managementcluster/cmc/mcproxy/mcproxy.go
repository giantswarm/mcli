package mcproxy

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Hostname string
	Port     string
}

const (
	ProxyHostnameKey = "proxy_hostname"
	ProxyPortKey     = "proxy_port"
)

func GetHTTPSProxy(kustomizationFile string) (Config, error) {
	log.Debug().Msg("Getting HTTPS proxy configuration")

	hostname, err := key.GetValue(ProxyHostnameKey, kustomizationFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get proxy hostname.\n%w", err)
	}

	port, err := key.GetValue(ProxyPortKey, kustomizationFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get proxy port.\n%w", err)
	}

	return Config{
		Hostname: hostname,
		Port:     port,
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
	log.Debug().Msg("Creating Kustomization for proxy")
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
