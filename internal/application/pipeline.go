package application

import (
	"context"
	"strings"

	"Philograph/internal/domain/model"
	"Philograph/internal/domain/service"
	"Philograph/internal/port"
)

// PipelineResult は分析パイプラインの結果。
type PipelineResult struct {
	Graph  *model.Graph
	Terms  []*model.Term
	Config model.AnalysisConfig
}

// Pipeline はNLP分析パイプラインのオーケストレーター。
type Pipeline struct {
	tokenizer port.Tokenizer
	analyzer  service.GraphAnalyzer
	listener  ProgressListener
}

// NewPipeline は新しいPipelineを返す。
func NewPipeline(tokenizer port.Tokenizer, analyzer service.GraphAnalyzer, listener ProgressListener) *Pipeline {
	return &Pipeline{
		tokenizer: tokenizer,
		analyzer:  analyzer,
		listener:  listener,
	}
}

// Run は分析パイプラインを実行する。
func (p *Pipeline) Run(ctx context.Context, text string, config model.AnalysisConfig) (*PipelineResult, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ErrEmptyText
	}

	// Step 1: Sentence split
	p.notify(Progress{Stage: StageSplitting, Percentage: 10, Message: "Splitting sentences..."})
	splitter := service.NewSentenceSplitter(config.Language)
	sentences := splitter.Split(text)
	if len(sentences) == 0 {
		return nil, ErrEmptyText
	}

	// Step 2: Tokenize
	p.notify(Progress{Stage: StageTokenizing, Percentage: 20, Message: "Tokenizing..."})
	var allTokens []model.Token
	for sentIdx, sent := range sentences {
		if err := ctx.Err(); err != nil {
			return nil, ErrAnalysisCancelled
		}
		tokens, err := p.tokenizer.Tokenize(ctx, sent)
		if err != nil {
			return nil, err
		}
		for i := range tokens {
			tokens[i].SentenceIndex = sentIdx
		}
		allTokens = append(allTokens, tokens...)
	}

	// Step 3: Filter tokens
	p.notify(Progress{Stage: StageFiltering, Percentage: 40, Message: "Filtering tokens..."})
	filter := service.NewTokenFilter(config.TargetPOS, config.StopWords, 2)
	filtered := filter.Filter(allTokens)

	// Step 4: Build vocabulary
	p.notify(Progress{Stage: StageVocabulary, Percentage: 50, Message: "Building vocabulary..."})
	termMap, terms := service.BuildVocabulary(filtered)

	// Apply min frequency filter
	terms = filterByMinFrequency(terms, config.MinFrequency)
	if len(terms) == 0 {
		return nil, ErrNoTerms
	}

	// Rebuild termMap with only surviving terms
	keptTerms := make(map[string]bool, len(terms))
	for _, t := range terms {
		keptTerms[t.Text] = true
	}
	for key := range termMap {
		if !keptTerms[key] {
			delete(termMap, key)
		}
	}

	// Re-filter tokens to only include kept terms
	var keptTokens []model.Token
	for _, t := range filtered {
		base := t.Base
		if base == "" {
			base = t.Surface
		}
		if keptTerms[base] {
			keptTokens = append(keptTokens, t)
		}
	}

	// Step 5: Extract co-occurrences
	p.notify(Progress{Stage: StageCooccur, Percentage: 60, Message: "Extracting co-occurrences..."})
	coBuilder := service.NewCooccurrenceBuilder(config.WindowSize)
	pairs := coBuilder.Build(keptTokens, termMap)

	// Step 6: Statistical filter
	p.notify(Progress{Stage: StageStatistics, Percentage: 70, Message: "Computing statistics..."})
	totalWindows := coBuilder.WindowCount(len(keptTokens))
	statFilter := service.NewStatisticalFilter(config.Metric, config.MinCooccurrence)
	sortedPairs := statFilter.Filter(pairs, terms, totalWindows)

	// Step 7: Build graph
	p.notify(Progress{Stage: StageGraphBuild, Percentage: 80, Message: "Building graph..."})
	graphBuilder := service.NewGraphBuilder(config.MaxNodes, p.analyzer)
	graph, err := graphBuilder.Build(sortedPairs, terms)
	if err != nil {
		return nil, err
	}

	p.notify(Progress{Stage: StageComplete, Percentage: 100, Message: "Complete"})

	return &PipelineResult{
		Graph:  graph,
		Terms:  terms,
		Config: config,
	}, nil
}

func (p *Pipeline) notify(prog Progress) {
	if p.listener != nil {
		p.listener(prog)
	}
}

func filterByMinFrequency(terms []*model.Term, minFreq int) []*model.Term {
	if minFreq <= 0 {
		return terms
	}
	result := make([]*model.Term, 0, len(terms))
	for _, t := range terms {
		if t.Frequency >= minFreq {
			result = append(result, t)
		}
	}
	return result
}
