package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Pane struct {
	Type     string   `yaml:"type"`
	Root     string   `yaml:"root"`
	Commands []string `yaml:"commands"`
}

type Window struct {
	Name     string   `yaml:"name"`
	Root     string   `yaml:"root"`
	Layout   string   `yaml:"layout"`
	Commands []string `yaml:"commands"`
	Panes    []Pane   `yaml:"panes"`
}

type Session struct {
	Name        string   `yaml:"name"`
	Root        string   `yaml:"root"`
	BeforeStart []string `yaml:"before_start"`
	AfterStart  []string `yaml:"after_start"`
	Windows     []Window `yaml:"windows"`
}

type Config struct {
	Root        string    `yaml:"root"`
	BeforeStart []string  `yaml:"before_start"`
	AfterStop   []string  `yaml:"after_stop"`
	Sessions    []Session `yaml:"sessions"`
}

func Load(name string) (*Config, error) {
	if strings.HasPrefix(name, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		name = strings.Replace(name, "~", homeDir, 1)
	}

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
