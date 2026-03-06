package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleGraph() *model.Graph {
	return &model.Graph{
		Nodes: []model.Node{
			{ID: 0, Label: "cat", Frequency: 10, DegreeCentrality: 0.5, CommunityID: 0},
			{ID: 1, Label: "dog", Frequency: 8, DegreeCentrality: 0.3, CommunityID: 1},
		},
		Edges: []model.Edge{
			{SourceID: 0, TargetID: 1, Weight: 0.75, RawCount: 5},
		},
	}
}

func TestJSONExporter_Export(t *testing.T) {
	exporter := NewJSONExporter()
	var buf bytes.Buffer

	err := exporter.Export(&buf, sampleGraph())
	require.NoError(t, err)

	// Verify valid JSON
	var result map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	nodes := result["nodes"].([]interface{})
	assert.Len(t, nodes, 2)

	edges := result["edges"].([]interface{})
	assert.Len(t, edges, 1)

	// Check snake_case keys
	node := nodes[0].(map[string]interface{})
	_, ok := node["degree_centrality"]
	assert.True(t, ok, "should use snake_case keys")
}

func TestJSONExporter_ContentType(t *testing.T) {
	e := NewJSONExporter()
	assert.Equal(t, "application/json", e.ContentType())
	assert.Equal(t, ".json", e.FileExtension())
}

func TestGEXFExporter_Export(t *testing.T) {
	exporter := NewGEXFExporter()
	var buf bytes.Buffer

	err := exporter.Export(&buf, sampleGraph())
	require.NoError(t, err)

	output := buf.String()
	assert.True(t, strings.Contains(output, "<?xml"))
	assert.True(t, strings.Contains(output, "<gexf"))
	assert.True(t, strings.Contains(output, "cat"))
	assert.True(t, strings.Contains(output, "dog"))
}

func TestGEXFExporter_ContentType(t *testing.T) {
	e := NewGEXFExporter()
	assert.Equal(t, "application/xml", e.ContentType())
	assert.Equal(t, ".gexf", e.FileExtension())
}
