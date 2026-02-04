package internal

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey       string `yaml:"api_key"`
	Model        string `yaml:"model"`
	Conventional bool   `yaml:"conventional"`
	GitMoji      bool   `yaml:"gitmoji"`
	Why          bool   `yaml:"why"`
	Language     string `yaml:"language"`
	MaxLength    int    `yaml:"max_length"`
}

func DefaultConfig() *Config {
	return &Config{
		Model:        "gpt-4o-mini",
		Conventional: true,
		GitMoji:      false,
		Language:     "en",
		MaxLength:    72,
	}
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gcommit", "config.yaml"), nil
}

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	path, err := ConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
