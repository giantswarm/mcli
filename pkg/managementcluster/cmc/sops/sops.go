package sops

type Config struct {
	Cluster   string
	AgePubKey string
}

func GetSopsConfig(file string) Config {
	return Config{}
}

func GetSopsFile(c Config) (string, error) {
	return "", nil
}
