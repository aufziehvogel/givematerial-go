package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type applicationConfig struct {
	AnkiFile             string              `json:"anki_file"`
	AnkiDecksForLanguage map[string][]string `json:"anki_decks"`
	WanikaniApiToken     string              `json:"wanikani_api_token"`
}

func loadConfig(path string) (applicationConfig, error) {
	var config applicationConfig

	file, err := os.Open(path)

	if err != nil {
		return applicationConfig{}, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	err = json.NewDecoder(reader).Decode(&config)
	return config, err
}
