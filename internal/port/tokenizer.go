package port

import (
	"context"

	"Philograph/internal/domain/model"
)

// Tokenizer は形態素解析を行うポートインターフェース。
type Tokenizer interface {
	Tokenize(ctx context.Context, sentence string) ([]model.Token, error)
	Language() model.Language
}
