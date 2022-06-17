package givematlib

type Language string

func (l Language) ShortCode() string {
	return string(l)
}

func MakeLanguageFromShortCode(shortCode string) Language {
	// TOOD: Handle situation when cannot be converted
	return Language(shortCode)
}

const (
	LANG_SPANISH  Language = "es"
	LANG_CROATIAN Language = "hr"
	LANG_JAPANESE Language = "ja"
)
