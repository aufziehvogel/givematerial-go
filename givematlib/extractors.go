package givematlib

import (
	"log"
	"os/exec"
	"strings"
)

type Extractor interface {
	ExtractLearnables(text string) []string
}

type KanjiExtractor struct{}
type ExternalExtractor struct {
	programCall []string
}

func (e *KanjiExtractor) ExtractLearnables(text string) []string {
	var kanji []string

	for _, char := range text {
		if (char >= '\u3400' && char <= '\u4dbf') || (char >= '\u4e00' && char <= '\u9faf') {
			kanji = append(kanji, string(char))
		}
	}

	return kanji
}

func (e *ExternalExtractor) ExtractLearnables(text string) []string {
	commandName := e.programCall[0]
	args := e.programCall[1:]
	cmd := exec.Command(commandName, args...)

	cmd.Stdin = strings.NewReader(text)
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	return strings.Split(string(output), "\n")
}

func GetExtractorForLanguage(language Language) Extractor {
	switch language {
	case LANG_JAPANESE:
		return new(KanjiExtractor)
	case LANG_SPANISH:
		return &ExternalExtractor{
			programCall: []string{"python", "contrib/lemmatizer.py", "es"},
		}
	case LANG_CROATIAN:
		return &ExternalExtractor{
			programCall: []string{"python", "contrib/lemmatizer.py", "hr"},
		}
	default:
		log.Panicf("Language %v is not supported", language)
		return nil
	}
}
