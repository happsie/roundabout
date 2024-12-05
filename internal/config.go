package internal

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Config struct {
	Port              string    `yaml:"port"`
	DefaultTargetHost string    `yaml:"defaultTargetHost"`
	Services          []Service `yaml:"services,flow"`
}

type Service struct {
	Name       string   `yaml:"name"`
	Paths      []string `yaml:"paths,flow"`
	TargetHost string   `yaml:"targetHost"`
}

func LoadConfig(path string) (*Config, error) {
	slog.Debug("reading config.yml")
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = yaml.Unmarshal(f, conf)
	if err != nil {
		return nil, err
	}
	slog.Debug("config loaded successfully", "config", *conf)
	return conf, nil
}

func (c *Config) Save() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile("config.yml", b, os.ModePerm)
	if err != nil {
		return err
	}
	slog.Debug("config file saved", "raw_string", string(b))
	return nil
}
