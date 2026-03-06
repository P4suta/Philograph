package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Philograph/internal/application"
	"Philograph/internal/domain/model"
	"Philograph/internal/infrastructure/export"
	"Philograph/internal/infrastructure/whitespace"
	"Philograph/internal/port"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestHandler() *Handler {
	tok := whitespace.NewTokenizer()
	pipeline := application.NewPipeline(tok, nil, nil)
	config := model.AnalysisConfig{
		WindowSize:      5,
		MinFrequency:    2,
		MinCooccurrence: 1,
		TargetPOS:       []model.POS{model.POSNoun},
		Metric:          model.MetricNPMI,
		MaxNodes:        150,
		Language:        model.LangEnglish,
	}
	session := application.NewSession(pipeline, config)

	exporters := map[string]port.Exporter{
		"json": export.NewJSONExporter(),
		"gexf": export.NewGEXFExporter(),
	}

	return NewHandler(session, exporters)
}

func TestHandleHealth(t *testing.T) {
	h := setupTestHandler()
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	h.HandleHealth(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "ok", result["status"])
}

func TestHandleGetConfig(t *testing.T) {
	h := setupTestHandler()
	req := httptest.NewRequest("GET", "/api/v1/config", nil)
	w := httptest.NewRecorder()

	h.HandleGetConfig(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleGetResult_NoResult(t *testing.T) {
	h := setupTestHandler()
	req := httptest.NewRequest("GET", "/api/v1/result", nil)
	w := httptest.NewRecorder()

	h.HandleGetResult(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleAnalyze(t *testing.T) {
	h := setupTestHandler()

	text := "The cat sat on the mat. The cat chased the dog. The dog barked at the cat. The cat and the dog played together. The mat was soft for the cat."
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)
	part.Write([]byte(text))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/analyze", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.HandleAnalyze(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.True(t, result["nodes"].(float64) > 0)
}

func TestHandleUpdateConfig(t *testing.T) {
	h := setupTestHandler()

	body := `{"window_size": 10, "max_nodes": 50}`
	req := httptest.NewRequest("PATCH", "/api/v1/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.HandleUpdateConfig(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleExport_NoResult(t *testing.T) {
	h := setupTestHandler()
	req := httptest.NewRequest("GET", "/api/v1/export?format=json", nil)
	w := httptest.NewRecorder()

	h.HandleExport(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMiddleware_Recovery(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	handler := recoveryMiddleware(panicHandler)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestMiddleware_Logging(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := loggingMiddleware(inner)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
