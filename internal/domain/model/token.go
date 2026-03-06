package model

// POS は品詞を表す文字列型。
type POS string

const (
	POSNoun      POS = "名詞"
	POSVerb      POS = "動詞"
	POSAdjective POS = "形容詞"
	POSAdverb    POS = "副詞"
	POSParticle  POS = "助詞"
	POSAuxVerb   POS = "助動詞"
	POSSymbol    POS = "記号"
	POSOther     POS = "その他"
)

// Token は形態素解析で得られた1トークンを表す。
type Token struct {
	Surface            string
	Base               string
	POS                POS
	SentenceIndex      int
	PositionInSentence int
}

// IsContent は内容語（名詞・動詞・形容詞・副詞）かどうかを返す。
func (t Token) IsContent() bool {
	switch t.POS {
	case POSNoun, POSVerb, POSAdjective, POSAdverb:
		return true
	default:
		return false
	}
}
