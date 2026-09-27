package config

import (
	"encoding/json"
	"os"

	"github.com/mohammad-rizwan-hussain/gator/internal/database"
)

type Config struct {
	DB_URL      string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

type State struct {
	DB     *database.Queries
	Config *Config
}

// const CONFIG_FILE_PATH = "/mnt/c/Users/perve/OneDrive/Desktop/Boot.dev/Gator/gatorconfig.json"
const CONFIG_FILE_PATH = "/home/mohammad/.gatorconfig.json"

func Read() (*Config, error) {
	// Read the config file
	file, err := os.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (cfg *Config) SetUser(current_user_name string) error {
	cfg.CurrentUser = current_user_name

	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	return CONFIG_FILE_PATH, nil
}

func write(cfg Config) error {
	file, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(CONFIG_FILE_PATH, file, 0644)
}
