package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Input struct {
	Id      *string `yaml:"id"`
	Glob    string  `yaml:"glob"`
	Preload bool    `yaml:"preload"`
	Minify  bool    `yaml:"minify"`
}

type Config struct {
	Assets []Input `yaml:"assets"`
	Out    *string `yaml:"out"`
}

func ReadConfig(path string) (*Config, error) {
	configDoc, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := Config{}

	err = yaml.Unmarshal(configDoc, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
