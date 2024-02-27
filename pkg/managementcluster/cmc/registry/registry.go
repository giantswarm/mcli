package registry

type Config struct {
	Values    string
	AgePubKey string
}

func GetRegistryConfig(file string) Config {
	return Config{}
}

func GetRegistryFile(c Config) (string, error) {
	return "", nil
}
