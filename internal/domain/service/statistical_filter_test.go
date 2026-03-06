package service

import (
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
)

func TestStatisticalFilter_Filter(t *testing.T) {
	terms := []*model.Term{
		{ID: 0, Text: "cat", Frequency: 10},
		{ID: 1, Text: "dog", Frequency: 8},
		{ID: 2, Text: "bird", Frequency: 3},
	}

	pairs := map[string]*model.CooccurrencePair{
		"0:1": {TermAID: 0, TermBID: 1, RawCount: 5},
		"0:2": {TermAID: 0, TermBID: 2, RawCount: 1}, // below min cooccurrence
		"1:2": {TermAID: 1, TermBID: 2, RawCount: 3},
	}

	filter := NewStatisticalFilter(model.MetricNPMI, 2)
	result := filter.Filter(pairs, terms, 100)

	// Only pairs with RawCount >= 2 should pass
	assert.Len(t, result, 2)

	// Check that PMI/NPMI/Jaccard were computed
	for _, p := range result {
		assert.NotZero(t, p.PMI)
		assert.NotZero(t, p.NPMI)
		assert.NotZero(t, p.Jaccard)
	}
}

func TestStatisticalFilter_SortByMetric(t *testing.T) {
	terms := []*model.Term{
		{ID: 0, Text: "a", Frequency: 10},
		{ID: 1, Text: "b", Frequency: 5},
		{ID: 2, Text: "c", Frequency: 20},
	}

	pairs := map[string]*model.CooccurrencePair{
		"0:1": {TermAID: 0, TermBID: 1, RawCount: 4},
		"0:2": {TermAID: 0, TermBID: 2, RawCount: 3},
	}

	filter := NewStatisticalFilter(model.MetricFrequency, 2)
	result := filter.Filter(pairs, terms, 100)

	assert.Len(t, result, 2)
	// Frequency sorted descending: 4 > 3
	assert.Equal(t, 4, result[0].RawCount)
	assert.Equal(t, 3, result[1].RawCount)
}
