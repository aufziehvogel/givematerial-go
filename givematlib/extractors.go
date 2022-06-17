package givematlib

import (
	"log"
	"strconv"
)

type Extractor interface {
	ExtractLearnables(text string) []string
}

type KanjiExtractor struct{}

func (e *KanjiExtractor) ExtractLearnables(text string) []string {
	var kanji []string

	for _, char := range text {
		number := int(char)
		lower1 := hexToInt("3400")
		upper1 := hexToInt("4dbf")
		lower2 := hexToInt("4e00")
		upper2 := hexToInt("9faf")

		if (number >= lower1 && number <= upper1) || (number >= lower2 && number <= upper2) {
			kanji = append(kanji, string(char))
		}
	}

	return kanji
}

func hexToInt(hex string) int {
	number, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		log.Fatalf("Could not convert %q to int", hex)
	}
	return int(number)
}
