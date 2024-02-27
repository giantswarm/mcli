package certmanager

type Config struct {
	Region          string
	Role            string
	AccessKeyID     string
	SecretAccessKey string
}

func GetCertManagerConfig(file string) Config {
	return Config{}
}

func GetCertManagerFile(c Config) (string, error) {
	return "", nil
}
