package autotokenizer

import (
	"context"
	"sync"

	"Philograph/internal/domain/model"
	"Philograph/internal/port"
)

// LanguageAware はTokenizerに対して言語を動的に設定可能にするインターフェース。
type LanguageAware interface {
	SetLanguage(lang model.Language)
}

// AutoTokenizer はテキストの言語に応じて適切なTokenizerに委譲するラッパー。
type AutoTokenizer struct {
	mu       sync.RWMutex
	lang     model.Language
	japanese port.Tokenizer
	english  port.Tokenizer
}

// NewAutoTokenizer は新しいAutoTokenizerを返す。
func NewAutoTokenizer(japanese port.Tokenizer, english port.Tokenizer) *AutoTokenizer {
	return &AutoTokenizer{
		japanese: japanese,
		english:  english,
	}
}

// SetLanguage は使用する言語を設定する。
func (a *AutoTokenizer) SetLanguage(lang model.Language) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lang = lang
}

// Language は現在設定されている言語を返す。
func (a *AutoTokenizer) Language() model.Language {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.lang == model.LangUnknown {
		return model.LangJapanese
	}
	return a.lang
}

// Tokenize はSetLanguageで設定された言語のTokenizerに委譲する。
func (a *AutoTokenizer) Tokenize(ctx context.Context, sentence string) ([]model.Token, error) {
	a.mu.RLock()
	lang := a.lang
	a.mu.RUnlock()

	if lang == model.LangJapanese {
		return a.japanese.Tokenize(ctx, sentence)
	}
	return a.english.Tokenize(ctx, sentence)
}
