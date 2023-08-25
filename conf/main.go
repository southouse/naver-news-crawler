package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Telegram struct {
	ApiKey string `yaml:"apiKey"`
	ChatId string `yaml:"chatId"`
}

type Config struct {
	Telegram Telegram `yaml:"telegram"`
}

func GetConfig(fileName string) (*Config, error) {
	fileAbs, _ := filepath.Abs(fileName)
	yamlFile, err := os.ReadFile(fileAbs)
	if err != nil {
		log.Fatalln(err)
	}
	config := &Config{}

	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		log.Fatalln(err)
	}

	return config, nil
}
