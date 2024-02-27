package capz

type Config struct {
	Namespace        string
	IdentityUA       string
	IdentitySP       string
	IdentityStaticSP string
	AgePubKey        string
}

func GetCAPZConfig(file string) Config {
	return Config{}
}

func GetCAPZFile(c Config) (string, error) {
	return "", nil
}
