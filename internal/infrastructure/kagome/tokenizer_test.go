package kagome

import (
	"context"
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizer_Tokenize(t *testing.T) {
	tok, err := NewTokenizer()
	require.NoError(t, err)

	tokens, err := tok.Tokenize(context.Background(), "吾輩は猫である")
	require.NoError(t, err)

	assert.True(t, len(tokens) > 0)

	// 「猫」should be a noun
	var foundCat bool
	for _, tk := range tokens {
		if tk.Surface == "猫" {
			assert.Equal(t, model.POSNoun, tk.POS)
			foundCat = true
		}
	}
	assert.True(t, foundCat, "should find 猫 token")
}

func TestTokenizer_BaseForm(t *testing.T) {
	tok, err := NewTokenizer()
	require.NoError(t, err)

	tokens, err := tok.Tokenize(context.Background(), "走った")
	require.NoError(t, err)

	// 「走った」→ base form should be 「走る」
	var found bool
	for _, tk := range tokens {
		if tk.Surface == "走っ" || tk.Surface == "走った" {
			if tk.Base == "走る" {
				found = true
			}
		}
	}
	assert.True(t, found, "should find base form 走る")
}

func TestTokenizer_Language(t *testing.T) {
	tok, err := NewTokenizer()
	require.NoError(t, err)
	assert.Equal(t, model.LangJapanese, tok.Language())
}
