package givematlib

import (
	"bufio"
	"encoding/json"
	"os"
)

type ApplicationConfig struct {
	AnkiFile             string              `json:"anki_file"`
	AnkiDecksForLanguage map[string][]string `json:"anki_decks"`
	WanikaniApiToken     string              `json:"wanikani_api_token"`
}

func LoadConfig(path string) (ApplicationConfig, error) {
	var config ApplicationConfig

	file, err := os.Open(path)

	if err != nil {
		return ApplicationConfig{}, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	err = json.NewDecoder(reader).Decode(&config)
	return config, err
}
