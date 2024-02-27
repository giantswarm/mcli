package defaultapps

type Config struct {
	Cluster                 string
	Name                    string
	Catalog                 string
	Version                 string
	Namespace               string
	Values                  string
	PrivateCA               bool
	Provider                string
	CertManagerDNSChallenge bool
	MCAppsPreventDeletion   bool
}

func GetDefaultAppsConfig(file string) Config {
	// todo: get the clusterapps from the file
	return Config{}
}

func GetDefaultAppsFile(c Config) (string, error) {
	// todo: get the clusterapps file from the clusterapps
	return "", nil
}
