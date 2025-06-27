package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DdbUrl          string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshalling: %v", err)
	}
	file, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	fullFile := home + "/" + configFileName
	return fullFile, nil
}

func Read() (Config, error) {
	fullFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonData, err := os.ReadFile(fullFile)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	err = json.Unmarshal(jsonData, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}
