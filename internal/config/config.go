package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	FB             `yaml:"fb"`
	Bot            `yaml:"bot"`
	DatabaseConfig `yaml:"db"`
}
type FB struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
	Secret string `yaml:"secret"`
	Prompt string `yaml:"prompt"`
}
type Bot struct {
	Token   string `yaml:"token"`
	Channel string `yaml:"channel"`
	ChatId  int64  `yaml:"chat_id"`
	Timeout int    `yaml:"timeout"`
}

type DatabaseConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Name           string `yaml:"name"`
	MaxConnections int32  `yaml:"maxConnections"`
	Sslmode        string `yaml:"sslmode"`
}

func New(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
