package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken_IsContent(t *testing.T) {
	tests := []struct {
		pos      POS
		expected bool
	}{
		{POSNoun, true},
		{POSVerb, true},
		{POSAdjective, true},
		{POSAdverb, true},
		{POSParticle, false},
		{POSAuxVerb, false},
		{POSSymbol, false},
		{POSOther, false},
	}

	for _, tt := range tests {
		tok := Token{POS: tt.pos}
		assert.Equal(t, tt.expected, tok.IsContent(), "POS: %s", tt.pos)
	}
}

func TestNormalizePairOrder(t *testing.T) {
	a, b := NormalizePairOrder(5, 3)
	assert.Equal(t, 3, a)
	assert.Equal(t, 5, b)

	a, b = NormalizePairOrder(2, 7)
	assert.Equal(t, 2, a)
	assert.Equal(t, 7, b)
}

func TestCooccurrencePair_Key(t *testing.T) {
	p := CooccurrencePair{TermAID: 3, TermBID: 7}
	assert.Equal(t, "3:7", p.Key())
}

func TestDetectLanguage(t *testing.T) {
	assert.Equal(t, LangJapanese, DetectLanguage("吾輩は猫である。名前はまだ無い。"))
	assert.Equal(t, LangEnglish, DetectLanguage("The quick brown fox jumps over the lazy dog."))
	assert.Equal(t, LangUnknown, DetectLanguage(""))
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, 5, cfg.WindowSize)
	assert.Equal(t, 3, cfg.MinFrequency)
	assert.Equal(t, 150, cfg.MaxNodes)
	assert.Equal(t, MetricNPMI, cfg.Metric)
}

func TestGraph_Counts(t *testing.T) {
	g := Graph{
		Nodes: []Node{{ID: 0}, {ID: 1}},
		Edges: []Edge{{SourceID: 0, TargetID: 1}},
	}
	assert.Equal(t, 2, g.NodeCount())
	assert.Equal(t, 1, g.EdgeCount())
}
