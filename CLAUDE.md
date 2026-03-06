# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Philograph is a Go CLI tool that builds co-occurrence networks from book-length texts and visualizes them interactively in the browser. It takes a UTF-8 text file, performs morphological analysis (Japanese via kagome, English via whitespace tokenization), extracts co-occurrence pairs, computes graph metrics, and serves an interactive Sigma.js/WebGL visualization. Everything ships as a single binary with no external runtime dependencies.

## Build & Development Commands

```bash
# Build
go build ./cmd/philograph/

# Run
go run ./cmd/philograph/ <textfile.txt>

# Test
go test ./...

# Run a single test
go test ./internal/domain/service/ -run TestCooccurrenceBuilder

# Lint (if golangci-lint is installed)
golangci-lint run
```

## Architecture

The project follows **Clean Architecture** with **Ports & Adapters (Hexagonal)** pattern:

- **`cmd/philograph/`** — Entry point. Manual dependency injection (no DI framework). Detects language, wires tokenizer and all services, starts HTTP server.
- **`internal/domain/model/`** — Value objects with zero external dependencies: `Token`, `Term`, `CooccurrencePair`, `Graph`/`Node`/`Edge`, `AnalysisConfig`, `Language`, `Metric` enums.
- **`internal/domain/service/`** — Stateless domain services: sentence splitting, token filtering, co-occurrence extraction, statistical filtering (PMI/NPMI/Jaccard), graph building with centrality & Louvain community detection.
- **`internal/port/`** — Interfaces: `Tokenizer` (morphological analysis), `Exporter` (graph data output). Domain depends on these; infrastructure implements them.
- **`internal/infrastructure/`** — Adapter implementations: `kagome/` (Japanese tokenizer), `whitespace/` (English tokenizer), `graph/` (gonum adapter for centrality/community), `export/` (JSON, GEXF exporters).
- **`internal/application/`** — Pipeline orchestration (`Pipeline`), session/cache management (`Session`), progress notification (`ProgressListener`).
- **`internal/api/`** — HTTP server, REST handlers, WebSocket hub for real-time progress, middleware. Uses `net/http` enhanced routing (Go 1.22+).
- **`web/`** — Frontend assets embedded via `//go:embed`. Vanilla JS (ES Modules) with Sigma.js v2 + Graphology. No build step.

## NLP Pipeline Flow

Text → Sentence Split → Tokenize (via port) → Filter Tokens (content words only) → Build Vocabulary → Extract Co-occurrences (windowed, within sentence boundaries) → Statistical Filter → Build Graph → Compute Centrality & Communities → `Graph` model

## Key Design Decisions

- **Language detection** uses Unicode codepoint distribution (CJK ratio > 30% → Japanese).
- **Co-occurrence pairs** are order-normalized (`TermAID < TermBID`) since the relationship is undirected.
- **POS tags** are mapped from kagome's tag system to domain-specific constants (e.g., `POSNoun = "名詞"`).
- **Ephemeral design** — no persistence layer. Analysis results live in memory only.
- **`AnalysisConfig.DefaultConfig()`** provides sensible defaults (window=5, minFreq=3, metric=NPMI, maxNodes=150).

## Tech Stack

- **Go 1.23+** with `embed.FS`, `log/slog`, enhanced `net/http` routing
- **kagome v2** — Pure Go Japanese morphological analyzer (IPAdic dictionary bundled in binary, ~50MB)
- **gonum** — Graph computation (centrality, Louvain community detection), sparse matrix ops
- **nhooyr.io/websocket** — WebSocket for progress notifications
- **cobra** — CLI framework
- **testify** — Test assertions
- **Sigma.js v2 + Graphology** — WebGL graph rendering (frontend)
- **GoReleaser** — Cross-platform binary builds

## Development Environment

### 1Password + GitHub 認証

このプロジェクトはWSL2上で開発しており、認証情報は1Passwordで管理している。

- **Git commit署名**: 1PasswordのSSH Agent統合を使用。Windows側の`op-ssh-sign-wsl.exe`がWSLから呼び出される。`commit.gpgsign=true`で全コミットに自動署名。
- **GitHub CLI (`gh`)**: `gh auth login`で認証済み。HTTPS経由でgit操作を行う。credential helperとして`gh auth git-credential`を使用。
- **1Password CLI (`op`)**: WSL内には未インストール。必要であれば`op plugin init gh`でGitHub CLIとの連携が可能。現状は`gh`単体で認証している。
- **SSH鍵**: 1Passwordが管理するed25519鍵をsigning keyとして使用（`gpg.format=ssh`）。

### git設定（グローバル）

```
user.name=P4suta
user.email=42543015+P4suta@users.noreply.github.com
gpg.format=ssh
gpg.ssh.program=/mnt/c/Users/livec/AppData/Local/Microsoft/WindowsApps/op-ssh-sign-wsl.exe
commit.gpgsign=true
core.sshcommand=ssh.exe
credential helper=gh auth git-credential
```

## Git運用規約

### ブランチ戦略

`main`ブランチを唯一の長寿命ブランチとし、トピックブランチを切って作業する。

| ブランチ | 命名規則 | 用途 |
|---------|---------|------|
| `main` | — | 安定版。直接コミット禁止。マージのみ。 |
| feature/ | `feature/<短い説明>` | 新機能の開発 |
| fix/ | `fix/<短い説明>` | バグ修正 |
| refactor/ | `refactor/<短い説明>` | リファクタリング |
| docs/ | `docs/<短い説明>` | ドキュメントのみの変更 |
| chore/ | `chore/<短い説明>` | ビルド設定、依存関係、CI等 |

- ブランチ名は英語小文字、ハイフン区切り（例: `feature/sentence-splitter`）
- 1つのブランチ = 1つの論理的な変更単位（レイヤー1つ分、機能1つ分など）
- 作業完了後は`main`へマージし、マージ済みブランチは削除する

### コミット規約

[Conventional Commits](https://www.conventionalcommits.org/) に準拠する。

```
<type>(<scope>): <summary>

<body>（任意）

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
```

**type** (必須):
- `feat` — 新機能
- `fix` — バグ修正
- `refactor` — 機能変更なしのコード改善
- `test` — テストの追加・修正
- `docs` — ドキュメントのみの変更
- `chore` — ビルド、CI、依存関係の変更
- `style` — フォーマット変更（動作に影響なし）

**scope** (推奨): 変更対象のパッケージやレイヤー
- 例: `model`, `service`, `api`, `kagome`, `web`, `pipeline`, `cli`

**summary**: 命令形の英語で簡潔に（50文字以内目安）

**コミット粒度のガイドライン**:
- 1コミット = 1つの論理的変更。「何を・なぜ」が明確に説明できる単位。
- コンパイルが通り、テストがパスする状態でコミットする。
- 大きな機能は複数コミットに分割する:
  1. モデル/型定義 → 2. ドメインロジック → 3. インフラ実装 → 4. API/結合 → 5. テスト追加
- やってはいけないこと: 複数の無関係な変更を1コミットに混ぜる、「WIP」だけのコミットメッセージ

### マージ方針

- トピックブランチから`main`へは **`git merge --no-ff`** を使う（マージコミットを残す）
- マージ前に`main`の最新を取り込む: `git merge main`（トピックブランチ上で）
- コンフリクトはトピックブランチ側で解消してからマージする

### Claudeへの指示（Git操作）

- コミットを求められたら、上記の規約に従ったメッセージを作成すること
- ブランチ作成を求められたら、命名規則に従うこと
- `main`に直接コミットしないこと（初回セットアップコミットを除く）
- コミット前に`go build ./cmd/philograph/`と`go test ./...`が通ることを確認すること（ビルド可能な状態であれば）
- マージ済みブランチの削除を提案すること
