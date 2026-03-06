package service

import (
	"testing"

	"Philograph/internal/domain/model"

	"github.com/stretchr/testify/assert"
)

func TestTokenFilter_Filter(t *testing.T) {
	filter := NewTokenFilter(
		[]model.POS{model.POSNoun, model.POSVerb},
		[]string{"stop", "word"},
		2,
	)

	tokens := []model.Token{
		{Surface: "hello", Base: "hello", POS: model.POSNoun},
		{Surface: "a", Base: "a", POS: model.POSNoun},         // too short
		{Surface: "stop", Base: "stop", POS: model.POSNoun},   // stopword
		{Surface: "run", Base: "running", POS: model.POSVerb},
		{Surface: "の", Base: "の", POS: model.POSParticle},    // wrong POS
	}

	result := filter.Filter(tokens)
	assert.Len(t, result, 2)
	assert.Equal(t, "hello", result[0].Base)
	assert.Equal(t, "running", result[1].Base)
}

func TestTokenFilter_CaseInsensitiveStopwords(t *testing.T) {
	filter := NewTokenFilter(
		[]model.POS{model.POSNoun},
		[]string{"The"},
		2,
	)

	tokens := []model.Token{
		{Surface: "the", Base: "the", POS: model.POSNoun},
		{Surface: "THE", Base: "THE", POS: model.POSNoun},
	}

	result := filter.Filter(tokens)
	assert.Len(t, result, 0)
}
