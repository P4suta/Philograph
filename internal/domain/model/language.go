package model

import "unicode"

// Language はテキストの言語を表す。
type Language int

const (
	LangUnknown  Language = iota
	LangJapanese
	LangEnglish
)

// DetectLanguage はテキストの先頭部分のUnicode符号点分布から言語を推定する。
func DetectLanguage(text string) Language {
	const sampleSize = 1000
	cjkCount, totalCount := 0, 0

	for _, r := range text {
		if totalCount >= sampleSize {
			break
		}
		if unicode.IsLetter(r) {
			totalCount++
			if unicode.Is(unicode.Han, r) ||
				unicode.Is(unicode.Hiragana, r) ||
				unicode.Is(unicode.Katakana, r) {
				cjkCount++
			}
		}
	}

	if totalCount == 0 {
		return LangUnknown
	}
	if float64(cjkCount)/float64(totalCount) > 0.3 {
		return LangJapanese
	}
	return LangEnglish
}
