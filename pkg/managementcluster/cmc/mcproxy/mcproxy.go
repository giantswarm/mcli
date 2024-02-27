package mcproxy

type Config struct {
	HostName string
	Port     int
}

func GetMCProxyConfig(allowNetPolFile string, sourceControllerFile string) Config {
	return Config{}
}

func GetAllowNetPolFile(c Config) (string, error) {
	return "", nil
}

func GetSourceControllerFile(c Config) (string, error) {
	return "", nil
}
