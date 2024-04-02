package pushcmc

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
)

const (
	CertManagerRegionKey         = "cert_manager_route53_region"
	CertManagerRoleKey           = "cert_manager_route53_role_base"
	CertManagerAccessKey         = "cert_manager_route53_accesskey"
	CertManagerSecretKey         = "cert_manager_route53_secretaccesskey"
	CloudDirectorRefreshTokenKey = "REFRESH_TOKEN"
	AzureClientIDKey             = "clientId"
	AzureClientSecretKey         = "clientSecret"
	AzureTenantIDKey             = "tenantId"
	azureClientIDUA              = "UA_clientId"
	azureTenantIDUA              = "UA_tenantId"
	azureResourceIDUA            = "UA_id"
)

const (
	CloudDirectorCredentialsFile = "cloud-director.sh"
	VsphereCredentialsFile       = "vsphere-credentials.yaml"
	AzureCredentialsFile         = "capz.secrets.sh"
	AzureIdentityFile            = "capz.identity.sh"
)

func GetSecrets(cluster string) []string {
	return []string{
		key.ClusterValuesFile,
		key.GetDeployKey(cluster),
		key.GetPassphrase(cluster),
		key.GetKnownHosts(cluster),
		key.GetDeployKey(fmt.Sprintf("%s-ccr", cluster)),
		key.GetPassphrase(fmt.Sprintf("%s-ccr", cluster)),
		key.GetKnownHosts(fmt.Sprintf("%s-ccr", cluster)),
		key.GetDeployKey(fmt.Sprintf("%s-scr", cluster)),
		key.GetPassphrase(fmt.Sprintf("%s-scr", cluster)),
	}
}

// We read the secrets from the provided secrets folder location to get the values
// This is assuming that the secrets are created beforehand or pulled from lastpass within mc-bootstrap
// Todo: Implement the actual secret management within mcli

func (c *Config) ReadSecretFlags() error {
	if c.Input != nil {
		log.Debug().Msg("Input file provided, skipping reading secrets folder")
		return nil
	}

	var secrets map[string]string
	var err error
	{
		commonSecretsFile := fmt.Sprintf("%s/%s", c.Flags.SecretFolder, key.CommonSecretsFile)
		clusterSecretFile := fmt.Sprintf("%s/%s", c.Flags.SecretFolder, key.GetClusterSecretFile(c.Cluster))

		// check if the two files exist inside the secrets folder
		if _, err := os.Stat(commonSecretsFile); err != nil {
			return fmt.Errorf("common secrets file %s can not be accessed", commonSecretsFile)
		}
		if _, err := os.Stat(clusterSecretFile); err != nil {
			return fmt.Errorf("cluster secrets file %s can not be accessed", clusterSecretFile)
		}
		// read the secrets from the files
		secrets, err = readFlagsFromFile(commonSecretsFile)
		if err != nil {
			return err
		}
		clusterSecrets, err := readFlagsFromFile(clusterSecretFile)
		if err != nil {
			return err
		}
		// merge the common and cluster secrets - cluster secrets will override common secrets
		for k, v := range clusterSecrets {
			secrets[k] = v
		}
	}

	// these secret files are compulsory and should exist
	for _, secret := range GetSecrets(c.Cluster) {
		v, err := c.ReadFileFromSecretFolder(secret)
		if err != nil {
			return err
		}
		secrets[secret] = v
	}

	if key.IsProviderVsphere(c.Provider) {
		vsphereCredentialsFile, err := c.ReadFileFromSecretFolder(VsphereCredentialsFile)
		if err != nil {
			return err
		}
		secrets[VsphereCredentialsFile] = vsphereCredentialsFile
	} else if key.IsProviderVCD(c.Provider) {
		clouddirectorfile := fmt.Sprintf("%s/%s", c.Flags.SecretFolder, CloudDirectorCredentialsFile)
		vcdFlags, err := readFlagsFromFile(clouddirectorfile)
		if err != nil {
			return err
		}
		secrets[CloudDirectorRefreshTokenKey] = vcdFlags[CloudDirectorRefreshTokenKey]
	} else if key.IsProviderAzure(c.Provider) {
		azurefile := fmt.Sprintf("%s/%s", c.Flags.SecretFolder, AzureCredentialsFile)
		azureFlags, err := readFlagsFromFile(azurefile)
		if err != nil {
			return err
		}
		secrets[AzureClientIDKey] = azureFlags[AzureClientIDKey]
		secrets[AzureClientSecretKey] = azureFlags[AzureClientSecretKey]
		secrets[AzureTenantIDKey] = azureFlags[AzureTenantIDKey]
		azureidentity := fmt.Sprintf("%s/%s", c.Flags.SecretFolder, AzureIdentityFile)
		azureIdentityFlags, err := readFlagsFromFile(azureidentity)
		if err != nil {
			return err
		}
		secrets[azureClientIDUA] = azureIdentityFlags[azureClientIDUA]
		secrets[azureTenantIDUA] = azureIdentityFlags[azureTenantIDUA]
		secrets[azureResourceIDUA] = azureIdentityFlags[azureResourceIDUA]
	}

	if c.Flags.ConfigureContainerRegistries {
		containerRegistry, err := c.ReadFileFromSecretFolder(key.GetContainerRegistriesFile(c.Cluster))
		if err != nil {
			return err
		}
		secrets[key.GetContainerRegistriesFile(c.Cluster)] = containerRegistry
	}
	return c.SetSecretFlags(secrets)
}

func (c *Config) SetSecretFlags(secrets map[string]string) error {
	for k, v := range secrets {
		switch k {
		case CertManagerRegionKey:
			if c.Flags.Secrets.CertManagerRoute53Region == "" {
				c.Flags.Secrets.CertManagerRoute53Region = v
			}
		case CertManagerRoleKey:
			if c.Flags.Secrets.CertManagerRoute53Role == "" {
				c.Flags.Secrets.CertManagerRoute53Role = fmt.Sprintf("%s-%s", v, c.Cluster)
			}
		case CertManagerAccessKey:
			if c.Flags.Secrets.CertManagerRoute53AccessKeyID == "" {
				c.Flags.Secrets.CertManagerRoute53AccessKeyID = v
			}
		case CertManagerSecretKey:
			if c.Flags.Secrets.CertManagerRoute53SecretAccessKey == "" {
				c.Flags.Secrets.CertManagerRoute53SecretAccessKey = v
			}
		case key.ClusterValuesFile:
			if c.Flags.Secrets.ClusterValues == "" {
				c.Flags.Secrets.ClusterValues = v
			}
		case key.GetDeployKey(c.Cluster):
			if c.Flags.Secrets.SSHDeployKey.Identity == "" {
				c.Flags.Secrets.SSHDeployKey.Identity = v
			}
		case key.GetPassphrase(c.Cluster):
			if c.Flags.Secrets.SSHDeployKey.Passphrase == "" {
				c.Flags.Secrets.SSHDeployKey.Passphrase = v
			}
		case key.GetKnownHosts(c.Cluster):
			if c.Flags.Secrets.SSHDeployKey.KnownHosts == "" {
				c.Flags.Secrets.SSHDeployKey.KnownHosts = v
			}
		case key.GetDeployKey(fmt.Sprintf("%s-ccr", c.Cluster)):
			if c.Flags.Secrets.CustomerDeployKey.Identity == "" {
				c.Flags.Secrets.CustomerDeployKey.Identity = v
			}
		case key.GetPassphrase(fmt.Sprintf("%s-ccr", c.Cluster)):
			if c.Flags.Secrets.CustomerDeployKey.Passphrase == "" {
				c.Flags.Secrets.CustomerDeployKey.Passphrase = v
			}
		case key.GetKnownHosts(fmt.Sprintf("%s-ccr", c.Cluster)):
			if c.Flags.Secrets.CustomerDeployKey.KnownHosts == "" {
				c.Flags.Secrets.CustomerDeployKey.KnownHosts = v
				c.Flags.Secrets.SharedDeployKey.KnownHosts = v
			}
		case key.GetDeployKey(fmt.Sprintf("%s-scr", c.Cluster)):
			if c.Flags.Secrets.SharedDeployKey.Identity == "" {
				c.Flags.Secrets.SharedDeployKey.Identity = v
			}
		case key.GetPassphrase(fmt.Sprintf("%s-scr", c.Cluster)):
			if c.Flags.Secrets.SharedDeployKey.Passphrase == "" {
				c.Flags.Secrets.SharedDeployKey.Passphrase = v
			}
		case VsphereCredentialsFile:
			if c.Flags.Secrets.VSphereCredentials == "" {
				c.Flags.Secrets.VSphereCredentials = v
			}
		case CloudDirectorRefreshTokenKey:
			if c.Flags.Secrets.CloudDirectorRefreshToken == "" {
				c.Flags.Secrets.CloudDirectorRefreshToken = v
			}
		case AzureClientIDKey:
			if c.Flags.Secrets.Azure.ClientID == "" {
				c.Flags.Secrets.Azure.ClientID = v
			}
		case AzureClientSecretKey:
			if c.Flags.Secrets.Azure.ClientSecret == "" {
				c.Flags.Secrets.Azure.ClientSecret = v
			}
		case AzureTenantIDKey:
			if c.Flags.Secrets.Azure.TenantID == "" {
				c.Flags.Secrets.Azure.TenantID = v
			}
		case azureClientIDUA:
			if c.Flags.Secrets.Azure.UAClientID == "" {
				c.Flags.Secrets.Azure.UAClientID = v
			}
		case azureTenantIDUA:
			if c.Flags.Secrets.Azure.UATenantID == "" {
				c.Flags.Secrets.Azure.UATenantID = v
			}
		case azureResourceIDUA:
			if c.Flags.Secrets.Azure.UAResourceID == "" {
				c.Flags.Secrets.Azure.UAResourceID = v
			}
		case key.GetContainerRegistriesFile(c.Cluster):
			if c.Flags.Secrets.ContainerRegistryConfiguration == "" {
				c.Flags.Secrets.ContainerRegistryConfiguration = v
			}
		default:
			log.Debug().Msgf("secret flag %s does not exist or is already set", k)
		}
	}
	log.Debug().Msg("secrets set")
	return nil
}

func (c *Config) ReadFileFromSecretFolder(file string) (string, error) {
	path := fmt.Sprintf("%s/%s", c.Flags.SecretFolder, file)
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("file %s can not be accessed", path)
	}
	return readFile(path)
}

func readFlagsFromFile(file string) (map[string]string, error) {
	// read the file and return the key value pairs
	flags := make(map[string]string)

	s, err := readFile(file)
	if err != nil {
		return nil, err
	}
	// regex to parse the key value pairs key=value
	re := regexp.MustCompile(`(\w+)=['"]?(\S+)`)
	pairs := re.FindAllStringSubmatch(s, -1)
	for _, pair := range pairs {
		value, _ := strings.CutSuffix(pair[2], "'")
		value, _ = strings.CutSuffix(value, "\"")
		value, _ = strings.CutSuffix(value, "\n")
		flags[pair[1]] = value
	}
	return flags, nil
}

func readFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", file, err)
	}
	s := string(data)
	return s, nil
}
