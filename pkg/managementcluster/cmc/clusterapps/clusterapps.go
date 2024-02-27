package clusterapps

type Config struct {
	Cluster                      string
	Name                         string
	Catalog                      string
	Version                      string
	Namespace                    string
	Values                       string
	ConfigureContainerRegistries bool
	Provider                     string
	MCAppsPreventDeletion        bool
}

func GetClusterAppsConfig(file string) Config {
	// todo: get the clusterapps from the file
	return Config{}
}

func GetClusterAppsFile(c Config) (string, error) {
	// todo: get the clusterapps file from the clusterapps
	return "", nil
}
