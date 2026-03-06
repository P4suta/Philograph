package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"Philograph/internal/api"
	"Philograph/internal/application"
	"Philograph/internal/domain/model"
	"Philograph/internal/infrastructure/export"
	"Philograph/internal/infrastructure/graphanalyzer"
	kagometok "Philograph/internal/infrastructure/kagome"
	"Philograph/internal/infrastructure/whitespace"
	"Philograph/internal/port"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_EnglishPipeline(t *testing.T) {
	data, err := os.ReadFile("testdata/english_sample.txt")
	require.NoError(t, err)

	tok := whitespace.NewTokenizer()
	analyzer := graphanalyzer.NewAnalyzer()
	pipeline := application.NewPipeline(tok, analyzer, nil)

	config := model.DefaultConfig()
	config.Language = model.LangEnglish
	config.MinFrequency = 2
	config.MinCooccurrence = 1

	result, err := pipeline.Run(context.Background(), string(data), config)
	require.NoError(t, err)

	assert.True(t, result.Graph.NodeCount() > 0, "should have nodes")
	assert.True(t, result.Graph.EdgeCount() > 0, "should have edges")

	for _, n := range result.Graph.Nodes {
		assert.True(t, n.DegreeCentrality >= 0, "node %s should have non-negative degree", n.Label)
	}

	var buf bytes.Buffer
	exporter := export.NewJSONExporter()
	err = exporter.Export(&buf, result.Graph)
	require.NoError(t, err)

	var exported map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &exported)
	require.NoError(t, err)
	assert.NotEmpty(t, exported["nodes"])
}

func TestIntegration_JapanesePipeline(t *testing.T) {
	data, err := os.ReadFile("testdata/japanese_sample.txt")
	require.NoError(t, err)

	tok, err := kagometok.NewTokenizer()
	require.NoError(t, err)

	analyzer := graphanalyzer.NewAnalyzer()
	pipeline := application.NewPipeline(tok, analyzer, nil)

	config := model.DefaultConfig()
	config.Language = model.LangJapanese
	config.MinFrequency = 2
	config.MinCooccurrence = 1

	result, err := pipeline.Run(context.Background(), string(data), config)
	require.NoError(t, err)

	assert.True(t, result.Graph.NodeCount() > 0, "should have nodes")
	assert.True(t, result.Graph.EdgeCount() > 0, "should have edges")
}

func TestIntegration_APIServer(t *testing.T) {
	tok := whitespace.NewTokenizer()
	analyzer := graphanalyzer.NewAnalyzer()
	wsHub := api.NewWSHub()
	pipeline := application.NewPipeline(tok, analyzer, wsHub.ProgressListener())

	config := model.DefaultConfig()
	config.Language = model.LangEnglish
	config.MinFrequency = 2
	config.MinCooccurrence = 1

	session := application.NewSession(pipeline, config)
	exporters := map[string]port.Exporter{
		"json": export.NewJSONExporter(),
		"gexf": export.NewGEXFExporter(),
	}

	handler := api.NewHandler(session, exporters)
	server := api.NewServer(handler, wsHub, 0)

	serverPort, err := server.Start()
	require.NoError(t, err)
	defer server.Shutdown(context.Background())

	base := fmt.Sprintf("http://localhost:%d", serverPort)

	// Health check
	resp, err := http.Get(base + "/api/v1/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Upload and analyze
	data, _ := os.ReadFile("testdata/english_sample.txt")
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write(data)
	writer.Close()

	resp, err = http.Post(base+"/api/v1/analyze", writer.FormDataContentType(), &buf)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Get result
	resp, err = http.Get(base + "/api/v1/result")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Export JSON
	resp, err = http.Get(base + "/api/v1/export?format=json")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Export GEXF
	resp, err = http.Get(base + "/api/v1/export?format=gexf")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}
