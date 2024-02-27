package taylorbot

type Config struct {
	User      string
	Token     string
	AgePubKey string
}

func GetTaylorBotConfig(file string) Config {
	return Config{}
}

func GetTaylorBotFile(c Config) (string, error) {
	return "", nil
}
