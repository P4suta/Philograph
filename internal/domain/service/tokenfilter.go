package service

import (
	"strings"
	"unicode/utf8"

	"Philograph/internal/domain/model"
)

// TokenFilter はトークンをフィルタリングする。
type TokenFilter struct {
	targetPOS     map[model.POS]bool
	stopWords     map[string]bool
	minCharLength int
}

// NewTokenFilter は新しいTokenFilterを返す。
func NewTokenFilter(targetPOS []model.POS, stopWords []string, minCharLength int) *TokenFilter {
	posMap := make(map[model.POS]bool, len(targetPOS))
	for _, p := range targetPOS {
		posMap[p] = true
	}

	swMap := make(map[string]bool, len(stopWords))
	for _, w := range stopWords {
		swMap[strings.ToLower(w)] = true
	}

	if minCharLength <= 0 {
		minCharLength = 2
	}

	return &TokenFilter{
		targetPOS:     posMap,
		stopWords:     swMap,
		minCharLength: minCharLength,
	}
}

// Filter はトークンのスライスからフィルタリング条件に合うものだけを返す。
func (f *TokenFilter) Filter(tokens []model.Token) []model.Token {
	result := make([]model.Token, 0, len(tokens)/2)
	for _, t := range tokens {
		if f.accept(t) {
			result = append(result, t)
		}
	}
	return result
}

func (f *TokenFilter) accept(t model.Token) bool {
	if !f.targetPOS[t.POS] {
		return false
	}

	base := t.Base
	if base == "" {
		base = t.Surface
	}

	if utf8.RuneCountInString(base) < f.minCharLength {
		return false
	}

	if f.stopWords[strings.ToLower(base)] {
		return false
	}

	return true
}

// DefaultStopWordsJapanese は日本語のデフォルトストップワード。
var DefaultStopWordsJapanese = []string{
	"これ", "それ", "あれ", "この", "その", "あの",
	"ここ", "そこ", "あそこ",
	"こちら", "そちら", "あちら",
	"よう", "こと", "もの", "ため", "とき", "ところ",
	"なに", "なん",
	"わたし", "あなた", "かれ", "かのじょ",
	"われ", "われわれ",
	"する", "いる", "ある", "なる", "できる",
	"れる", "られる", "せる", "させる",
}

// DefaultStopWordsEnglish は英語のデフォルトストップワード。
var DefaultStopWordsEnglish = []string{
	"the", "be", "to", "of", "and", "a", "in", "that", "have", "i",
	"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
	"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
	"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
	"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
	"when", "make", "can", "like", "no", "just", "him", "know", "take",
	"people", "into", "year", "your", "good", "some", "could", "them", "see",
	"other", "than", "then", "now", "look", "only", "come", "its", "over",
	"think", "also", "back", "after", "use", "two", "how", "our", "work",
	"first", "well", "way", "even", "new", "want", "because", "any", "these",
	"give", "day", "most", "us", "was", "were", "been", "has", "had", "are",
	"is", "am", "did", "does", "being", "having",
}
