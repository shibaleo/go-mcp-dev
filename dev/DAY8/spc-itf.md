# MCPist インターフェース仕様書（spc-itf）

## ドキュメント管理情報

| 項目 | 値 |
|------|-----|
| Status | `draft` |
| Version | v1.0 (DAY8) |
| Note | コンポーネント間インターフェース定義 |

---

## 概要

本ドキュメントは、spc-sys.mdで定義されたコンポーネント間のインターフェースを規定する。

MCP Serverは内部コンポーネント（Auth Middleware, MCP Handler, Module Registry, Modules）の集合であり、やり取りの主体とはならない。各内部コンポーネントが直接外部とやり取りする。

### コンポーネント一覧

| # | コンポーネント | 略称 | 備考 |
|---|---------------|------|------|
| 1 | MCP Client | CLI | 実装範囲外 |
| 2 | Auth Server | AUTH | |
| 3 | Auth Middleware | MW | MCP Server内部 |
| 4 | MCP Handler | HDL | MCP Server内部 |
| 5 | Module Registry | REG | MCP Server内部 |
| 6 | Modules | MOD | MCP Server内部 |
| 7 | Entitlement Store | ENT | |
| 8 | Token Vault | TVL | |
| 9 | User Console | CON | |
| 10 | External API Server | EXT | 実装範囲外 |

---

## 1. MCP Client（CLI）

| 相手 | 方向 | やり取り |
|------|------|----------|
| Auth Server | CLI → AUTH | OAuth 2.1認証フロー（認可コード取得、トークン交換） |
| MCP Server | CLI → SRV | MCP Protocol（JSON-RPC over SSE） |
| Auth Middleware | - | 直接やり取りなし（MCP Server経由） |
| MCP Handler | - | 直接やり取りなし（MCP Server経由） |
| Module Registry | - | 直接やり取りなし（MCP Server経由） |
| Modules | - | 直接やり取りなし（MCP Server経由） |
| Entitlement Store | - | 直接やり取りなし |
| Token Vault | - | 直接やり取りなし |
| User Console | - | 直接やり取りなし |
| External API Server | - | 直接やり取りなし |

---

## 2. Auth Server（AUTH）

| 相手                  | 方向         | やり取り                     |
| ------------------- | ---------- | ------------------------ |
| MCP Client          | AUTH ← CLI | OAuth 2.1認証リクエスト受付       |
| MCP Server          | AUTH ← SRV | JWKS公開鍵の提供（JWT検証用）       |
| Auth Middleware     | AUTH ← MW  | JWKS取得（公開鍵キャッシュ）         |
| MCP Handler         | -          | 直接やり取りなし                 |
| Module Registry     | -          | 直接やり取りなし                 |
| Modules             | -          | 直接やり取りなし                 |
| Entitlement Store   | AUTH → ENT | ユーザー情報の参照・作成（OAuth登録時）   |
| Token Vault         | -          | 直接やり取りなし                 |
| User Console        | AUTH ← CON | OAuth 2.1認証フロー（ユーザーログイン） |
| External API Server | -          | 直接やり取りなし                 |

---

## 3. MCP Server（SRV）

| 相手 | 方向 | やり取り |
|------|------|----------|
| MCP Client | SRV ← CLI | MCP Protocolリクエスト受付 |
| Auth Server | SRV → AUTH | JWKS取得（JWT検証用公開鍵） |
| Auth Middleware | SRV ↔ MW | 内部コンポーネント（リクエスト処理委譲） |
| MCP Handler | SRV ↔ HDL | 内部コンポーネント（リクエスト処理委譲） |
| Module Registry | SRV ↔ REG | 内部コンポーネント（リクエスト処理委譲） |
| Modules | SRV ↔ MOD | 内部コンポーネント（リクエスト処理委譲） |
| Entitlement Store | SRV → ENT | 権限情報の参照 |
| Token Vault | SRV → TVL | トークンの取得 |
| User Console | - | 直接やり取りなし |
| External API Server | - | 直接やり取りなし（Modules経由） |

---

## 4. Auth Middleware（MW）

| 相手                  | 方向        | やり取り                                  |
| ------------------- | --------- | ------------------------------------- |
| MCP Client          | -         | 直接やり取りなし（MCP Server経由）                |
| Auth Server         | MW → AUTH | JWKS取得（公開鍵キャッシュ、定期更新）                 |
| MCP Server          | MW ↔ SRV  | 内部コンポーネント（親コンポーネント）                   |
| MCP Handler         | MW → HDL  | 認証済みリクエストの転送（user_id付きcontext）        |
| Module Registry     | -         | 直接やり取りなし                              |
| Modules             | -         | 直接やり取りなし                              |
| Entitlement Store   | MW → ENT  | アカウント状態の確認（active/suspended/disabled） |
| Token Vault         | -         | 直接やり取りなし                              |
| User Console        | -         | 直接やり取りなし                              |
| External API Server | -         | 直接やり取りなし                              |

---

## 5. MCP Handler（HDL）

| 相手 | 方向 | やり取り |
|------|------|----------|
| MCP Client | - | 直接やり取りなし（MCP Server経由） |
| Auth Server | - | 直接やり取りなし |
| MCP Server | HDL ↔ SRV | 内部コンポーネント（親コンポーネント） |
| Auth Middleware | HDL ← MW | 認証済みリクエストの受信 |
| Module Registry | HDL → REG | メタツール呼び出し（get_module_schema, call, batch） |
| Modules | - | 直接やり取りなし（Module Registry経由） |
| Entitlement Store | HDL → ENT | 権限チェック（Permission Gate/Filter） |
| Token Vault | - | 直接やり取りなし |
| User Console | - | 直接やり取りなし |
| External API Server | - | 直接やり取りなし |

---

## 6. Module Registry（REG）

| 相手 | 方向 | やり取り |
|------|------|----------|
| MCP Client | - | 直接やり取りなし |
| Auth Server | - | 直接やり取りなし |
| MCP Server | REG ↔ SRV | 内部コンポーネント（親コンポーネント） |
| Auth Middleware | - | 直接やり取りなし |
| MCP Handler | REG ← HDL | メタツールリクエスト受信 |
| Modules | REG → MOD | ツール実行委譲、スキーマ取得 |
| Entitlement Store | REG → ENT | Permission Filter（許可ツールのフィルタリング） |
| Token Vault | - | 直接やり取りなし（Modules経由） |
| User Console | - | 直接やり取りなし |
| External API Server | - | 直接やり取りなし |

---

## 7. Modules（MOD）

| 相手 | 方向 | やり取り |
|------|------|----------|
| MCP Client | - | 直接やり取りなし |
| Auth Server | - | 直接やり取りなし |
| MCP Server | MOD ↔ SRV | 内部コンポーネント（親コンポーネント） |
| Auth Middleware | - | 直接やり取りなし |
| MCP Handler | - | 直接やり取りなし（Module Registry経由） |
| Module Registry | MOD ← REG | ツール実行リクエスト受信 |
| Entitlement Store | - | 直接やり取りなし |
| Token Vault | MOD → TVL | OAuthトークン取得（user_id + service） |
| User Console | - | 直接やり取りなし |
| External API Server | MOD → EXT | 外部API呼び出し（HTTPS + Bearer Token） |

---

## 8. Entitlement Store（ENT）

| 相手 | 方向 | やり取り |
|------|------|----------|
| MCP Client | - | 直接やり取りなし |
| Auth Server | ENT ← AUTH | ユーザー情報の参照・作成 |
| MCP Server | ENT ← SRV | 権限情報の参照 |
| Auth Middleware | ENT ← MW | アカウント状態の参照 |
| MCP Handler | ENT ← HDL | 権限情報の参照 |
| Module Registry | ENT ← REG | 権限情報の参照 |
| Modules | - | 直接やり取りなし |
| Token Vault | - | 直接やり取りなし |
| User Console | ENT ← CON | 設定の書き込み（課金、ツール有効/無効） |
| External API Server | - | 直接やり取りなし |

---

## 9. Token Vault（TVL）

| 相手 | 方向 | やり取り |
|------|------|----------|
| MCP Client | - | 直接やり取りなし |
| Auth Server | - | 直接やり取りなし |
| MCP Server | TVL ← SRV | トークン取得リクエスト |
| Auth Middleware | - | 直接やり取りなし |
| MCP Handler | - | 直接やり取りなし |
| Module Registry | - | 直接やり取りなし |
| Modules | TVL ← MOD | トークン取得リクエスト（user_id + service） |
| Entitlement Store | - | 直接やり取りなし |
| User Console | TVL ← CON | OAuthトークン登録（OAuth連携完了時） |
| External API Server | TVL → EXT | トークンリフレッシュ（OAuth refresh_token使用） |

---

## 10. User Console（CON）

| 相手                  | 方向         | やり取り                |
| ------------------- | ---------- | ------------------- |
| MCP Client          | -          | 直接やり取りなし            |
| Auth Server         | CON → AUTH | ユーザー認証（ログイン）        |
| MCP Server          | -          | 直接やり取りなし            |
| Auth Middleware     | -          | 直接やり取りなし            |
| MCP Handler         | -          | 直接やり取りなし            |
| Module Registry     | -          | 直接やり取りなし            |
| Modules             | -          | 直接やり取りなし            |
| Entitlement Store   | CON → ENT  | 設定の読み書き             |
| Token Vault         | CON → TVL  | OAuthトークン登録         |
| External API Server | CON → EXT  | OAuth認可フロー（認可コード取得） |

---

## 11. External API Server（EXT）

| 相手 | 方向 | やり取り |
|------|------|----------|
| MCP Client | - | 直接やり取りなし |
| Auth Server | - | 直接やり取りなし |
| MCP Server | - | 直接やり取りなし |
| Auth Middleware | - | 直接やり取りなし |
| MCP Handler | - | 直接やり取りなし |
| Module Registry | - | 直接やり取りなし |
| Modules | EXT ← MOD | API呼び出し受付（HTTPS） |
| Entitlement Store | - | 直接やり取りなし |
| Token Vault | EXT ← TVL | トークンリフレッシュリクエスト受付 |
| User Console | EXT ← CON | OAuth認可フロー受付 |

---

## 関連ドキュメント

| ドキュメント | 内容 |
|-------------|------|
| [spc-sys.md](./spc-sys.md) | システム仕様書（コンポーネント定義） |
| [dsn-module-registry.md](../DAY7/dsn-module-registry.md) | Module Registry設計 |
| [dsn-permission-system.md](../DAY7/dsn-permission-system.md) | 権限システム設計 |
