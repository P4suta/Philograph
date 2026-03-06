package export

import (
	"encoding/xml"
	"fmt"
	"io"

	"Philograph/internal/domain/model"
)

// GEXF XML structure
type gexfDoc struct {
	XMLName xml.Name  `xml:"gexf"`
	Xmlns   string    `xml:"xmlns,attr"`
	Version string    `xml:"version,attr"`
	Graph   gexfGraph `xml:"graph"`
}

type gexfGraph struct {
	DefaultEdgeType string         `xml:"defaultedgetype,attr"`
	Attributes      gexfAttributes `xml:"attributes"`
	Nodes           gexfNodes      `xml:"nodes"`
	Edges           gexfEdges      `xml:"edges"`
}

type gexfAttributes struct {
	Class string          `xml:"class,attr"`
	Attrs []gexfAttribute `xml:"attribute"`
}

type gexfAttribute struct {
	ID    string `xml:"id,attr"`
	Title string `xml:"title,attr"`
	Type  string `xml:"type,attr"`
}

type gexfNodes struct {
	Nodes []gexfNode `xml:"node"`
}

type gexfNode struct {
	ID        string         `xml:"id,attr"`
	Label     string         `xml:"label,attr"`
	AttValues gexfAttValues  `xml:"attvalues"`
}

type gexfEdges struct {
	Edges []gexfEdge `xml:"edge"`
}

type gexfEdge struct {
	ID     string  `xml:"id,attr"`
	Source string  `xml:"source,attr"`
	Target string  `xml:"target,attr"`
	Weight float64 `xml:"weight,attr"`
}

type gexfAttValues struct {
	Values []gexfAttValue `xml:"attvalue"`
}

type gexfAttValue struct {
	For   string `xml:"for,attr"`
	Value string `xml:"value,attr"`
}

// GEXFExporter はGEXF形式でグラフデータをエクスポートする。
type GEXFExporter struct{}

// NewGEXFExporter は新しいGEXFExporterを返す。
func NewGEXFExporter() *GEXFExporter {
	return &GEXFExporter{}
}

// Export はグラフをGEXF形式で書き出す。
func (e *GEXFExporter) Export(w io.Writer, g *model.Graph) error {
	doc := gexfDoc{
		Xmlns:   "http://www.gexf.net/1.2draft",
		Version: "1.2",
		Graph: gexfGraph{
			DefaultEdgeType: "undirected",
			Attributes: gexfAttributes{
				Class: "node",
				Attrs: []gexfAttribute{
					{ID: "0", Title: "frequency", Type: "integer"},
					{ID: "1", Title: "degree_centrality", Type: "float"},
					{ID: "2", Title: "betweenness_centrality", Type: "float"},
					{ID: "3", Title: "eigenvector_centrality", Type: "float"},
					{ID: "4", Title: "community_id", Type: "integer"},
				},
			},
		},
	}

	nodes := make([]gexfNode, len(g.Nodes))
	for i, n := range g.Nodes {
		nodes[i] = gexfNode{
			ID:    fmt.Sprintf("%d", n.ID),
			Label: n.Label,
			AttValues: gexfAttValues{
				Values: []gexfAttValue{
					{For: "0", Value: fmt.Sprintf("%d", n.Frequency)},
					{For: "1", Value: fmt.Sprintf("%.6f", n.DegreeCentrality)},
					{For: "2", Value: fmt.Sprintf("%.6f", n.BetweennessCentrality)},
					{For: "3", Value: fmt.Sprintf("%.6f", n.EigenvectorCentrality)},
					{For: "4", Value: fmt.Sprintf("%d", n.CommunityID)},
				},
			},
		}
	}
	doc.Graph.Nodes = gexfNodes{Nodes: nodes}

	edges := make([]gexfEdge, len(g.Edges))
	for i, edge := range g.Edges {
		edges[i] = gexfEdge{
			ID:     fmt.Sprintf("%d", i),
			Source: fmt.Sprintf("%d", edge.SourceID),
			Target: fmt.Sprintf("%d", edge.TargetID),
			Weight: edge.Weight,
		}
	}
	doc.Graph.Edges = gexfEdges{Edges: edges}

	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(doc)
}

// ContentType はMIMEタイプを返す。
func (e *GEXFExporter) ContentType() string {
	return "application/xml"
}

// FileExtension はファイル拡張子を返す。
func (e *GEXFExporter) FileExtension() string {
	return ".gexf"
}
