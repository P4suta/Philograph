package model

// Graph は共起ネットワークを表す。
type Graph struct {
	Nodes []Node
	Edges []Edge
}

// NodeCount はノード数を返す。
func (g Graph) NodeCount() int { return len(g.Nodes) }

// EdgeCount はエッジ数を返す。
func (g Graph) EdgeCount() int { return len(g.Edges) }

// Node はネットワーク上の1語彙を表す。
type Node struct {
	ID                    int
	Label                 string
	Frequency             int
	DegreeCentrality      float64
	BetweennessCentrality float64
	EigenvectorCentrality float64
	CommunityID           int
}

// Edge はネットワーク上の共起関係を表す。
type Edge struct {
	SourceID int
	TargetID int
	Weight   float64
	RawCount int
}
