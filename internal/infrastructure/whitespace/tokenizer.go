package whitespace

import (
	"context"
	"strings"
	"unicode"

	"Philograph/internal/domain/model"
)

// Tokenizer は英語テキストをホワイトスペースで分割するトークナイザー。
type Tokenizer struct{}

// NewTokenizer は新しいTokenizerを返す。
func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

// Language はこのトークナイザーの対象言語を返す。
func (t *Tokenizer) Language() model.Language {
	return model.LangEnglish
}

// Tokenize はテキストをトークン化する。
func (t *Tokenizer) Tokenize(_ context.Context, sentence string) ([]model.Token, error) {
	fields := strings.Fields(sentence)
	tokens := make([]model.Token, 0, len(fields))

	for i, f := range fields {
		word := stripPunctuation(f)
		if word == "" {
			continue
		}
		lower := strings.ToLower(word)
		tokens = append(tokens, model.Token{
			Surface:            word,
			Base:               lower,
			POS:                model.POSNoun,
			PositionInSentence: i,
		})
	}

	return tokens, nil
}

func stripPunctuation(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSymbol(r)
	})
}
