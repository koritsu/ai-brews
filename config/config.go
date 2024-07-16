package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	DBPath        string `yaml:"db_path"`
	ServerPort    string `yaml:"server_port"`
	AdminUsername string `yaml:"admin_username"`
}

func LoadConfig() (*Config, error) {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
