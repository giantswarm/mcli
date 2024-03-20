package issuer

import "github.com/rs/zerolog/log"

// TODO: this does not seem ideal but it's how its in mc-bootstrap and we go with it for now
func GetIssuerFile() string {
	log.Debug().Msg("Creating issuer file")
	return `apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: private-giantswarm
  labels:
	giantswarm.io/service-type: "managed"
spec:
  ca:
	secretName: private-giantswarm-secret
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: private-giantswarm-ca
  namespace: kube-system
spec:
  isCA: true
  commonName: gigantic.internal
  secretName: private-giantswarm-secret
  privateKey:
	algorithm: ECDSA
	size: 256
  issuerRef:
	name: selfsigned-giantswarm
	kind: ClusterIssuer
	group: cert-manager.io`
}
