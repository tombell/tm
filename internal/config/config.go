package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name string `yaml:"name"`
	Root string `yaml:"root"`

	Sessions []struct {
		Name    string `yaml:"name"`
		Root    string `yaml:"root"`
		Windows []struct {
			Name  string `yaml:"name"`
			Root  string `yaml:"root"`
			Panes []struct {
				Root string `yaml:"root"`
			} `yaml:"panes"`
		} `yaml:"windows"`
	} `yaml:"sessions"`
}

func Load(name string) (*Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read file failed: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("yaml unmarshal failed: %w", err)
	}

	return &cfg, nil
}
