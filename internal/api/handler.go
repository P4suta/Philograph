package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"Philograph/internal/application"
	"Philograph/internal/domain/model"
	"Philograph/internal/port"
)

const maxUploadSize = 50 << 20 // 50 MB

// Handler はREST APIハンドラー。
type Handler struct {
	session   *application.Session
	exporters map[string]port.Exporter
}

// NewHandler は新しいHandlerを返す。
func NewHandler(session *application.Session, exporters map[string]port.Exporter) *Handler {
	return &Handler{
		session:   session,
		exporters: exporters,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

// HandleAnalyze はファイルアップロード→分析。
func (h *Handler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file required: "+err.Error())
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusRequestEntityTooLarge, "file too large")
		return
	}

	text := string(data)
	result, err := h.session.Analyze(r.Context(), text)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case application.ErrEmptyText:
			status = http.StatusBadRequest
		case application.ErrNoTerms:
			status = http.StatusUnprocessableEntity
		case application.ErrFileTooLarge:
			status = http.StatusRequestEntityTooLarge
		case application.ErrAnalysisCancelled:
			status = http.StatusServiceUnavailable
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"nodes": result.Graph.NodeCount(),
		"edges": result.Graph.EdgeCount(),
	})
}

// HandleGetConfig は現在の設定を返す。
func (h *Handler) HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	config := h.session.Config()
	writeJSON(w, http.StatusOK, config)
}

// HandleUpdateConfig はパラメータ変更→再分析。
func (h *Handler) HandleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var update struct {
		WindowSize      *int    `json:"window_size"`
		MinFrequency    *int    `json:"min_frequency"`
		MinCooccurrence *int    `json:"min_cooccurrence"`
		MaxNodes        *int    `json:"max_nodes"`
		Metric          *string `json:"metric"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	config := h.session.Config()

	if update.WindowSize != nil {
		config.WindowSize = *update.WindowSize
	}
	if update.MinFrequency != nil {
		config.MinFrequency = *update.MinFrequency
	}
	if update.MinCooccurrence != nil {
		config.MinCooccurrence = *update.MinCooccurrence
	}
	if update.MaxNodes != nil {
		config.MaxNodes = *update.MaxNodes
	}
	if update.Metric != nil {
		switch strings.ToLower(*update.Metric) {
		case "pmi":
			config.Metric = model.MetricPMI
		case "npmi":
			config.Metric = model.MetricNPMI
		case "jaccard":
			config.Metric = model.MetricJaccard
		case "frequency":
			config.Metric = model.MetricFrequency
		}
	}

	h.session.UpdateConfig(config)

	if !h.session.HasText() {
		writeJSON(w, http.StatusOK, map[string]string{"status": "config updated"})
		return
	}

	result, err := h.session.Reanalyze(r.Context())
	if err != nil {
		slog.Error("reanalysis failed", "error", err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"nodes": result.Graph.NodeCount(),
		"edges": result.Graph.EdgeCount(),
	})
}

// HandleGetResult は最新結果を返す。
func (h *Handler) HandleGetResult(w http.ResponseWriter, r *http.Request) {
	result := h.session.Result()
	if result == nil {
		writeError(w, http.StatusNotFound, "no analysis result available")
		return
	}
	writeJSON(w, http.StatusOK, result.Graph)
}

// HandleExport はグラフをエクスポートする。
func (h *Handler) HandleExport(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	exporter, ok := h.exporters[format]
	if !ok {
		writeError(w, http.StatusBadRequest, "unsupported format: "+format)
		return
	}

	result := h.session.Result()
	if result == nil {
		writeError(w, http.StatusNotFound, "no analysis result available")
		return
	}

	w.Header().Set("Content-Type", exporter.ContentType())
	w.Header().Set("Content-Disposition", "attachment; filename=graph"+exporter.FileExtension())
	if err := exporter.Export(w, result.Graph); err != nil {
		slog.Error("export failed", "error", err)
	}
}

// HandleHealth はヘルスチェック。
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
