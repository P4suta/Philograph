package application

import (
	"context"
	"testing"

	"Philograph/internal/domain/model"
	"Philograph/internal/infrastructure/whitespace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPipeline_Run(t *testing.T) {
	tok := whitespace.NewTokenizer()
	pipeline := NewPipeline(tok, nil, nil)

	text := "The cat sat on the mat. The cat chased the dog. The dog barked at the cat. The cat and the dog played together. The mat was soft for the cat."
	config := model.AnalysisConfig{
		WindowSize:      5,
		MinFrequency:    2,
		MinCooccurrence: 1,
		TargetPOS:       []model.POS{model.POSNoun},
		Metric:          model.MetricNPMI,
		MaxNodes:        150,
		Language:        model.LangEnglish,
	}

	result, err := pipeline.Run(context.Background(), text, config)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Graph.NodeCount() > 0, "should have nodes")
	assert.True(t, result.Graph.EdgeCount() > 0, "should have edges")
}

func TestPipeline_EmptyText(t *testing.T) {
	tok := whitespace.NewTokenizer()
	pipeline := NewPipeline(tok, nil, nil)

	_, err := pipeline.Run(context.Background(), "", model.DefaultConfig())
	assert.ErrorIs(t, err, ErrEmptyText)
}

func TestPipeline_Cancelled(t *testing.T) {
	tok := whitespace.NewTokenizer()
	pipeline := NewPipeline(tok, nil, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	text := "The cat sat on the mat. The cat chased the dog."
	config := model.DefaultConfig()
	config.Language = model.LangEnglish

	_, err := pipeline.Run(ctx, text, config)
	assert.ErrorIs(t, err, ErrAnalysisCancelled)
}

func TestPipeline_ProgressNotification(t *testing.T) {
	tok := whitespace.NewTokenizer()
	var stages []Stage
	listener := func(p Progress) {
		stages = append(stages, p.Stage)
	}
	pipeline := NewPipeline(tok, nil, listener)

	text := "The cat sat on the mat. The cat chased the dog. The dog barked at the cat. The cat and the dog played together. The mat was soft for the cat."
	config := model.AnalysisConfig{
		WindowSize:      5,
		MinFrequency:    2,
		MinCooccurrence: 1,
		TargetPOS:       []model.POS{model.POSNoun},
		Metric:          model.MetricNPMI,
		MaxNodes:        150,
		Language:        model.LangEnglish,
	}

	_, err := pipeline.Run(context.Background(), text, config)
	require.NoError(t, err)

	assert.Contains(t, stages, StageSplitting)
	assert.Contains(t, stages, StageComplete)
}
