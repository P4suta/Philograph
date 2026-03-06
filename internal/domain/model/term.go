package model

// Term は語彙（レンマ化済み）を表す。
type Term struct {
	ID        int
	Text      string
	POS       POS
	Frequency int
}
