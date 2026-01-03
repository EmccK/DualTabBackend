"use client"

import { useState } from "react"
import Link from "next/link"
import { usePathname, useRouter } from "next/navigation"
import { LayoutGrid, Image, Search, FolderOpen, LogOut, KeyRound, ImageIcon, Settings } from "lucide-react"
import { clearToken, authApi } from "@/lib/api"
import { cn } from "@/lib/utils"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useToast } from "@/hooks/use-toast"

const navItems = [
  { href: "/dashboard", label: "仪表盘", icon: LayoutGrid },
  { href: "/dashboard/icons", label: "图标管理", icon: Image },
  { href: "/dashboard/categories", label: "分类管理", icon: FolderOpen },
  { href: "/dashboard/search-engines", label: "搜索引擎", icon: Search },
  { href: "/dashboard/wallpapers", label: "壁纸管理", icon: ImageIcon },
  { href: "/dashboard/settings", label: "系统配置", icon: Settings },
]

export function Sidebar() {
  const { toast } = useToast()
  const pathname = usePathname()
  const router = useRouter()
  const [showPasswordDialog, setShowPasswordDialog] = useState(false)
  const [oldPassword, setOldPassword] = useState("")
  const [newPassword, setNewPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const handleLogout = () => {
    clearToken()
    router.push("/")
  }

  const handleChangePassword = async () => {
    setError("")

    if (newPassword !== confirmPassword) {
      setError("两次输入的密码不一致")
      return
    }

    if (newPassword.length < 6) {
      setError("新密码至少6位")
      return
    }

    setLoading(true)
    try {
      await authApi.changePassword({
        old_password: oldPassword,
        new_password: newPassword,
      })
      setShowPasswordDialog(false)
      setOldPassword("")
      setNewPassword("")
      setConfirmPassword("")
      toast({
        title: "修改成功",
        description: "密码已成功修改",
        variant: "success",
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : "修改失败")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex h-screen w-64 flex-col bg-white border-r border-gray-200">
      {/* Logo区域 */}
      <div className="flex h-16 items-center gap-3 px-6 border-b border-gray-100">
        <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-orange-500 to-amber-500 flex items-center justify-center shadow-md shadow-orange-200">
          <LayoutGrid className="w-5 h-5 text-white" />
        </div>
        <h1 className="text-lg font-semibold text-gray-800">DualTab</h1>
      </div>

      {/* 导航菜单 */}
      <nav className="flex-1 p-4 space-y-1">
        {navItems.map((item) => {
          const isActive = pathname === item.href
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium transition-all duration-200",
                isActive
                  ? "bg-gradient-to-r from-orange-500 to-amber-500 text-white shadow-md shadow-orange-200"
                  : "text-gray-600 hover:bg-orange-50 hover:text-orange-600"
              )}
            >
              <item.icon className="h-5 w-5" />
              {item.label}
            </Link>
          )
        })}
      </nav>

      {/* 底部操作区 */}
      <div className="border-t border-gray-100 p-4 space-y-1">
        <button
          onClick={() => setShowPasswordDialog(true)}
          className="flex w-full items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium text-gray-600 transition-all duration-200 hover:bg-orange-50 hover:text-orange-600"
        >
          <KeyRound className="h-5 w-5" />
          修改密码
        </button>
        <button
          onClick={handleLogout}
          className="flex w-full items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium text-gray-600 transition-all duration-200 hover:bg-red-50 hover:text-red-500"
        >
          <LogOut className="h-5 w-5" />
          退出登录
        </button>
      </div>

      {/* 修改密码对话框 */}
      <Dialog open={showPasswordDialog} onOpenChange={setShowPasswordDialog}>
        <DialogContent className="bg-white">
          <DialogHeader>
            <DialogTitle className="text-gray-900">修改密码</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <Label htmlFor="oldPassword" className="text-gray-700">当前密码</Label>
              <Input
                id="oldPassword"
                type="password"
                value={oldPassword}
                onChange={(e) => setOldPassword(e.target.value)}
                placeholder="请输入当前密码"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="newPassword" className="text-gray-700">新密码</Label>
              <Input
                id="newPassword"
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="请输入新密码（至少6位）"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirmPassword" className="text-gray-700">确认新密码</Label>
              <Input
                id="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="请再次输入新密码"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            {error && <p className="text-sm text-red-500">{error}</p>}
            <Button
              className="w-full bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white"
              onClick={handleChangePassword}
              disabled={loading}
            >
              {loading ? "提交中..." : "确认修改"}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
