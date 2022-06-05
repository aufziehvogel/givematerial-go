package givematlib

import (
	"fmt"
)

func UpdateAnki(
	ankiFile string,
	decksForLanguage map[string][]string,
	messageChan chan<- string,
	closeChannel bool,
) error {
	if closeChannel {
		defer close(messageChan)
	}

	for language, decks := range decksForLanguage {
		p := AnkiProvider{
			Decks:    decks,
			AnkiFile: ankiFile,
		}
		messageChan <- fmt.Sprintf(
			"Reading status from Anki for %s (decks: %v)",
			language,
			decks,
		)
		learnables, err := p.FetchLearnables()
		if err != nil {
			return err
		}

		messageChan <- fmt.Sprintf(
			"Saving Anki status for %s (%d words)",
			language,
			len(learnables),
		)
		err = SaveLearnableStatus("anki", language, learnables)
		if err != nil {
			return err
		}
	}

	messageChan <- "Anki status update complete"

	return nil
}

func UpdateWanikani(
	apiToken string,
	messageChan chan<- string,
	closeChannel bool,
) error {
	if closeChannel {
		defer close(messageChan)
	}

	w := NewWanikaniProvider(apiToken)

	messageChan <- "Reading status from Wanikani"
	learnables, err := w.FetchLearnables()
	if err != nil {
		return err
	}

	messageChan <- "Saving Wanikani status"
	err = SaveLearnableStatus("wanikani", "ja", learnables)
	messageChan <- "Wanikani status update complete"
	return err
}
