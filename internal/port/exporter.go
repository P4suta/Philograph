package port

import (
	"io"

	"Philograph/internal/domain/model"
)

// Exporter はグラフデータをエクスポートするポートインターフェース。
type Exporter interface {
	Export(w io.Writer, graph *model.Graph) error
	ContentType() string
	FileExtension() string
}
