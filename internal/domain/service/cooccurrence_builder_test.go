package service

import (
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCooccurrenceBuilder_Build(t *testing.T) {
	tokens := []model.Token{
		{Surface: "cat", Base: "cat", POS: model.POSNoun, SentenceIndex: 0, PositionInSentence: 0},
		{Surface: "sat", Base: "sit", POS: model.POSVerb, SentenceIndex: 0, PositionInSentence: 1},
		{Surface: "mat", Base: "mat", POS: model.POSNoun, SentenceIndex: 0, PositionInSentence: 2},
	}

	termMap, _ := BuildVocabulary(tokens)
	builder := NewCooccurrenceBuilder(5)
	pairs := builder.Build(tokens, termMap)

	// 3 tokens → 3 pairs: cat-sit, cat-mat, sit-mat
	assert.Len(t, pairs, 3)

	for _, p := range pairs {
		assert.Equal(t, 1, p.RawCount)
		assert.True(t, p.TermAID < p.TermBID, "pair should be normalized")
	}
}

func TestCooccurrenceBuilder_SentenceBoundary(t *testing.T) {
	tokens := []model.Token{
		{Surface: "cat", Base: "cat", POS: model.POSNoun, SentenceIndex: 0, PositionInSentence: 0},
		{Surface: "dog", Base: "dog", POS: model.POSNoun, SentenceIndex: 1, PositionInSentence: 0},
	}

	termMap, _ := BuildVocabulary(tokens)
	builder := NewCooccurrenceBuilder(5)
	pairs := builder.Build(tokens, termMap)

	// Different sentences → no co-occurrence
	assert.Len(t, pairs, 0)
}

func TestBuildVocabulary(t *testing.T) {
	tokens := []model.Token{
		{Surface: "cat", Base: "cat", POS: model.POSNoun},
		{Surface: "cats", Base: "cat", POS: model.POSNoun},
		{Surface: "dog", Base: "dog", POS: model.POSNoun},
	}

	termMap, terms := BuildVocabulary(tokens)

	require.Len(t, terms, 2)
	assert.Equal(t, 2, termMap["cat"].Frequency)
	assert.Equal(t, 1, termMap["dog"].Frequency)
}

func TestCooccurrenceBuilder_WindowCount(t *testing.T) {
	builder := NewCooccurrenceBuilder(5)
	assert.Equal(t, 0, builder.WindowCount(0))
	assert.Equal(t, 0, builder.WindowCount(1))
	assert.Equal(t, 1, builder.WindowCount(3))
	assert.Equal(t, 6, builder.WindowCount(10))
}
