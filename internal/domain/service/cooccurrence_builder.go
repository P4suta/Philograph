package service

import (
	"Philograph/internal/domain/model"
)

// CooccurrenceBuilder はトークン列からスライディングウィンドウ方式で共起ペアを抽出する。
type CooccurrenceBuilder struct {
	windowSize int
}

// NewCooccurrenceBuilder は新しいCooccurrenceBuilderを返す。
func NewCooccurrenceBuilder(windowSize int) *CooccurrenceBuilder {
	if windowSize <= 0 {
		windowSize = 5
	}
	return &CooccurrenceBuilder{windowSize: windowSize}
}

// BuildVocabulary はフィルタ済みトークンから語彙マップを構築する。
func BuildVocabulary(tokens []model.Token) (map[string]*model.Term, []*model.Term) {
	termMap := make(map[string]*model.Term)
	var terms []*model.Term
	nextID := 0

	for i := range tokens {
		base := tokens[i].Base
		if base == "" {
			base = tokens[i].Surface
		}

		if _, exists := termMap[base]; !exists {
			term := &model.Term{
				ID:        nextID,
				Text:      base,
				POS:       tokens[i].POS,
				Frequency: 0,
			}
			termMap[base] = term
			terms = append(terms, term)
			nextID++
		}
		termMap[base].Frequency++
	}

	return termMap, terms
}

// Build はフィルタ済みトークンと語彙マップから共起ペアを抽出する。
// 文境界を考慮し、同一文内のトークン間でのみ共起を計算する。
func (b *CooccurrenceBuilder) Build(tokens []model.Token, termMap map[string]*model.Term) map[string]*model.CooccurrencePair {
	pairs := make(map[string]*model.CooccurrencePair)

	// Group tokens by sentence
	sentences := groupBySentence(tokens)

	for _, sentTokens := range sentences {
		for i := 0; i < len(sentTokens); i++ {
			baseI := sentTokens[i].Base
			if baseI == "" {
				baseI = sentTokens[i].Surface
			}
			termI, okI := termMap[baseI]
			if !okI {
				continue
			}

			end := i + b.windowSize
			if end > len(sentTokens) {
				end = len(sentTokens)
			}

			for j := i + 1; j < end; j++ {
				baseJ := sentTokens[j].Base
				if baseJ == "" {
					baseJ = sentTokens[j].Surface
				}
				termJ, okJ := termMap[baseJ]
				if !okJ {
					continue
				}

				if termI.ID == termJ.ID {
					continue
				}

				aID, bID := model.NormalizePairOrder(termI.ID, termJ.ID)
				pair := model.CooccurrencePair{
					TermAID: aID,
					TermBID: bID,
				}
				key := pair.Key()

				if existing, ok := pairs[key]; ok {
					existing.RawCount++
				} else {
					pair.RawCount = 1
					pairs[key] = &pair
				}
			}
		}
	}

	return pairs
}

// groupBySentence はトークンをSentenceIndexごとにグループ化する。
func groupBySentence(tokens []model.Token) [][]model.Token {
	if len(tokens) == 0 {
		return nil
	}

	sentenceMap := make(map[int][]model.Token)
	maxIdx := 0
	for _, t := range tokens {
		sentenceMap[t.SentenceIndex] = append(sentenceMap[t.SentenceIndex], t)
		if t.SentenceIndex > maxIdx {
			maxIdx = t.SentenceIndex
		}
	}

	result := make([][]model.Token, 0, len(sentenceMap))
	for i := 0; i <= maxIdx; i++ {
		if s, ok := sentenceMap[i]; ok {
			result = append(result, s)
		}
	}
	return result
}

// WindowCount はトークン数から抽出可能なウィンドウ数を計算する。
func (b *CooccurrenceBuilder) WindowCount(tokenCount int) int {
	if tokenCount <= 1 {
		return 0
	}
	if tokenCount <= b.windowSize {
		return 1
	}
	return tokenCount - b.windowSize + 1
}
