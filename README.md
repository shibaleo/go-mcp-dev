# go-mcp-dev

Go言語で実装したMCP (Model Context Protocol) サーバー。Claude CodeなどのLLMクライアントから外部APIを操作するためのシングルテナント用ツール。

## 特徴

- **Go + SSE**: 軽量・高速なJSON-RPC 2.0 over Server-Sent Events
- **シングルテナント**: 固定シークレット認証、個人用途に最適化
- **$0運用**: Koyeb Free Tier + GitHub Actions pingでコールドスタート回避
- **オブザーバビリティ**: Grafana Cloud Lokiへのリアルタイムログ送信

## 対応モジュール

| モジュール | 状態 | 説明 |
|-----------|------|------|
| Supabase | 実装済 | Management API（プロジェクト一覧、SQL実行） |
| Notion | 予定 | ページ・データベース操作 |
| GitHub | 予定 | リポジトリ、Issue、PR操作 |
| Jira | 予定 | Issue/Project操作 |
| Confluence | 予定 | Space/Page操作 |

## エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/health` | ヘルスチェック |
| POST | `/mcp` | JSON-RPC 2.0 over SSE |

## 利用可能なツール

### Supabase

- `supabase_list_projects` - プロジェクト一覧取得
- `supabase_run_query` - SQL実行

## セットアップ

### 前提条件

- Go 1.22+
- Docker
- Koyeb CLI

### ローカル開発

```bash
# 環境変数設定
cp .env.example .env.development
# .env.development を編集

# 起動
docker-compose up

# テスト
go test ./...
```

### 本番デプロイ

```bash
# 初回デプロイ
bash deploy.sh

# 環境変数更新
koyeb service update go-mcp-dev/go-mcp-dev --env KEY=VALUE
```

## 環境変数

| 変数 | 説明 |
|------|------|
| `INTERNAL_SECRET` | MCP認証用Bearer token |
| `SUPABASE_ACCESS_TOKEN` | Supabase Management API token |
| `GRAFANA_LOKI_URL` | Loki Push API エンドポイント |
| `GRAFANA_LOKI_USER` | Grafana Cloud ユーザーID |
| `GRAFANA_LOKI_API_KEY` | Grafana Cloud API Key |

## ユースケース

```
┌─────────────────────────────────────────────────────────────────┐
│  開発者（shibaleo）                                              │
│                                                                 │
│  ┌─────────────────┐                                            │
│  │   Claude Code   │  「Supabaseのテーブル一覧を見せて」         │
│  │   (VSCode拡張)  │  「このSQLを実行して」                      │
│  └────────┬────────┘                                            │
│           │                                                     │
│           │ MCP Protocol (JSON-RPC 2.0 over SSE)                │
│           │ Authorization: Bearer <INTERNAL_SECRET>             │
│           ↓                                                     │
│  ┌─────────────────┐                                            │
│  │  go-mcp-dev     │  shibaleo-dev.mcpist.app                   │
│  │  (Koyeb)        │  - 認証検証                                │
│  └────────┬────────┘  - ツール実行                              │
│           │           - Lokiへログ送信                          │
│           │                                                     │
│           │ 各サービスのAPIトークン                              │
│           ↓                                                     │
│  ┌─────────────────────────────────────────────┐                │
│  │  External APIs                              │                │
│  │  - Supabase Management API                  │                │
│  │  - Notion API (予定)                        │                │
│  │  - GitHub API (予定)                        │                │
│  │  - Jira/Confluence API (予定)               │                │
│  └─────────────────────────────────────────────┘                │
└─────────────────────────────────────────────────────────────────┘
```

### 想定ユーザー

- **対象**: 開発者本人（シングルテナント）
- **用途**: Claude Codeから自然言語で外部サービスを操作
- **例**:
  - 「go-mcp-demoプロジェクトのテーブル一覧を見せて」
  - 「usersテーブルにカラムを追加するSQLを実行して」
  - 「Notionの今週のタスクを取得して」（予定）

## CI/CD

- **PR時**: ビルド・テスト実行
- **main push時**: Koyeb自動再デプロイ + ヘルスチェック
- **45分ごと**: GitHub Actionsからping（スリープ回避）

## ライセンス

Private
