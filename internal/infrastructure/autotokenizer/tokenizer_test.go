package autotokenizer

import (
	"context"
	"testing"

	"Philograph/internal/domain/model"
	"Philograph/internal/infrastructure/whitespace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubTokenizer はテスト用のスタブトークナイザー。
type stubTokenizer struct {
	lang   model.Language
	called bool
}

func (s *stubTokenizer) Tokenize(_ context.Context, sentence string) ([]model.Token, error) {
	s.called = true
	return []model.Token{{Surface: sentence, Base: sentence, POS: model.POSNoun}}, nil
}

func (s *stubTokenizer) Language() model.Language {
	return s.lang
}

func TestAutoTokenizer_DelegatesToJapanese(t *testing.T) {
	ja := &stubTokenizer{lang: model.LangJapanese}
	en := &stubTokenizer{lang: model.LangEnglish}

	auto := NewAutoTokenizer(ja, en)
	auto.SetLanguage(model.LangJapanese)

	_, err := auto.Tokenize(context.Background(), "テスト")
	require.NoError(t, err)
	assert.True(t, ja.called)
	assert.False(t, en.called)
}

func TestAutoTokenizer_DelegatesToEnglish(t *testing.T) {
	ja := &stubTokenizer{lang: model.LangJapanese}
	en := &stubTokenizer{lang: model.LangEnglish}

	auto := NewAutoTokenizer(ja, en)
	auto.SetLanguage(model.LangEnglish)

	_, err := auto.Tokenize(context.Background(), "hello world")
	require.NoError(t, err)
	assert.False(t, ja.called)
	assert.True(t, en.called)
}

func TestAutoTokenizer_Language(t *testing.T) {
	en := whitespace.NewTokenizer()
	auto := NewAutoTokenizer(nil, en)

	// Default (LangUnknown) returns LangJapanese as fallback
	assert.Equal(t, model.LangJapanese, auto.Language())

	auto.SetLanguage(model.LangEnglish)
	assert.Equal(t, model.LangEnglish, auto.Language())
}

func TestAutoTokenizer_ImplementsLanguageAware(t *testing.T) {
	auto := NewAutoTokenizer(nil, nil)
	var _ LanguageAware = auto
}
