# go-mcp-dev

## MCPist: Managing Context Personally

**MCPist is a personal context management tool for power users who want to precisely control what an LLM can see and do, beyond the average context provided by MCP.**

MCP is the protocol.
MCPist is how I live with it.

Go言語で実装したMCP (Model Context Protocol) サーバー。Claude CodeなどのLLMクライアントから外部APIを操作するためのシングルテナント用ツール。

---

## Why MCPist? The Vendor Lock-in Problem

Model Context Protocol (MCP) は「AI界のUSB-C」として期待されていますが、**2025-11-25の仕様改定により構造的な問題が明確になりました**。

### The Problem: Who Benefits from Portability?

MCP 2025-11-25で導入された**Enterprise-Managed Authorization (XAA)** は、エンタープライズIdP（Okta等）との連携を標準化しました。これにより企業IT管理者は「1箇所」でMCP接続を管理できます。

しかし、その「1箇所」はMCPホスト（Claude Desktop, Cursor等）に紐づいています:

```
User → MCP Host (SSO) → Enterprise IdP → Token → MCP Server
            ↑
      ここでロックイン
```

**結果:**
- LLMホストを変更 → IdP連携全体を再構築
- 100個のコネクタ → 100個の認証再設定
- 移行コスト >> 継続コスト → ベンダーロックイン

詳細: [MCPist Position Statement: MCPの認証モデルとベンダーロックインリスク](https://github.com/shibaleo/go-mcp-dev/blob/main/docs/position-statement.md) *(予定)*

### The Solution: Server-Side Authentication

MCPistは認証をサーバー側（Token Broker）で管理:

```
User → MCP Host → MCP Server (Token Broker) → External APIs
                       ↑
                  認証をここで管理
```

**利点:**
- ✅ **LLMホスト変更 → 認証はそのまま** - Claude Desktop → Cursor への移行でも認証再設定不要
- ✅ **IdP連携もサーバー側** - ホスト依存なし
- ✅ **LLM選択の自由** - Claude/GPT/Gemini/Ollama を自由に切り替え（Phase 2）
- ✅ **真のポータビリティ** - MCPの本来の理念を実現

**Phase 1（現在）:** Claude Desktop/Cursor用のポータブルMCPサーバー
**Phase 2（予定）:** Host層を含めた完全ポータビリティ（MCPist Desktop）

参照:
- [ADR-005: RLSに依存しない認可設計](./dev/ADR-005-no-rls-dependency.md) - Token Brokerによるサーバー側認証
- [ADR-007: Host層を含めた完全ポータビリティの実現](./dev/ADR-007-host-layer-portability.md) - MCPist Desktopのビジョン
- [ADR-008: LLM API抽象化層によるクライアント側ポータビリティ](./dev/ADR-008-llm-api-abstraction.md) - LLM選択の自由

---

## 特徴

- **Go + SSE**: 軽量・高速なJSON-RPC 2.0 over Server-Sent Events
- **シングルテナント**: 固定シークレット認証、個人用途に最適化
- **$0運用**: Koyeb Free Tier + GitHub Actions pingでコールドスタート回避
- **オブザーバビリティ**: Grafana Cloud Lokiへのリアルタイムログ送信
- **84ツール**: 5モジュールで合計84のAPIツールを提供

## 対応モジュール

| モジュール | ツール数 | 説明 |
|-----------|---------|------|
| Supabase | 18 | Management API（プロジェクト管理、SQL実行、マイグレーション、ログ、ストレージ） |
| Notion | 15 | ページ・データベース・ブロック・コメント操作 |
| GitHub | 24 | リポジトリ、Issue、PR、Actions、コード検索 |
| Jira | 14 | Issue/Project操作、コメント、ワークログ |
| Confluence | 13 | Space/Page操作、CQL検索、ラベル |
| **合計** | **84** | |

## エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/health` | ヘルスチェック |
| POST | `/mcp` | JSON-RPC 2.0 over SSE |

## メタツール

LLMは2つのメタツールを通じて全84ツールにアクセスします（Lazy Loading）。

### get_module_schema
モジュールのツール定義を取得。各モジュールにつき1セッション1回のみ呼び出し。

```json
{
  "module": "supabase"
}
```

### call_module_tool
モジュールのツールを実行。

```json
{
  "module": "supabase",
  "tool_name": "run_query",
  "params": {
    "project_ref": "xxxxx",
    "query": "SELECT * FROM users LIMIT 10"
  }
}
```

## 各モジュールのツール一覧

### Supabase (18ツール)
- **Account**: list_organizations, list_projects, get_project
- **Database**: list_tables, run_query, list_migrations, apply_migration
- **Debugging**: get_logs, get_security_advisors, get_performance_advisors
- **Development**: get_project_url, get_api_keys, generate_typescript_types
- **Edge Functions**: list_edge_functions, get_edge_function
- **Storage**: list_storage_buckets, get_storage_config

### Notion (15ツール)
- **Search**: search
- **Pages**: get_page, get_page_content, create_page, update_page
- **Databases**: get_database, query_database
- **Blocks**: append_blocks, delete_block
- **Comments**: list_comments, add_comment
- **Users**: list_users, get_user, get_bot_user

### GitHub (24ツール)
- **User**: get_user
- **Repos**: list_repos, get_repo, list_branches, list_commits, get_file_content
- **Issues**: list_issues, get_issue, create_issue, update_issue, add_issue_comment
- **PRs**: list_prs, get_pr, create_pr, list_pr_commits, list_pr_files, list_pr_reviews
- **Search**: search_repos, search_code, search_issues
- **Actions**: list_workflows, list_workflow_runs, get_workflow_run

### Jira (14ツール)
- **User**: get_myself
- **Projects**: list_projects, get_project
- **Issues**: search, get_issue, create_issue, update_issue
- **Transitions**: get_transitions, transition_issue
- **Comments**: get_comments, add_comment
- **Worklogs**: get_worklogs, add_worklog

### Confluence (13ツール)
- **Spaces**: list_spaces, get_space
- **Pages**: get_pages, get_page, create_page, update_page, delete_page
- **Search**: search (CQL)
- **Comments**: get_page_comments, add_page_comment
- **Labels**: get_page_labels, add_page_label

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
| `NOTION_TOKEN` | Notion Integration token |
| `GITHUB_TOKEN` | GitHub Personal Access Token |
| `JIRA_DOMAIN` | Atlassian domain (例: xxx.atlassian.net) |
| `JIRA_EMAIL` | Atlassian account email |
| `JIRA_API_TOKEN` | Jira API token |
| `CONFLUENCE_DOMAIN` | Atlassian domain (Jiraと同じ) |
| `CONFLUENCE_EMAIL` | Atlassian account email |
| `CONFLUENCE_API_TOKEN` | Confluence API token |
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
│  │  - Notion API                               │                │
│  │  - GitHub API                               │                │
│  │  - Jira/Confluence API                      │                │
│  └─────────────────────────────────────────────┘                │
└─────────────────────────────────────────────────────────────────┘
```

### 想定ユーザー

- **対象**: 開発者本人（シングルテナント）
- **用途**: Claude Codeから自然言語で外部サービスを操作
- **例**:
  - 「go-mcp-demoプロジェクトのテーブル一覧を見せて」
  - 「usersテーブルにカラムを追加するSQLを実行して」
  - 「Notionのページを検索して」
  - 「GitHubのIssue一覧を取得して」
  - 「JiraでTODOのIssueを検索して」

## CI/CD

- **PR時**: ビルド・テスト実行
- **main push時**: Koyeb自動再デプロイ + ヘルスチェック + APIバージョンチェック
- **45分ごと**: GitHub Actionsからping（スリープ回避）

## APIバージョン管理

各モジュールは外部APIの公式バージョン文字列をそのまま記録する（semverではない）。

| モジュール | APIVersion | 形式 |
|-----------|------------|------|
| Supabase | `v1` | URLパス (`/v1/`) |
| Notion | `2022-06-28` | ヘッダー (`Notion-Version`) |
| GitHub | `2022-11-28` | ヘッダー (`X-GitHub-Api-Version`) |
| Jira | `3` | URLパス (`/rest/api/3`) |
| Confluence | `v2` | URLパス (`/wiki/api/v2`) |

### バージョンチェック

```bash
# ローカルで実行
go run ./cmd/version-check

# CIで自動実行（main push時）
# 各APIに実際にリクエストを送り、バージョン互換性を確認
```

- `TestedAt`: 最終動作確認日（手動更新）
- 不一致検出時: CIが失敗（exit code 2）

## ライセンス

Private
