package givematlib

import "testing"

func TestKanjiExtractor(t *testing.T) {
	extractor := new(KanjiExtractor)
	learnables := extractor.ExtractLearnables("日本語を話す")
	want := []string{"日", "本", "語", "話"}

	compareLearnableEquals(learnables, want, t)
}

func TestExternalExtractor(t *testing.T) {
	extractor := ExternalExtractor{
		// Replace spaces with newlines
		programCall: []string{"tr", "' '", "'\n'"},
	}
	learnables := extractor.ExtractLearnables("this is a simple test")
	want := []string{"this", "is", "a", "simple", "test"}

	compareLearnableEquals(learnables, want, t)
}

func compareLearnableEquals(learnables []string, want []string, t *testing.T) {
	if len(learnables) != len(want) {
		t.Errorf("Expected %d kanji, got %d", len(want), len(learnables))
	}

	for i, learnable := range learnables {
		if learnable != want[i] {
			t.Errorf(
				"Extracted kanjis were incorrect, got %v, want %v",
				learnables,
				want,
			)
		}
	}
}
