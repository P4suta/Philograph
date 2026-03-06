package kagome

import (
	"context"

	"Philograph/internal/domain/model"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

// Tokenizer はkagomeを使った日本語形態素解析トークナイザー。
type Tokenizer struct {
	tok *tokenizer.Tokenizer
}

// NewTokenizer は新しいTokenizerを返す。
func NewTokenizer() (*Tokenizer, error) {
	tok, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}
	return &Tokenizer{tok: tok}, nil
}

// Language はこのトークナイザーの対象言語を返す。
func (t *Tokenizer) Language() model.Language {
	return model.LangJapanese
}

// Tokenize はテキストを形態素解析してトークンのスライスを返す。
func (t *Tokenizer) Tokenize(_ context.Context, sentence string) ([]model.Token, error) {
	morphs := t.tok.Tokenize(sentence)

	tokens := make([]model.Token, 0, len(morphs))
	for i, m := range morphs {
		features := m.Features()
		if len(features) == 0 {
			continue
		}

		pos := mapPOS(features[0])
		base := m.Surface
		if len(features) >= 7 && features[6] != "*" {
			base = features[6]
		}

		tokens = append(tokens, model.Token{
			Surface:            m.Surface,
			Base:               base,
			POS:                pos,
			PositionInSentence: i,
		})
	}

	return tokens, nil
}

func mapPOS(pos string) model.POS {
	switch pos {
	case "名詞":
		return model.POSNoun
	case "動詞":
		return model.POSVerb
	case "形容詞":
		return model.POSAdjective
	case "副詞":
		return model.POSAdverb
	case "助詞":
		return model.POSParticle
	case "助動詞":
		return model.POSAuxVerb
	case "記号":
		return model.POSSymbol
	default:
		return model.POSOther
	}
}
