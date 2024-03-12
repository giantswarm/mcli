package capvcd

type Config struct {
	Namespace   string
	CloudConfig string
}

func GetCAPVCDConfig(file string) Config {
	return Config{}
}

func GetCAPVCDFile(c Config) (string, error) {
	return "", nil
}
