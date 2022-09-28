package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name        string   `yaml:"name"`
	Root        string   `yaml:"root"`
	BeforeStart []string `yaml:"before_start"`

	Sessions []struct {
		Name     string   `yaml:"name"`
		Root     string   `yaml:"root"`
		Commands []string `yaml:"commands"`

		Windows []struct {
			Name     string   `yaml:"name"`
			Root     string   `yaml:"root"`
			Layout   string   `yaml:"layout"`
			Commands []string `yaml:"commands"`

			Panes []struct {
				Type     string   `yaml:"type"`
				Root     string   `yaml:"root"`
				Commands []string `yaml:"commands"`
			} `yaml:"panes"`
		} `yaml:"windows"`
	} `yaml:"sessions"`
}

func Load(name string) (*Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
