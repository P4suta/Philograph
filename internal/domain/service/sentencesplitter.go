package service

import (
	"strings"
	"unicode/utf8"

	"Philograph/internal/domain/model"
)

// SentenceSplitter はテキストを文単位に分割する。
type SentenceSplitter struct {
	language model.Language
}

// NewSentenceSplitter は新しいSentenceSplitterを返す。
func NewSentenceSplitter(lang model.Language) *SentenceSplitter {
	return &SentenceSplitter{language: lang}
}

// Split はテキストを文のスライスに分割する。
func (s *SentenceSplitter) Split(text string) []string {
	if s.language == model.LangJapanese {
		return s.splitJapanese(text)
	}
	return s.splitEnglish(text)
}

func (s *SentenceSplitter) splitJapanese(text string) []string {
	var sentences []string
	var buf strings.Builder

	for _, r := range text {
		buf.WriteRune(r)
		if r == '。' || r == '！' || r == '？' || r == '\n' {
			sentence := strings.TrimSpace(buf.String())
			if sentence != "" {
				sentences = append(sentences, sentence)
			}
			buf.Reset()
		}
	}

	if remaining := strings.TrimSpace(buf.String()); remaining != "" {
		sentences = append(sentences, remaining)
	}

	return sentences
}

// englishAbbreviations は略語として扱うサフィックス。
var englishAbbreviations = []string{
	"Mr.", "Mrs.", "Ms.", "Dr.", "Prof.",
	"Jr.", "Sr.", "St.",
	"Inc.", "Corp.", "Ltd.", "Co.",
	"vs.", "etc.", "e.g.", "i.e.",
	"U.S.", "U.K.",
}

func (s *SentenceSplitter) splitEnglish(text string) []string {
	var sentences []string
	var buf strings.Builder

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		buf.WriteRune(r)

		if r == '\n' {
			sentence := strings.TrimSpace(buf.String())
			if sentence != "" {
				sentences = append(sentences, sentence)
			}
			buf.Reset()
			continue
		}

		if r == '.' || r == '!' || r == '?' {
			// Check if next char is a space or end of text (sentence boundary)
			isEnd := i == len(runes)-1
			hasSpace := !isEnd && (runes[i+1] == ' ' || runes[i+1] == '\n' || runes[i+1] == '\r')

			if isEnd || hasSpace {
				if r == '.' && s.isAbbreviation(buf.String()) {
					continue
				}
				sentence := strings.TrimSpace(buf.String())
				if sentence != "" {
					sentences = append(sentences, sentence)
				}
				buf.Reset()
			}
		}
	}

	if remaining := strings.TrimSpace(buf.String()); remaining != "" {
		sentences = append(sentences, remaining)
	}

	return sentences
}

func (s *SentenceSplitter) isAbbreviation(text string) bool {
	for _, abbr := range englishAbbreviations {
		if strings.HasSuffix(text, abbr) {
			// Ensure the abbreviation is at a word boundary
			prefix := text[:len(text)-len(abbr)]
			if prefix == "" {
				return true
			}
			lastRune, _ := utf8.DecodeLastRuneInString(prefix)
			if lastRune == ' ' || lastRune == '\n' || lastRune == '\t' {
				return true
			}
		}
	}
	return false
}
