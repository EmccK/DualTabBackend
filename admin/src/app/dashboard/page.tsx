"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { iconApi, categoryApi, searchEngineApi } from "@/lib/api"
import { Image, FolderOpen, Search, ArrowRight } from "lucide-react"
import Link from "next/link"

export default function DashboardPage() {
  const [stats, setStats] = useState({
    icons: 0,
    categories: 0,
    searchEngines: 0,
  })

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [iconsRes, categoriesRes, searchEnginesRes] = await Promise.all([
          iconApi.list({ page: 1, size: 1 }),
          categoryApi.list(),
          searchEngineApi.list(),
        ])
        setStats({
          icons: iconsRes.total,
          categories: categoriesRes.list.length,
          searchEngines: searchEnginesRes.list.length,
        })
      } catch (error) {
        console.error("获取统计数据失败:", error)
      }
    }
    fetchStats()
  }, [])

  // 统计卡片配置
  const statCards = [
    {
      title: "推荐图标",
      value: stats.icons,
      description: "已添加的推荐书签",
      icon: Image,
      href: "/dashboard/icons",
      gradient: "from-orange-500 to-amber-500",
      bgGradient: "from-orange-500/10 to-amber-500/10",
    },
    {
      title: "分类数量",
      value: stats.categories,
      description: "图标分类总数",
      icon: FolderOpen,
      href: "/dashboard/categories",
      gradient: "from-rose-500 to-pink-500",
      bgGradient: "from-rose-500/10 to-pink-500/10",
    },
    {
      title: "搜索引擎",
      value: stats.searchEngines,
      description: "可用搜索引擎数量",
      icon: Search,
      href: "/dashboard/search-engines",
      gradient: "from-amber-500 to-yellow-500",
      bgGradient: "from-amber-500/10 to-yellow-500/10",
    },
  ]

  return (
    <div className="space-y-8">
      {/* 页面标题 */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">仪表盘</h1>
        <p className="text-gray-500 mt-2">欢迎使用 DualTab 后台管理系统</p>
      </div>

      {/* 统计卡片 */}
      <div className="grid gap-6 md:grid-cols-3">
        {statCards.map((card) => (
          <Link key={card.title} href={card.href}>
            <Card className="group relative overflow-hidden bg-white border border-gray-200 shadow-sm hover:shadow-md transition-all duration-300 cursor-pointer">
              <CardHeader className="relative flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium text-gray-600">{card.title}</CardTitle>
                <div className={`w-10 h-10 rounded-xl bg-gradient-to-br ${card.gradient} flex items-center justify-center shadow-md`}>
                  <card.icon className="h-5 w-5 text-white" />
                </div>
              </CardHeader>
              <CardContent className="relative">
                <div className="text-3xl font-bold text-gray-900">{card.value}</div>
                <p className="text-sm text-gray-500 mt-1">{card.description}</p>
                <div className="flex items-center gap-1 mt-3 text-sm font-medium text-gray-400 group-hover:text-orange-500 transition-colors">
                  查看详情
                  <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
                </div>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      {/* 快速开始 */}
      <Card className="bg-white border border-gray-200 shadow-sm">
        <CardHeader>
          <CardTitle className="text-lg font-semibold text-gray-900">快速开始</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <Link href="/dashboard/categories" className="group flex items-center gap-4 p-4 rounded-xl bg-gray-50 hover:bg-orange-50 border border-gray-100 hover:border-orange-200 transition-all">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-orange-500 to-amber-500 flex items-center justify-center shadow-md">
                <span className="text-lg font-bold text-white">1</span>
              </div>
              <div className="flex-1">
                <p className="font-medium text-gray-900">分类管理</p>
                <p className="text-sm text-gray-500">管理图标分类</p>
              </div>
              <ArrowRight className="w-5 h-5 text-gray-400 group-hover:text-orange-500 group-hover:translate-x-1 transition-all" />
            </Link>
            <Link href="/dashboard/icons" className="group flex items-center gap-4 p-4 rounded-xl bg-gray-50 hover:bg-orange-50 border border-gray-100 hover:border-orange-200 transition-all">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-orange-500 to-amber-500 flex items-center justify-center shadow-md">
                <span className="text-lg font-bold text-white">2</span>
              </div>
              <div className="flex-1">
                <p className="font-medium text-gray-900">图标管理</p>
                <p className="text-sm text-gray-500">添加推荐书签</p>
              </div>
              <ArrowRight className="w-5 h-5 text-gray-400 group-hover:text-orange-500 group-hover:translate-x-1 transition-all" />
            </Link>
            <Link href="/dashboard/search-engines" className="group flex items-center gap-4 p-4 rounded-xl bg-gray-50 hover:bg-orange-50 border border-gray-100 hover:border-orange-200 transition-all">
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-orange-500 to-amber-500 flex items-center justify-center shadow-md">
                <span className="text-lg font-bold text-white">3</span>
              </div>
              <div className="flex-1">
                <p className="font-medium text-gray-900">搜索引擎</p>
                <p className="text-sm text-gray-500">配置搜索引擎</p>
              </div>
              <ArrowRight className="w-5 h-5 text-gray-400 group-hover:text-orange-500 group-hover:translate-x-1 transition-all" />
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
