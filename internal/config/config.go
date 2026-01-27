package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName string = ".gatorconfig.json"

func (cfg Config) Read() (Config, error) {
	configFilename, err := getFileName()
	f, err := os.Open(configFilename)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	cfgFile, err := io.ReadAll(f)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(cfgFile, &cfg)
	if err != nil {
		fmt.Println("Error:", err)
		return Config{}, err
	}
	return cfg, nil
}

func (cfg Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	if err := write(cfg); err != nil {
		fmt.Printf("Error while writing JSON object data.")
		return err
	}
	return nil
}

func write(cfg Config) error {
	configBlob, err := json.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling Config struct as JSON")
		return err
	}
	configPath, err := getFileName()
	if err != nil {
		fmt.Printf("Error setting output filepath.")
		return err
	}
	if err = os.WriteFile(configPath, configBlob, 0644); err != nil {
		fmt.Printf("Error writing JSON object to file.")
		return err
	}
	return nil
}

func getFileName() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error fetching home directory address.")
		return "", err
	}
	configPath := homeDir + "/" + configFileName
	return configPath, nil
}
