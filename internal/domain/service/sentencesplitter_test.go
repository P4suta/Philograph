package service

import (
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
)

func TestSentenceSplitter_SplitJapanese(t *testing.T) {
	splitter := NewSentenceSplitter(model.LangJapanese)

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic sentences",
			input:    "吾輩は猫である。名前はまだ無い。",
			expected: []string{"吾輩は猫である。", "名前はまだ無い。"},
		},
		{
			name:     "exclamation and question marks",
			input:    "本当ですか？はい！",
			expected: []string{"本当ですか？", "はい！"},
		},
		{
			name:     "newline as delimiter",
			input:    "第一段落\n第二段落",
			expected: []string{"第一段落", "第二段落"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "trailing text without delimiter",
			input:    "これはテスト",
			expected: []string{"これはテスト"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitter.Split(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSentenceSplitter_SplitEnglish(t *testing.T) {
	splitter := NewSentenceSplitter(model.LangEnglish)

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic sentences",
			input:    "Hello world. How are you?",
			expected: []string{"Hello world.", "How are you?"},
		},
		{
			name:     "abbreviation handling",
			input:    "Dr. Smith went home. He was tired.",
			expected: []string{"Dr. Smith went home.", "He was tired."},
		},
		{
			name:     "exclamation",
			input:    "Wow! That is amazing.",
			expected: []string{"Wow!", "That is amazing."},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitter.Split(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
