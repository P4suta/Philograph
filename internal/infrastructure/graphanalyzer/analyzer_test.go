package graphanalyzer

import (
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTriangleGraph() *model.Graph {
	return &model.Graph{
		Nodes: []model.Node{
			{ID: 0, Label: "a"},
			{ID: 1, Label: "b"},
			{ID: 2, Label: "c"},
		},
		Edges: []model.Edge{
			{SourceID: 0, TargetID: 1, Weight: 1.0},
			{SourceID: 1, TargetID: 2, Weight: 1.0},
			{SourceID: 0, TargetID: 2, Weight: 1.0},
		},
	}
}

func TestAnalyzer_ComputeCentrality(t *testing.T) {
	analyzer := NewAnalyzer()
	g := makeTriangleGraph()

	err := analyzer.ComputeCentrality(g)
	require.NoError(t, err)

	// In a triangle, all nodes have equal degree centrality = 2/(3-1) = 1.0
	for _, n := range g.Nodes {
		assert.InDelta(t, 1.0, n.DegreeCentrality, 0.001, "node %s degree", n.Label)
	}

	// Betweenness should be 0 for all nodes in a triangle (no shortest paths through intermediaries)
	for _, n := range g.Nodes {
		assert.InDelta(t, 0.0, n.BetweennessCentrality, 0.001, "node %s betweenness", n.Label)
	}

	// Eigenvector/PageRank should be roughly equal for all nodes
	for _, n := range g.Nodes {
		assert.True(t, n.EigenvectorCentrality > 0, "node %s should have positive eigenvector centrality", n.Label)
	}
}

func TestAnalyzer_ComputeCentrality_Star(t *testing.T) {
	// Star graph: center node connected to 3 leaves
	g := &model.Graph{
		Nodes: []model.Node{
			{ID: 0, Label: "center"},
			{ID: 1, Label: "leaf1"},
			{ID: 2, Label: "leaf2"},
			{ID: 3, Label: "leaf3"},
		},
		Edges: []model.Edge{
			{SourceID: 0, TargetID: 1, Weight: 1.0},
			{SourceID: 0, TargetID: 2, Weight: 1.0},
			{SourceID: 0, TargetID: 3, Weight: 1.0},
		},
	}

	analyzer := NewAnalyzer()
	err := analyzer.ComputeCentrality(g)
	require.NoError(t, err)

	// Center should have highest degree centrality
	assert.True(t, g.Nodes[0].DegreeCentrality > g.Nodes[1].DegreeCentrality)
	// Center should have highest betweenness
	assert.True(t, g.Nodes[0].BetweennessCentrality > g.Nodes[1].BetweennessCentrality)
}

func TestAnalyzer_DetectCommunities(t *testing.T) {
	analyzer := NewAnalyzer()
	g := makeTriangleGraph()

	err := analyzer.DetectCommunities(g)
	require.NoError(t, err)

	// All nodes in a triangle should be in the same community
	assert.Equal(t, g.Nodes[0].CommunityID, g.Nodes[1].CommunityID)
	assert.Equal(t, g.Nodes[1].CommunityID, g.Nodes[2].CommunityID)
}

func TestAnalyzer_EmptyGraph(t *testing.T) {
	analyzer := NewAnalyzer()
	g := &model.Graph{}

	assert.NoError(t, analyzer.ComputeCentrality(g))
	assert.NoError(t, analyzer.DetectCommunities(g))
}
