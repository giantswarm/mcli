package capv

type Config struct {
	Namespace   string
	CloudConfig string
	AgePubKey   string
}

func GetCAPVConfig(file string) Config {
	return Config{}
}

func GetCAPVFile(c Config) (string, error) {
	return "", nil
}
