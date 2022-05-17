package config

type Telegram struct {
	Token string `yaml:"token" env:"TOKEN"`
}

type Config struct {
	Telegram `yaml:"telegram"`
}
