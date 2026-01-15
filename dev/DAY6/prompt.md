# MCPist Admin UI 修正プロンプト

## 修正項目一覧

1. マイ接続ページ: 課金で無効なサービスを下部にまとめる
2. プロファイル一覧ページの追加
3. 全ページに戻るボタンを追加
4. サイドバーのスライダーデザイン修正
5. ProfilesPageのkey propエラー修正
6. Toolsページの表示を管理者/ユーザーで分ける
7. ダッシュボードのカードにリンクを追加

---

## 1. マイ接続ページの修正 (`components/my-connections/my-connections-content.tsx`)

### 変更内容

課金状態により無効なサービスを下部にまとめて表示する。

```
┌─────────────────────────────────────────────────────────────┐
│ マイ接続                                          [← 戻る]  │
│ サービスとの接続を管理します                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ 【利用可能なサービス】                                       │
│ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │
│ │ Google   │ │  Notion  │ │  GitHub  │ │  Google  │        │
│ │ Calendar │ │          │ │          │ │  Drive   │        │
│ │  接続済  │ │  接続済  │ │  未接続  │ │  未接続  │        │
│ └──────────┘ └──────────┘ └──────────┘ └──────────┘        │
│                                                              │
│ ─────────────────────────────────────────────────────────── │
│                                                              │
│ 【プランのアップグレードが必要】                Collapse可能  │
│ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ グレー │
│ │  Slack   │ │   Jira   │ │  freee   │ │  Zaim    │        │
│ │  🔒 PRO  │ │  🔒 PRO  │ │  🔒 MAX  │ │  🔒 MAX  │        │
│ └──────────┘ └──────────┘ └──────────┘ └──────────┘        │
└─────────────────────────────────────────────────────────────┘
```

### 実装要件

- 現在のプランで利用可能なサービスを上部に表示
- プランで利用不可のサービスは「プランのアップグレードが必要」セクションとして下部に表示
- 下部セクションはCollapsibleで折りたたみ可能
- グレーアウト表示 + プランバッジ（PRO/MAX）
- クリックするとアップグレード促進ダイアログを表示

---

## 2. プロファイル一覧ページの修正 (`app/(admin)/profiles/page.tsx`)

### 問題点

現在のProfilesページは権限マトリクスのみ表示しており、プロファイルの一覧・管理機能がない。

### 変更内容

プロファイル一覧と権限マトリクスをタブで切り替えられるようにする。

```
┌─────────────────────────────────────────────────────────────┐
│ プロファイル設定                                  [← 戻る]  │
│ プロファイルと権限を管理                                    │
├─────────────────────────────────────────────────────────────┤
│ [一覧] [権限マトリクス]                                     │
├─────────────────────────────────────────────────────────────┤
│                                    （一覧タブ選択時）        │
│ [+ 新規プロファイル]                                        │
│                                                              │
│ ┌─────────────────────────────────────────────────────────┐│
│ │ 開発者標準                                               ││
│ │ 開発者向けの標準的な権限セット                           ││
│ │ 適用ロール: 開発者                    [編集] [削除]      ││
│ └─────────────────────────────────────────────────────────┘│
│ ┌─────────────────────────────────────────────────────────┐│
│ │ 閲覧のみ                                                 ││
│ │ 読み取り専用の最小権限セット                             ││
│ │ 適用ロール: 閲覧者, 外部パートナー    [編集] [削除]      ││
│ └─────────────────────────────────────────────────────────┘│
│ ┌─────────────────────────────────────────────────────────┐│
│ │ フルアクセス                                             ││
│ │ 全てのツールへのアクセス権限                             ││
│ │ 適用ロール: 管理者                    [編集] [削除]      ││
│ └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

### 実装要件

- Tabsコンポーネントで「一覧」と「権限マトリクス」を切り替え
- 一覧タブ：プロファイルのカード一覧 + 新規作成・編集・削除機能
- 権限マトリクスタブ：既存のマトリクス表示
- プロファイル編集はDialogで実装

---

## 3. 全ページに戻るボタンを追加

### 変更対象ページ

- `/tools`
- `/tools/[module]`
- `/users`
- `/roles`
- `/profiles`
- `/service-auth`
- `/billing`
- `/requests`
- `/logs`
- `/tokens`
- `/my/connections`
- `/my/mcp-connection`
- `/my/preferences`

### 実装方法

各ページのヘッダー部分に戻るボタンを追加：

```tsx
import { ArrowLeft } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useRouter } from "next/navigation"

// ページヘッダー部分
<div className="flex items-center gap-4">
  <Button variant="ghost" size="icon" onClick={() => router.back()}>
    <ArrowLeft className="h-5 w-5" />
  </Button>
  <div>
    <h1 className="text-2xl font-bold text-foreground">ページタイトル</h1>
    <p className="text-muted-foreground mt-1">説明文</p>
  </div>
</div>
```

### 注意

- Dashboardページには戻るボタン不要（トップページのため）
- 戻るボタンは `router.back()` を使用

---

## 4. サイドバーのスライダーデザイン修正 (`components/sidebar.tsx`)

### 問題点

サイドバーの折りたたみボタン（ChevronLeft/ChevronRight）がデザインと合っていない。

### 変更内容

折りたたみボタンをより目立たないシンプルなデザインに変更：

```tsx
// 現在
<Button variant="ghost" size="sm" className={cn("w-full", collapsed && "justify-center px-2")}>

// 変更後
<Button
  variant="ghost"
  size="sm"
  className={cn(
    "w-full h-8 text-muted-foreground hover:text-foreground",
    collapsed && "justify-center px-2"
  )}
>
```

または、サイドバー端に細いバー/ハンドルとして実装：

```tsx
// サイドバーの右端に配置
<div
  className="absolute right-0 top-1/2 -translate-y-1/2 w-1 h-16 bg-border rounded-full cursor-pointer hover:bg-primary/50 transition-colors"
  onClick={() => onCollapsedChange?.(!collapsed)}
/>
```

---

## 5. ProfilesPageのkey propエラー修正 (`app/(admin)/profiles/page.tsx`)

### 問題点

```
Each child in a list should have a unique "key" prop.
```

`<>` (Fragment) に key が必要。

### 修正方法

```tsx
// 修正前
return (
  <>
    <TableRow key={service.id} ...>
    ...
  </>
)

// 修正後
return (
  <React.Fragment key={service.id}>
    <TableRow ...>
    ...
  </React.Fragment>
)
```

または、mapのインデックスを使用してFragmentにkeyを付与：

```tsx
{services.map((service, serviceIndex) => (
  <React.Fragment key={`service-${service.id}`}>
    <TableRow className="cursor-pointer hover:bg-muted/50" onClick={() => toggleService(service.id)}>
      ...
    </TableRow>
    {expandedServices.includes(service.id) && (
      moduleDetails[service.id]?.tools.map((tool) => (
        <TableRow key={`${service.id}-${tool.id}`}>
          ...
        </TableRow>
      ))
    )}
  </React.Fragment>
))}
```

---

## 6. Toolsページの表示を管理者/ユーザーで分ける (`app/(admin)/tools/page.tsx`)

### 問題点

「組織で利用可能なサービスを管理」という説明文や、プラン情報は管理者向けの内容であり、一般ユーザーには不適切。

### 変更内容

`useAuth()` の `isAdmin` を使って表示を切り替える：

**管理者向け表示：**
```
┌─────────────────────────────────────────────────────────────┐
│ サービス管理                                      [← 戻る]  │
│ 組織で利用可能なサービスを管理                               │
├─────────────────────────────────────────────────────────────┤
│ 現在のプラン: Free        ユーザー: 5/5名    [プランを変更]  │
├─────────────────────────────────────────────────────────────┤
│ （プラン別のサービス一覧 + アップグレード導線）              │
└─────────────────────────────────────────────────────────────┘
```

**一般ユーザー向け表示：**
```
┌─────────────────────────────────────────────────────────────┐
│ ツール                                            [← 戻る]  │
│ 利用可能なツールを確認                                      │
├─────────────────────────────────────────────────────────────┤
│ （自分のロールで利用可能なサービスのみ表示）                │
│ （プラン情報やアップグレードボタンは非表示）                │
└─────────────────────────────────────────────────────────────┘
```

### 実装要件

```tsx
const { isAdmin } = useAuth()

return (
  <div>
    <h1>{isAdmin ? "サービス管理" : "ツール"}</h1>
    <p>{isAdmin ? "組織で利用可能なサービスを管理" : "利用可能なツールを確認"}</p>

    {isAdmin && (
      // プラン情報バー
      <div>現在のプラン: {plan} ...</div>
    )}

    {/* サービス一覧 */}
    {isAdmin ? (
      // 管理者：プラン別にグループ化して表示
    ) : (
      // ユーザー：利用可能なサービスのみ表示
    )}
  </div>
)
```

---

## 7. ダッシュボードのカードにリンクを追加 (`app/(admin)/dashboard/page.tsx`)

### 変更内容

統計カードをクリック可能にして、関連ページへ遷移させる。

```tsx
import Link from "next/link"

// カードをLinkでラップ
<Link href="/tools">
  <Card className="cursor-pointer transition-all hover:border-primary/50 hover:shadow-md">
    <CardHeader>
      <CardTitle>接続済みサービス</CardTitle>
    </CardHeader>
    <CardContent>
      <div className="text-3xl font-bold">6</div>
    </CardContent>
  </Card>
</Link>
```

### リンク先マッピング

| カード | リンク先 |
|--------|----------|
| 接続済みサービス | `/tools` |
| アクティブユーザー | `/users` |
| API呼び出し（今日） | `/logs` |
| 保留中のリクエスト | `/requests` |
| プラン情報 | `/billing` |

### スタイル

- `cursor-pointer` を追加
- `hover:border-primary/50` でボーダー色変更
- `hover:shadow-md` でシャドウ追加
- `transition-all` でスムーズなアニメーション

---

## 共通事項

### インポート追加が必要なもの

```tsx
import React from "react"  // React.Fragment用
import { ArrowLeft } from "lucide-react"  // 戻るボタン用
import { useRouter } from "next/navigation"  // 戻るボタン用
import Link from "next/link"  // ダッシュボードカード用
```

### デザイン指針

- 戻るボタン: `variant="ghost"` `size="icon"` で控えめに
- ホバー効果: `hover:border-primary/50` `hover:shadow-md` で統一
- プランバッジ:
  - PRO: `bg-blue-500/20 text-blue-600 border-blue-500/30`
  - MAX: `bg-purple-500/20 text-purple-600 border-purple-500/30`
