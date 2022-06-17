package givematlib

type Extractor interface {
	ExtractLearnables(text string) []string
}

type KanjiExtractor struct{}
type ExternalExtractor struct{}

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
	// TODO: Call external program to run extraction on text
	return []string{}
}
