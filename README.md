# Philograph

**テキストから共起ネットワークを構築し、ブラウザ上でインタラクティブに可視化するCLIツール**

Philograph は書籍一冊分のテキストファイルを入力に、形態素解析・共起抽出・統計フィルタリング・グラフ分析を行い、Sigma.js/WebGL によるインタラクティブなネットワーク可視化をブラウザ上で提供します。すべてのアセットを単一バイナリに同梱し、外部ランタイム依存なしで動作します。

> プロジェクト名は古代ギリシャ語の φίλος（愛する）と γράφω（描く）に由来し、「概念の関係を描き出す」という本ツールの目的を表現しています。

## Features

- **ゼロ設定** — テキストファイル1つで即座に分析開始
- **ゼロ引数起動** — 引数なしでサーバー起動、ブラウザからファイルアップロードで分析開始
- **日本語/英語対応** — kagome (日本語形態素解析) / whitespace tokenizer (英語) を自動切替
- **動的言語切替** — AutoTokenizer によりアップロードごとにテキスト言語を自動検出・切替
- **統計フィルタリング** — PMI / NPMI / Jaccard によるノイズ除去
- **グラフ分析** — 中心性指標・Louvain コミュニティ検出
- **インタラクティブ可視化** — Sigma.js v2 + WebGL による数百ノード規模のリアルタイム描画
- **単一バイナリ配布** — Go embed でフロントエンド資材を内包、インストール不要
- **リアルタイム進捗** — WebSocket による分析進捗の通知

## Quick Start

```bash
# ビルド
go build ./cmd/philograph/

# サーバーのみ起動（ブラウザからファイルアップロードで分析開始）
./philograph

# ファイル指定で即座に分析（従来通り）
./philograph your-text.txt
```

## Architecture

Clean Architecture (Ports & Adapters) パターンを採用しています。

```
cmd/philograph/          # エントリポイント・DI
internal/
  domain/
    model/               # 値オブジェクト（Token, Term, Graph, etc.）
    service/             # ドメインサービス（共起抽出, フィルタリング, グラフ構築）
  port/                  # インターフェース（Tokenizer, Exporter）
  infrastructure/        # アダプタ実装（kagome, whitespace, autotokenizer, gonum, export）
  application/           # パイプライン, セッション管理
  api/                   # HTTPサーバー, REST, WebSocket
web/                     # フロントエンド（Sigma.js, go:embed）
```

## NLP Pipeline

```
テキスト → 文分割 → 形態素解析 → 内容語フィルタ → 語彙構築
→ 共起抽出（窓幅ベース, 文境界考慮） → 統計フィルタ
→ グラフ構築 → 中心性・コミュニティ検出 → 可視化
```

## Tech Stack

| Category | Technology |
|----------|-----------|
| Language | Go 1.23+ |
| NLP (Japanese) | kagome v2 (Pure Go, IPAdic) |
| Graph computation | gonum |
| WebSocket | nhooyr.io/websocket |
| CLI | cobra |
| Visualization | Sigma.js v2 + Graphology (WebGL) |
| Testing | testify |

## Development

```bash
go test ./...              # 全テスト実行
go build ./cmd/philograph/ # ビルド
```

詳細な設計については [`Philograph 技術設計書.md`](Philograph%20技術設計書.md) を参照してください。

## License

[MIT](LICENSE)
