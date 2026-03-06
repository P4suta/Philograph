package application

import (
	"context"
	"sync"

	"Philograph/internal/domain/model"
)

// Session は分析セッションを管理する。
type Session struct {
	mu       sync.RWMutex
	text     string
	config   model.AnalysisConfig
	result   *PipelineResult
	pipeline *Pipeline
}

// NewSession は新しいSessionを返す。
func NewSession(pipeline *Pipeline, config model.AnalysisConfig) *Session {
	return &Session{
		pipeline: pipeline,
		config:   config,
	}
}

// Analyze はテキストを分析する。
func (s *Session) Analyze(ctx context.Context, text string) (*PipelineResult, error) {
	s.mu.Lock()
	s.text = text
	config := s.config
	s.mu.Unlock()

	result, err := s.pipeline.Run(ctx, text, config)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.result = result
	s.mu.Unlock()

	return result, nil
}

// Reanalyze は保存済みテキストを現在の設定で再分析する。
func (s *Session) Reanalyze(ctx context.Context) (*PipelineResult, error) {
	s.mu.RLock()
	text := s.text
	s.mu.RUnlock()

	if text == "" {
		return nil, ErrEmptyText
	}

	return s.Analyze(ctx, text)
}

// UpdateConfig は設定を更新する。
func (s *Session) UpdateConfig(config model.AnalysisConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = config
}

// Config は現在の設定を返す。
func (s *Session) Config() model.AnalysisConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// Result は最新の分析結果を返す。
func (s *Session) Result() *PipelineResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.result
}

// HasText はテキストが設定済みかどうかを返す。
func (s *Session) HasText() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.text != ""
}
