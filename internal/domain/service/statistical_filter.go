package service

import (
	"math"
	"sort"

	"Philograph/internal/domain/model"
)

// StatisticalFilter は統計指標に基づいて共起ペアをフィルタリングする。
type StatisticalFilter struct {
	metric          model.Metric
	minCooccurrence int
}

// NewStatisticalFilter は新しいStatisticalFilterを返す。
func NewStatisticalFilter(metric model.Metric, minCooccurrence int) *StatisticalFilter {
	if minCooccurrence <= 0 {
		minCooccurrence = 2
	}
	return &StatisticalFilter{
		metric:          metric,
		minCooccurrence: minCooccurrence,
	}
}

// Filter は共起ペアの統計指標を計算し、最小共起回数でフィルタリングした上でソートして返す。
func (f *StatisticalFilter) Filter(
	pairs map[string]*model.CooccurrencePair,
	terms []*model.Term,
	totalWindows int,
) []*model.CooccurrencePair {
	if totalWindows <= 0 {
		totalWindows = 1
	}

	// Build frequency lookup by term ID
	freqByID := make(map[int]int, len(terms))
	for _, t := range terms {
		freqByID[t.ID] = t.Frequency
	}

	N := float64(totalWindows)

	var result []*model.CooccurrencePair
	for _, p := range pairs {
		if p.RawCount < f.minCooccurrence {
			continue
		}

		fA := float64(freqByID[p.TermAID])
		fB := float64(freqByID[p.TermBID])
		fAB := float64(p.RawCount)

		pAB := fAB / N
		pA := fA / N
		pB := fB / N

		if pA > 0 && pB > 0 && pAB > 0 {
			pmi := math.Log2(pAB / (pA * pB))
			p.PMI = pmi

			hAB := -math.Log2(pAB)
			if hAB > 0 {
				p.NPMI = pmi / hAB
			}
		}

		// Jaccard
		union := fA + fB - fAB
		if union > 0 {
			p.Jaccard = fAB / union
		}

		result = append(result, p)
	}

	// Sort by selected metric descending
	sort.Slice(result, func(i, j int) bool {
		return f.metricValue(result[i]) > f.metricValue(result[j])
	})

	return result
}

func (f *StatisticalFilter) metricValue(p *model.CooccurrencePair) float64 {
	switch f.metric {
	case model.MetricPMI:
		return p.PMI
	case model.MetricNPMI:
		return p.NPMI
	case model.MetricJaccard:
		return p.Jaccard
	case model.MetricFrequency:
		return float64(p.RawCount)
	default:
		return p.NPMI
	}
}
