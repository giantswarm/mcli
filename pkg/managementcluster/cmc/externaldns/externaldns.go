package externaldns

import "github.com/rs/zerolog/log"

func GetExternalDNSFile() string {
	log.Debug().Msg("Creating externalDNS file")
	return `kind: ConfigMap
apiVersion: v1
metadata:
  annotations:
    meta.helm.sh/release-name: external-dns
    meta.helm.sh/release-namespace: kube-system
  labels:
    app.kubernetes.io/managed-by: Helm
    k8s-app: external-dns
  name: external-dns-user-values
  namespace: kube-system
data:
  values: |
flavor: capi
provider: azure
clusterID: {{ .Values.clusterName }}
crd:
  install: false
externalDNS:
  namespaceFilter: \"\"
  sources:
  - ingress`
}
