package application

import (
	"context"
	"sync"
	"testing"

	"Philograph/internal/domain/model"
	"Philograph/internal/infrastructure/whitespace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSession_AnalyzeAndResult(t *testing.T) {
	tok := whitespace.NewTokenizer()
	pipeline := NewPipeline(tok, nil, nil)
	config := model.AnalysisConfig{
		WindowSize:      5,
		MinFrequency:    2,
		MinCooccurrence: 1,
		TargetPOS:       []model.POS{model.POSNoun},
		Metric:          model.MetricNPMI,
		MaxNodes:        150,
		Language:        model.LangEnglish,
	}
	session := NewSession(pipeline, config)

	text := "The cat sat on the mat. The cat chased the dog. The dog barked at the cat. The cat and the dog played together. The mat was soft for the cat."
	result, err := session.Analyze(context.Background(), text)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, result, session.Result())
	assert.True(t, session.HasText())
}

func TestSession_Reanalyze(t *testing.T) {
	tok := whitespace.NewTokenizer()
	pipeline := NewPipeline(tok, nil, nil)
	config := model.AnalysisConfig{
		WindowSize:      5,
		MinFrequency:    2,
		MinCooccurrence: 1,
		TargetPOS:       []model.POS{model.POSNoun},
		Metric:          model.MetricNPMI,
		MaxNodes:        150,
		Language:        model.LangEnglish,
	}
	session := NewSession(pipeline, config)

	// Reanalyze without text should fail
	_, err := session.Reanalyze(context.Background())
	assert.ErrorIs(t, err, ErrEmptyText)

	// Analyze first
	text := "The cat sat on the mat. The cat chased the dog. The dog barked at the cat. The cat and the dog played together. The mat was soft for the cat."
	_, err = session.Analyze(context.Background(), text)
	require.NoError(t, err)

	// Update config and reanalyze
	newConfig := config
	newConfig.MaxNodes = 10
	session.UpdateConfig(newConfig)

	result, err := session.Reanalyze(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSession_ConcurrentAccess(t *testing.T) {
	tok := whitespace.NewTokenizer()
	pipeline := NewPipeline(tok, nil, nil)
	config := model.AnalysisConfig{
		WindowSize:      5,
		MinFrequency:    2,
		MinCooccurrence: 1,
		TargetPOS:       []model.POS{model.POSNoun},
		Metric:          model.MetricNPMI,
		MaxNodes:        150,
		Language:        model.LangEnglish,
	}
	session := NewSession(pipeline, config)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = session.Config()
			_ = session.Result()
			_ = session.HasText()
		}()
	}
	wg.Wait()
}
