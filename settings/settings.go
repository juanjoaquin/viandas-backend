package settings

import (
	_ "embed"
	"os"

	"github.com/joho/godotenv"
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
	SSLMode  string `yaml:"sslmode"`
}

type Settings struct {
	Port string         `yaml:"port"`
	DB   DatabaseConfig `yaml:"database"`
}

func New() (*Settings, error) {
	_ = godotenv.Load()

	var s Settings

	if err := yaml.Unmarshal(settingsFile, &s); err != nil {
		return nil, err
	}

	applyEnvOverrides(&s)

	return &s, nil
}

func applyEnvOverrides(s *Settings) {
	if v := os.Getenv("PORT"); v != "" {
		s.Port = v
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
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		s.DB.SSLMode = v
	}
}
