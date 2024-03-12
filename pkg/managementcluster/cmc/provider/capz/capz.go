package capz

type Config struct {
	Namespace        string
	IdentityUA       string
	IdentitySP       string
	IdentityStaticSP string
}

func GetCAPZConfig(sp string, ua string, staticsp string) Config {
	return Config{}
}

func GetCAPZFile(c Config) (string, error) {
	return "", nil
}
