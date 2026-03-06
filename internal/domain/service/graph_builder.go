package service

import (
	"Philograph/internal/domain/model"
)

// GraphAnalyzer はグラフのメトリクス計算を行うポートインターフェース。
type GraphAnalyzer interface {
	ComputeCentrality(graph *model.Graph) error
	DetectCommunities(graph *model.Graph) error
}

type nodeCandidate struct {
	id   int
	freq int
}

// GraphBuilder は共起ペアからグラフを構築する。
type GraphBuilder struct {
	maxNodes int
	analyzer GraphAnalyzer
}

// NewGraphBuilder は新しいGraphBuilderを返す。
func NewGraphBuilder(maxNodes int, analyzer GraphAnalyzer) *GraphBuilder {
	if maxNodes <= 0 {
		maxNodes = 150
	}
	return &GraphBuilder{
		maxNodes: maxNodes,
		analyzer: analyzer,
	}
}

// Build は共起ペアと語彙からグラフを構築する。
func (b *GraphBuilder) Build(pairs []*model.CooccurrencePair, terms []*model.Term) (*model.Graph, error) {
	// Collect node IDs from pairs
	nodeIDs := make(map[int]bool)
	for _, p := range pairs {
		nodeIDs[p.TermAID] = true
		nodeIDs[p.TermBID] = true
	}

	// Build term lookup
	termByID := make(map[int]*model.Term, len(terms))
	for _, t := range terms {
		termByID[t.ID] = t
	}

	// Create nodes (apply max nodes cutoff based on frequency)
	var candidates []nodeCandidate
	for id := range nodeIDs {
		freq := 0
		if t, ok := termByID[id]; ok {
			freq = t.Frequency
		}
		candidates = append(candidates, nodeCandidate{id: id, freq: freq})
	}

	// Sort by frequency descending for cutoff
	sortNodeCandidates(candidates)

	if len(candidates) > b.maxNodes {
		candidates = candidates[:b.maxNodes]
	}

	// Create the kept-nodes set
	keptNodes := make(map[int]bool, len(candidates))
	for _, c := range candidates {
		keptNodes[c.id] = true
	}

	// Build nodes
	nodes := make([]model.Node, 0, len(candidates))
	for _, c := range candidates {
		t := termByID[c.id]
		label := ""
		freq := 0
		if t != nil {
			label = t.Text
			freq = t.Frequency
		}
		nodes = append(nodes, model.Node{
			ID:        c.id,
			Label:     label,
			Frequency: freq,
		})
	}

	// Build edges (only between kept nodes)
	var edges []model.Edge
	for _, p := range pairs {
		if !keptNodes[p.TermAID] || !keptNodes[p.TermBID] {
			continue
		}
		weight := metricWeight(p)
		edges = append(edges, model.Edge{
			SourceID: p.TermAID,
			TargetID: p.TermBID,
			Weight:   weight,
			RawCount: p.RawCount,
		})
	}

	graph := &model.Graph{
		Nodes: nodes,
		Edges: edges,
	}

	// Compute centrality and communities if analyzer is available
	if b.analyzer != nil {
		if err := b.analyzer.ComputeCentrality(graph); err != nil {
			return nil, err
		}
		if err := b.analyzer.DetectCommunities(graph); err != nil {
			return nil, err
		}
	}

	return graph, nil
}

func metricWeight(p *model.CooccurrencePair) float64 {
	if p.NPMI > 0 {
		return p.NPMI
	}
	return float64(p.RawCount)
}

func sortNodeCandidates(candidates []nodeCandidate) {
	for i := 1; i < len(candidates); i++ {
		for j := i; j > 0 && candidates[j].freq > candidates[j-1].freq; j-- {
			candidates[j], candidates[j-1] = candidates[j-1], candidates[j]
		}
	}
}
