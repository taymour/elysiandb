package configuration

import (
	"os"

	"github.com/taymour/elysiandb/internal/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Folder string `yaml:"folder"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("error:", err)
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal("error:", err)
		return nil, err
	}

	return &config, nil
}
