package cmd

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Project      string              `yaml:"project"`
	Version      string              `yaml:"version"`
	Build        BuildSettings       `yaml:"build"`
	Dependencies DependencyConfig    `yaml:"dependencies"`
	Tasks        map[string][]string `yaml:"tasks"`
}

type BuildSettings struct {
	Default BuildTarget            `yaml:"default"`
	Targets map[string]BuildTarget `yaml:"targets"`
}

type BuildTarget struct {
	OS      string   `yaml:"os"`
	Arch    string   `yaml:"arch"`
	Output  string   `yaml:"output"`
	Ldflags string   `yaml:"ldflags,omitempty"`
	Tags    []string `yaml:"tags,omitempty"`
	Cgo     bool     `yaml:"cgo,omitempty"`
}

type DependencyConfig struct {
	Check []string `yaml:"check"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return config, nil
}
