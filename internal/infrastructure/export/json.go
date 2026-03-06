package export

import (
	"encoding/json"
	"io"

	"Philograph/internal/domain/model"
)

type jsonGraph struct {
	Nodes []jsonNode `json:"nodes"`
	Edges []jsonEdge `json:"edges"`
}

type jsonNode struct {
	ID                    int     `json:"id"`
	Label                 string  `json:"label"`
	Frequency             int     `json:"frequency"`
	DegreeCentrality      float64 `json:"degree_centrality"`
	BetweennessCentrality float64 `json:"betweenness_centrality"`
	EigenvectorCentrality float64 `json:"eigenvector_centrality"`
	CommunityID           int     `json:"community_id"`
}

type jsonEdge struct {
	SourceID int     `json:"source_id"`
	TargetID int     `json:"target_id"`
	Weight   float64 `json:"weight"`
	RawCount int     `json:"raw_count"`
}

// JSONExporter はJSON形式でグラフデータをエクスポートする。
type JSONExporter struct{}

// NewJSONExporter は新しいJSONExporterを返す。
func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

// Export はグラフをJSON形式で書き出す。
func (e *JSONExporter) Export(w io.Writer, g *model.Graph) error {
	jg := jsonGraph{
		Nodes: make([]jsonNode, len(g.Nodes)),
		Edges: make([]jsonEdge, len(g.Edges)),
	}

	for i, n := range g.Nodes {
		jg.Nodes[i] = jsonNode{
			ID:                    n.ID,
			Label:                 n.Label,
			Frequency:             n.Frequency,
			DegreeCentrality:      n.DegreeCentrality,
			BetweennessCentrality: n.BetweennessCentrality,
			EigenvectorCentrality: n.EigenvectorCentrality,
			CommunityID:           n.CommunityID,
		}
	}

	for i, edge := range g.Edges {
		jg.Edges[i] = jsonEdge{
			SourceID: edge.SourceID,
			TargetID: edge.TargetID,
			Weight:   edge.Weight,
			RawCount: edge.RawCount,
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jg)
}

// ContentType はMIMEタイプを返す。
func (e *JSONExporter) ContentType() string {
	return "application/json"
}

// FileExtension はファイル拡張子を返す。
func (e *JSONExporter) FileExtension() string {
	return ".json"
}
