package settings

import (
	_ "embed"
	"os"

	"go.yaml.in/yaml/v3"
)

//go:embed settings.yml
var settingsFile []byte

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Settings struct {
	Port string         `yaml:"port"`
	DB   DatabaseConfig `yaml:"database"`
}

func New() (*Settings, error) {
	var s Settings

	if err := yaml.Unmarshal(settingsFile, &s); err != nil {
		return nil, err
	}

	if v := os.Getenv("DB_HOST"); v != "" {
		s.DB.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		s.DB.Port = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		s.DB.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		s.DB.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		s.DB.Name = v
	}

	return &s, nil
}
