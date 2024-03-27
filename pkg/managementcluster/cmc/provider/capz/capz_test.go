package capz

import (
	"reflect"
	"testing"
)

func TestGetCAPZConfig(t *testing.T) {
	testCases := []struct {
		name string
		sp   string
		ua   string
		sec  string

		expected    Config
		expectError bool
	}{
		{
			name: "case 0: simple",
			sp:   getValidCAPZSPInput(),
			ua:   getValidCAPZUAInput(),
			sec:  getValidCAPZSecretInput(),

			expected: Config{
				UAClientID:   "test-ua-client-id",
				UATenantID:   "test-ua-tenant-id",
				UAResourceID: "test/ua/resource/id",
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-id",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetCAPZConfig(tc.sp, tc.ua, tc.sec)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
				if !reflect.DeepEqual(actual, tc.expected) {
					t.Fatalf("expected %#v but got %#v", tc.expected, actual)
				}
			}
		})
	}
}

func getValidCAPZSPInput() string {
	return `apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AzureClusterIdentity
metadata:
  name: cluster-identity-static-sp
  namespace: test-namespace
  labels:
    clusterctl.cluster.x-k8s.io/move: "true"
spec:
  allowedNamespaces: {}
  clientID: test-client-id
  clientSecret:
    name: cluster-identity-secret-static-sp
    namespace: test-namespace
  tenantID: test-tenant-id
  type: ManualServicePrincipal`
}

func getValidCAPZUAInput() string {
	return `apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AzureClusterIdentity
metadata:
  name: cluster-identity
  namespace: test-namespace
spec:
  allowedNamespaces: {}
  clientID: test-ua-client-id
  tenantID: test-ua-tenant-id
  resourceID: test/ua/resource/id
  type: UserAssignedMSI`
}

func getValidCAPZSecretInput() string {
	return `apiVersion: v1
kind: Secret
metadata:
  name: cluster-identity-secret-static
  namespace: test-namespace
  labels:
    clusterctl.cluster.x-k8s.io/move: "true"
type: Opaque
data:
  clientSecret: test-client-secret`
}
