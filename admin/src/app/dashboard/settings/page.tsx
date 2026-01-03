"use client"

import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { configApi, type SystemConfig } from "@/lib/api"
import { Save, Cloud, Search, ImageIcon } from "lucide-react"
import { useToast } from "@/hooks/use-toast"

// 配置项定义
const CONFIG_ITEMS = [
  {
    key: "weather_api_key",
    label: "天气 API Key",
    description: "和风天气或 OpenWeather 的 API Key",
    type: "text",
    icon: Cloud,
    group: "weather",
  },
  {
    key: "weather_api_type",
    label: "天气 API 类型",
    description: "选择天气数据来源",
    type: "select",
    options: [
      { value: "qweather", label: "和风天气" },
      { value: "openweather", label: "OpenWeather" },
    ],
    icon: Cloud,
    group: "weather",
  },
  {
    key: "search_suggest_on",
    label: "搜索建议代理",
    description: "启用后，扩展将通过后端代理获取 Google 搜索建议",
    type: "switch",
    icon: Search,
    group: "search",
  },
  {
    key: "bing_wallpaper_on",
    label: "Bing 每日壁纸",
    description: "启用后，当壁纸库为空时自动获取 Bing 每日壁纸",
    type: "switch",
    icon: ImageIcon,
    group: "wallpaper",
  },
]

export default function SettingsPage() {
  const { toast } = useToast()
  const [configs, setConfigs] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  // 加载配置
  useEffect(() => {
    const fetchConfigs = async () => {
      try {
        const data = await configApi.list()
        const configMap: Record<string, string> = {}
        data.list.forEach((config: SystemConfig) => {
          configMap[config.key] = config.value
        })
        setConfigs(configMap)
      } catch (error) {
        console.error("获取配置失败:", error)
      } finally {
        setLoading(false)
      }
    }
    fetchConfigs()
  }, [])

  // 更新配置值
  const updateConfig = (key: string, value: string) => {
    setConfigs({ ...configs, [key]: value })
  }

  // 保存配置
  const handleSave = async () => {
    setSaving(true)
    try {
      const configsToSave = CONFIG_ITEMS.map((item) => ({
        key: item.key,
        value: configs[item.key] || "",
      }))
      await configApi.batchSet(configsToSave)
      toast({
        title: "保存成功",
        description: "系统配置已成功保存",
        variant: "success",
      })
    } catch (error) {
      console.error("保存失败:", error)
      toast({
        title: "保存失败",
        description: error instanceof Error ? error.message : "保存配置时发生错误",
        variant: "destructive",
      })
    } finally {
      setSaving(false)
    }
  }

  // 按分组组织配置项
  const groupedConfigs = {
    weather: CONFIG_ITEMS.filter((item) => item.group === "weather"),
    search: CONFIG_ITEMS.filter((item) => item.group === "search"),
    wallpaper: CONFIG_ITEMS.filter((item) => item.group === "wallpaper"),
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-gray-500">加载中...</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">系统配置</h1>
          <p className="text-gray-500 mt-2">配置后端服务的各项功能</p>
        </div>
        <Button
          onClick={handleSave}
          disabled={saving}
          className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white"
        >
          <Save className="h-4 w-4 mr-2" />
          {saving ? "保存中..." : "保存配置"}
        </Button>
      </div>

      {/* 天气配置 */}
      <Card className="bg-white border-gray-200">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center">
              <Cloud className="w-5 h-5 text-white" />
            </div>
            <div>
              <CardTitle className="text-gray-900">天气服务</CardTitle>
              <CardDescription>配置天气数据来源</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {groupedConfigs.weather.map((item) => (
            <div key={item.key} className="space-y-2">
              <Label className="text-gray-700">{item.label}</Label>
              {item.type === "select" ? (
                <select
                  value={configs[item.key] || "qweather"}
                  onChange={(e) => updateConfig(item.key, e.target.value)}
                  className="w-full h-10 rounded-md border border-gray-200 px-3 bg-white text-gray-900"
                >
                  {item.options?.map((opt) => (
                    <option key={opt.value} value={opt.value}>
                      {opt.label}
                    </option>
                  ))}
                </select>
              ) : (
                <Input
                  value={configs[item.key] || ""}
                  onChange={(e) => updateConfig(item.key, e.target.value)}
                  placeholder={item.description}
                  className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
                />
              )}
              <p className="text-sm text-gray-500">{item.description}</p>
            </div>
          ))}
        </CardContent>
      </Card>

      {/* 搜索配置 */}
      <Card className="bg-white border-gray-200">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-green-500 to-emerald-500 flex items-center justify-center">
              <Search className="w-5 h-5 text-white" />
            </div>
            <div>
              <CardTitle className="text-gray-900">搜索服务</CardTitle>
              <CardDescription>配置搜索相关功能</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {groupedConfigs.search.map((item) => (
            <div key={item.key} className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-gray-700">{item.label}</Label>
                <p className="text-sm text-gray-500">{item.description}</p>
              </div>
              <Switch
                checked={configs[item.key] === "true" || configs[item.key] === "1"}
                onCheckedChange={(checked) => updateConfig(item.key, checked ? "true" : "false")}
              />
            </div>
          ))}
        </CardContent>
      </Card>

      {/* 壁纸配置 */}
      <Card className="bg-white border-gray-200">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center">
              <ImageIcon className="w-5 h-5 text-white" />
            </div>
            <div>
              <CardTitle className="text-gray-900">壁纸服务</CardTitle>
              <CardDescription>配置壁纸相关功能</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {groupedConfigs.wallpaper.map((item) => (
            <div key={item.key} className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label className="text-gray-700">{item.label}</Label>
                <p className="text-sm text-gray-500">{item.description}</p>
              </div>
              <Switch
                checked={configs[item.key] === "true" || configs[item.key] === "1"}
                onCheckedChange={(checked) => updateConfig(item.key, checked ? "true" : "false")}
              />
            </div>
          ))}
        </CardContent>
      </Card>

      {/* 配置说明 */}
      <Card className="bg-orange-50 border-orange-200">
        <CardContent className="pt-6">
          <h3 className="font-medium text-orange-800 mb-2">配置说明</h3>
          <ul className="text-sm text-orange-700 space-y-1">
            <li>- 天气 API Key：需要到对应平台申请，和风天气免费版每天 1000 次调用</li>
            <li>- 搜索建议代理：由于 Google 搜索建议 API 在国内无法直接访问，需要通过后端代理</li>
            <li>- Bing 每日壁纸：当壁纸库为空时，自动从 Bing 获取每日壁纸作为备用</li>
          </ul>
        </CardContent>
      </Card>
    </div>
  )
}
