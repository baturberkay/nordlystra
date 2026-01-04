package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"nordlystra.app/bridge"
)

const configFile = "./config/hue_config.json"

type Config struct {
	BridgeIP string `json:"bridge_ip"`
	Username string `json:"username"`
}

func Load() (Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func Save(cfg Config) error {
	dir := "./config"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(&cfg)
}

func LoadOrCreate() Config {
	cfg, err := Load()
	if err != nil {
		fmt.Println("No configuration found. Please enter the Hue Bridge IP to set up:")
		fmt.Print("Enter Hue Bridge IP: ")

		var bridgeIP string
		fmt.Scanln(&bridgeIP)

		username, err := bridge.Authenticate(bridgeIP)
		if err != nil {
			log.Fatalf("Authentication failed: %v", err)
		}

		cfg = Config{BridgeIP: bridgeIP, Username: username}
		if err := Save(cfg); err != nil {
			log.Fatalf("Failed to save configuration: %v", err)
		}

		fmt.Println("Configuration saved successfully.")
	}

	return cfg
}
