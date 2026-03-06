package graphanalyzer

import (
	"Philograph/internal/domain/model"

	"gonum.org/v1/gonum/graph/community"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/simple"
)

// Analyzer はgonumを使ったグラフ分析器。
type Analyzer struct{}

// NewAnalyzer は新しいAnalyzerを返す。
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// ComputeCentrality は各ノードの中心性指標を計算する。
func (a *Analyzer) ComputeCentrality(g *model.Graph) error {
	if len(g.Nodes) == 0 {
		return nil
	}

	ug := simple.NewUndirectedGraph()
	nodeIDMap := make(map[int]int64) // model ID → gonum ID
	gonumToModel := make(map[int64]int)

	for i, n := range g.Nodes {
		id := int64(i)
		nodeIDMap[n.ID] = id
		gonumToModel[id] = i
		ug.AddNode(simple.Node(id))
	}

	for _, e := range g.Edges {
		fromID, okFrom := nodeIDMap[e.SourceID]
		toID, okTo := nodeIDMap[e.TargetID]
		if !okFrom || !okTo {
			continue
		}
		if ug.HasEdgeBetween(fromID, toID) {
			continue
		}
		ug.SetEdge(simple.Edge{F: simple.Node(fromID), T: simple.Node(toID)})
	}

	// Degree centrality
	maxDegree := float64(len(g.Nodes) - 1)
	if maxDegree < 1 {
		maxDegree = 1
	}
	for id, idx := range gonumToModel {
		degree := float64(ug.From(id).Len())
		g.Nodes[idx].DegreeCentrality = degree / maxDegree
	}

	// Betweenness centrality
	betweenness := network.Betweenness(ug)
	normFactor := 1.0
	n := len(g.Nodes)
	if n > 2 {
		normFactor = 2.0 / float64((n-1)*(n-2))
	}
	for id, val := range betweenness {
		if idx, ok := gonumToModel[id]; ok {
			g.Nodes[idx].BetweennessCentrality = val * normFactor
		}
	}

	// Eigenvector centrality (PageRank as approximation for undirected)
	// Use a simple power iteration approach via gonum's PageRank
	// For undirected graphs, PageRank approximates eigenvector centrality
	dg := simple.NewDirectedGraph()
	for i := range g.Nodes {
		dg.AddNode(simple.Node(int64(i)))
	}
	for _, e := range g.Edges {
		fromID, okFrom := nodeIDMap[e.SourceID]
		toID, okTo := nodeIDMap[e.TargetID]
		if !okFrom || !okTo {
			continue
		}
		// Add both directions for undirected
		dg.SetEdge(simple.Edge{F: simple.Node(fromID), T: simple.Node(toID)})
		dg.SetEdge(simple.Edge{F: simple.Node(toID), T: simple.Node(fromID)})
	}

	pagerank := network.PageRank(dg, 0.85, 1e-6)
	for id, val := range pagerank {
		if idx, ok := gonumToModel[id]; ok {
			g.Nodes[idx].EigenvectorCentrality = val
		}
	}

	return nil
}

// DetectCommunities はLouvain法でコミュニティを検出する。
func (a *Analyzer) DetectCommunities(g *model.Graph) error {
	if len(g.Nodes) == 0 {
		return nil
	}

	ug := simple.NewUndirectedGraph()
	nodeIDMap := make(map[int]int64)
	gonumToModel := make(map[int64]int)

	for i, n := range g.Nodes {
		id := int64(i)
		nodeIDMap[n.ID] = id
		gonumToModel[id] = i
		ug.AddNode(simple.Node(id))
	}

	for _, e := range g.Edges {
		fromID, okFrom := nodeIDMap[e.SourceID]
		toID, okTo := nodeIDMap[e.TargetID]
		if !okFrom || !okTo {
			continue
		}
		if ug.HasEdgeBetween(fromID, toID) {
			continue
		}
		ug.SetEdge(simple.Edge{F: simple.Node(fromID), T: simple.Node(toID)})
	}

	// Louvain community detection
	reduced := community.Modularize(ug, 1, nil)
	communities := reduced.Communities()

	for commID, comm := range communities {
		for _, n := range comm {
			if idx, ok := gonumToModel[n.ID()]; ok {
				g.Nodes[idx].CommunityID = commID
			}
		}
	}

	return nil
}
