"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { authApi, setToken } from "@/lib/api"
import { LayoutGrid } from "lucide-react"

export default function LoginPage() {
  const router = useRouter()
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)

    try {
      const data = await authApi.login({ username, password })
      setToken(data.token)
      router.push("/dashboard")
    } catch (err) {
      setError(err instanceof Error ? err.message : "登录失败")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-amber-50 via-orange-50 to-rose-50">
      {/* 背景装饰 */}
      <div className="absolute inset-0 bg-[url('data:image/svg+xml,%3Csvg width=%2260%22 height=%2260%22 viewBox=%220 0 60 60%22 fill=%22none%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Ccircle cx=%2230%22 cy=%2230%22 r=%221.5%22 fill=%22rgba(251,146,60,0.15)%22/%3E%3C/svg%3E')]" />

      <Card className="w-full max-w-md mx-4 shadow-2xl border border-orange-100 bg-white/80 backdrop-blur-sm">
        <CardHeader className="space-y-4 pb-6">
          {/* Logo区域 */}
          <div className="flex justify-center">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-orange-500 to-amber-500 flex items-center justify-center shadow-lg shadow-orange-200">
              <LayoutGrid className="w-8 h-8 text-white" />
            </div>
          </div>
          <div className="space-y-2">
            <CardTitle className="text-2xl font-bold text-center text-gray-800">
              DualTab 后台管理
            </CardTitle>
            <CardDescription className="text-center text-gray-500">
              请输入管理员账号登录
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-2">
              <Label htmlFor="username" className="text-gray-700">用户名</Label>
              <Input
                id="username"
                type="text"
                placeholder="请输入用户名"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
                className="h-11 bg-white border-gray-200 text-gray-900 placeholder:text-gray-400 focus:border-orange-400 focus:ring-orange-400 transition-colors"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password" className="text-gray-700">密码</Label>
              <Input
                id="password"
                type="password"
                placeholder="请输入密码"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="h-11 bg-white border-gray-200 text-gray-900 placeholder:text-gray-400 focus:border-orange-400 focus:ring-orange-400 transition-colors"
              />
            </div>
            {error && (
              <div className="text-sm text-red-500 text-center bg-red-50 py-2 rounded-lg">{error}</div>
            )}
            <Button
              type="submit"
              className="w-full h-11 bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white font-medium shadow-md hover:shadow-lg transition-all"
              disabled={loading}
            >
              {loading ? "登录中..." : "登录"}
            </Button>
          </form>
          <div className="mt-6 text-center text-sm text-gray-400">
            默认账号: admin / admin123
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
