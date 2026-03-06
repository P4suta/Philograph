package service

import (
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphBuilder_Build(t *testing.T) {
	terms := []*model.Term{
		{ID: 0, Text: "cat", Frequency: 10},
		{ID: 1, Text: "dog", Frequency: 8},
		{ID: 2, Text: "bird", Frequency: 5},
	}

	pairs := []*model.CooccurrencePair{
		{TermAID: 0, TermBID: 1, RawCount: 5, NPMI: 0.8},
		{TermAID: 0, TermBID: 2, RawCount: 3, NPMI: 0.5},
		{TermAID: 1, TermBID: 2, RawCount: 2, NPMI: 0.3},
	}

	builder := NewGraphBuilder(150, nil) // no analyzer
	graph, err := builder.Build(pairs, terms)
	require.NoError(t, err)

	assert.Equal(t, 3, graph.NodeCount())
	assert.Equal(t, 3, graph.EdgeCount())
}

func TestGraphBuilder_MaxNodesCutoff(t *testing.T) {
	terms := []*model.Term{
		{ID: 0, Text: "a", Frequency: 100},
		{ID: 1, Text: "b", Frequency: 50},
		{ID: 2, Text: "c", Frequency: 10},
	}

	pairs := []*model.CooccurrencePair{
		{TermAID: 0, TermBID: 1, RawCount: 5, NPMI: 0.8},
		{TermAID: 0, TermBID: 2, RawCount: 3, NPMI: 0.5},
		{TermAID: 1, TermBID: 2, RawCount: 2, NPMI: 0.3},
	}

	builder := NewGraphBuilder(2, nil) // max 2 nodes
	graph, err := builder.Build(pairs, terms)
	require.NoError(t, err)

	assert.Equal(t, 2, graph.NodeCount())
	// Only edges between kept nodes
	assert.Equal(t, 1, graph.EdgeCount())
}
