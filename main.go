package main

import (
	"givematerial/givematlib"
	"givematerial/gui"

	"log"
	"os"
	"path/filepath"
)

func configPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configFile := filepath.Join(
		homeDir,
		".config",
		"givematerial",
		"givematerial.json",
	)
	return configFile
}

func main() {
	configFile := configPath()
	config, err := givematlib.LoadConfig(configFile)
	if err == nil {
		log.Printf("Loaded configuration from %s", configFile)
	} else {
		log.Panicf("Could not load configuration from %s", configFile)
	}

	gui.Init(&config)
}
