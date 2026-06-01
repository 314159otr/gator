package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)
 const CONFIG_FILE_NAME = ".gatorconfig.json"

type Config struct {
	DblURL          string `json:db_url`
	CurrentUserName string `json:current_user_name`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (config *Config) SetUser(user string) error {
	config.CurrentUserName = user
	return write(config)
}

func write(config *Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)

	if err := encoder.Encode(config); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, CONFIG_FILE_NAME)
	return path, nil
}
