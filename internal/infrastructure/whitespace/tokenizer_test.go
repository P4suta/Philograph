package whitespace

import (
	"context"
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizer_Tokenize(t *testing.T) {
	tok := NewTokenizer()

	tokens, err := tok.Tokenize(context.Background(), "The cat sat on the mat.")
	require.NoError(t, err)

	assert.Len(t, tokens, 6)
	assert.Equal(t, "The", tokens[0].Surface)
	assert.Equal(t, "the", tokens[0].Base)
	assert.Equal(t, model.POSNoun, tokens[0].POS)

	// "mat." should have punctuation stripped
	assert.Equal(t, "mat", tokens[5].Surface)
	assert.Equal(t, "mat", tokens[5].Base)
}

func TestTokenizer_Language(t *testing.T) {
	tok := NewTokenizer()
	assert.Equal(t, model.LangEnglish, tok.Language())
}

func TestTokenizer_EmptyInput(t *testing.T) {
	tok := NewTokenizer()
	tokens, err := tok.Tokenize(context.Background(), "")
	require.NoError(t, err)
	assert.Empty(t, tokens)
}
