"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"
import { useAuth } from "@/lib/auth-context"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  LayoutDashboard,
  Wrench,
  Users,
  Shield,
  FileText,
  LogOut,
  Settings,
  ChevronLeft,
  ChevronRight,
  Key,
  ClipboardList,
  Layers,
} from "lucide-react"

const navItems = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard, adminOnly: false },
  { href: "/tools", label: "Tools", icon: Wrench, adminOnly: false },
  { href: "/tokens", label: "API Tokens", icon: Key, adminOnly: false },
  { href: "/users", label: "Users", icon: Users, adminOnly: true },
  { href: "/roles", label: "Roles", icon: Shield, adminOnly: true },
  { href: "/profiles", label: "Profiles", icon: Layers, adminOnly: true },
  { href: "/requests", label: "Requests", icon: ClipboardList, adminOnly: true },
  { href: "/logs", label: "Logs", icon: FileText, adminOnly: true },
]

interface SidebarProps {
  collapsed?: boolean
  onCollapsedChange?: (collapsed: boolean) => void
}

export function Sidebar({ collapsed = false, onCollapsedChange }: SidebarProps) {
  const pathname = usePathname()
  const { user, isAdmin } = useAuth()

  const filteredNavItems = navItems.filter((item) => !item.adminOnly || isAdmin)

  return (
    <aside
      className={cn(
        "flex flex-col h-full bg-sidebar border-r border-sidebar-border transition-all duration-300",
        collapsed ? "w-16" : "w-64",
      )}
    >
      {/* Logo */}
      <div className="flex items-center h-16 px-4 border-b border-sidebar-border">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
            <span className="text-primary-foreground font-bold text-sm">M</span>
          </div>
          {!collapsed && <span className="font-semibold text-sidebar-foreground">MCP Server</span>}
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 p-3 space-y-1">
        {filteredNavItems.map((item) => {
          const isActive = pathname === item.href || pathname.startsWith(item.href + "/")
          return (
            <Link key={item.href} href={item.href}>
              <Button
                variant={isActive ? "secondary" : "ghost"}
                className={cn(
                  "w-full justify-start gap-3 h-10",
                  isActive && "bg-sidebar-accent text-sidebar-accent-foreground",
                  collapsed && "justify-center px-2",
                )}
              >
                <item.icon className="h-5 w-5 shrink-0" />
                {!collapsed && <span>{item.label}</span>}
              </Button>
            </Link>
          )
        })}
      </nav>

      {/* Collapse Button */}
      <div className="px-3 pb-2">
        <Button
          variant="ghost"
          size="sm"
          className={cn("w-full", collapsed && "justify-center px-2")}
          onClick={() => onCollapsedChange?.(!collapsed)}
        >
          {collapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <>
              <ChevronLeft className="h-4 w-4 mr-2" />
              <span>折りたたむ</span>
            </>
          )}
        </Button>
      </div>

      {/* User Profile */}
      <div className="p-3 border-t border-sidebar-border">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className={cn("w-full h-auto p-2", collapsed ? "justify-center" : "justify-start")}>
              <Avatar className="h-8 w-8">
                <AvatarImage src={user?.avatar || "/placeholder.svg"} />
                <AvatarFallback className="bg-primary text-primary-foreground text-xs">
                  {user?.name?.slice(0, 2) || "U"}
                </AvatarFallback>
              </Avatar>
              {!collapsed && (
                <div className="ml-3 text-left">
                  <p className="text-sm font-medium text-sidebar-foreground truncate max-w-[140px]">{user?.name}</p>
                  <p className="text-xs text-muted-foreground truncate max-w-[140px]">{user?.email}</p>
                </div>
              )}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <DropdownMenuItem>
              <Settings className="mr-2 h-4 w-4" />
              設定
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-destructive">
              <LogOut className="mr-2 h-4 w-4" />
              ログアウト
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </aside>
  )
}
