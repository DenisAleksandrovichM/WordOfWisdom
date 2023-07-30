package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerHost         string `envconfig:"SERVER_HOST"`
	ServerPort         int    `envconfig:"SERVER_PORT"`
	HashcashZerosCount int    `envconfig:"HASHCASH_ZEROS_COUNT"`
	IncZerosCountLimit int    `envconfig:"INC_ZEROS_COUNT_LIMIT"`
}

func Load(path string) (*Config, error) {
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fileSafeClose(file)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return &config, err
	}
	err = envconfig.Process("", &config)
	return &config, err
}

func fileSafeClose(file *os.File) {
	if err := file.Close(); err != nil {
		log.Printf("Error closing file: %s\n", err)
	}
}
