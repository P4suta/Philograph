package model

import "fmt"

// CooccurrencePair は2語間の共起関係を表す。
type CooccurrencePair struct {
	TermAID  int
	TermBID  int
	RawCount int
	PMI      float64
	NPMI     float64
	Jaccard  float64
}

// Key はマップのキーとして使用するための正規化済み文字列を返す。
func (p CooccurrencePair) Key() string {
	return fmt.Sprintf("%d:%d", p.TermAID, p.TermBID)
}

// NormalizePairOrder はIDの順序を正規化する。
func NormalizePairOrder(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}
