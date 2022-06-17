package givematlib

type Extractor interface {
	ExtractLearnables() []string
}

type KanjiExtractor struct{}

func (e *KanjiExtractor) ExtractLearnables(text string) []string {
	return []string{}
}
