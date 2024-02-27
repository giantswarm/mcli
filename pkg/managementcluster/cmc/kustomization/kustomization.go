package kustomization

type Config struct {
	CertManagerDNSChallenge      bool
	Provider                     string
	PrivateCA                    bool
	ConfigureContainerRegistries bool
	CustomCoreDNS                bool
	DisableDenyAllNetPol         bool
	MCProxy                      bool
}

func GetKustomizationConfig(file string) Config {
	return Config{}
}

func GetKustomizationFile(c Config) (string, error) {
	return "", nil
}
