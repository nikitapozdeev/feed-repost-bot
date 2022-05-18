package config

type Config struct {
	Telegram Telegram `yaml:"telegram"`
	VK VK `yaml:"vk"`
}

type Telegram struct {
	Token string `yaml:"token"`
}

type VK struct {
	Host string `yaml:"host"`
	BasePath string `yaml:"basePath"`
	Version string `yaml:"version"`
	Token string `yaml:"token"`
}