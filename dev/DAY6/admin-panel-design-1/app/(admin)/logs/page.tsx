import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

const logs = [
  {
    id: "1",
    timestamp: "2026-01-15 10:32:15",
    user: "山田 太郎",
    action: "サービス連携",
    target: "Google Calendar",
    status: "success",
  },
  {
    id: "2",
    timestamp: "2026-01-15 10:28:03",
    user: "佐藤 花子",
    action: "ロール変更",
    target: "開発者 → 管理者",
    status: "success",
  },
  {
    id: "3",
    timestamp: "2026-01-15 09:45:22",
    user: "鈴木 一郎",
    action: "ログイン",
    target: "-",
    status: "success",
  },
  {
    id: "4",
    timestamp: "2026-01-15 09:12:08",
    user: "田中 美咲",
    action: "サービス解除",
    target: "Slack",
    status: "warning",
  },
  {
    id: "5",
    timestamp: "2026-01-14 18:30:45",
    user: "高橋 健太",
    action: "API呼び出し",
    target: "GitHub API",
    status: "error",
  },
]

export default function LogsPage() {
  const getStatusBadge = (status: string) => {
    switch (status) {
      case "success":
        return (
          <Badge variant="outline" className="bg-success/20 text-success border-success/30">
            成功
          </Badge>
        )
      case "warning":
        return (
          <Badge variant="outline" className="bg-warning/20 text-warning border-warning/30">
            警告
          </Badge>
        )
      case "error":
        return (
          <Badge variant="outline" className="bg-destructive/20 text-destructive border-destructive/30">
            エラー
          </Badge>
        )
      default:
        return <Badge variant="secondary">{status}</Badge>
    }
  }

  return (
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Logs</h1>
        <p className="text-muted-foreground mt-1">システムアクティビティログ</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">最近のアクティビティ</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          <div className="overflow-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[180px]">日時</TableHead>
                  <TableHead>ユーザー</TableHead>
                  <TableHead>アクション</TableHead>
                  <TableHead className="hidden md:table-cell">対象</TableHead>
                  <TableHead className="text-right">ステータス</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {logs.map((log) => (
                  <TableRow key={log.id}>
                    <TableCell className="font-mono text-sm text-muted-foreground">{log.timestamp}</TableCell>
                    <TableCell>{log.user}</TableCell>
                    <TableCell>{log.action}</TableCell>
                    <TableCell className="hidden md:table-cell text-muted-foreground">{log.target}</TableCell>
                    <TableCell className="text-right">{getStatusBadge(log.status)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
