package deploykey

type Config struct {
	AgePubKey  string
	Key        string
	Identity   string
	KnownHosts string
}

func GetDeployKeyConfig(file string) Config {
	return Config{}
}

func GetDeployKeyFile(c Config) (string, error) {
	return "", nil
}
