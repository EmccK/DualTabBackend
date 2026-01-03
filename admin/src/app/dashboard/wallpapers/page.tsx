"use client"

import { useEffect, useState, useCallback } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog"
import { ConfirmDialog } from "@/components/ui/confirm-dialog"
import { wallpaperApi, type Wallpaper, type CreateWallpaperRequest } from "@/lib/api"
import { Plus, Pencil, Trash2, Upload, ImageIcon } from "lucide-react"
import { useToast } from "@/hooks/use-toast"

export default function WallpapersPage() {
  const { toast } = useToast()
  const [wallpapers, setWallpapers] = useState<Wallpaper[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [deletingId, setDeletingId] = useState<number | null>(null)
  const [editingWallpaper, setEditingWallpaper] = useState<Wallpaper | null>(null)
  const [uploading, setUploading] = useState(false)
  const [formData, setFormData] = useState<CreateWallpaperRequest>({
    title: "",
    url: "",
    thumb_url: "",
    source: "",
    sort_order: 0,
    is_active: true,
  })

  const fetchWallpapers = useCallback(async () => {
    setLoading(true)
    try {
      const data = await wallpaperApi.list({ page, size: 12 })
      setWallpapers(data.list)
      setTotal(data.total)
    } catch (error) {
      console.error("获取壁纸列表失败:", error)
    } finally {
      setLoading(false)
    }
  }, [page])

  useEffect(() => {
    fetchWallpapers()
  }, [fetchWallpapers])

  const handleCreate = () => {
    setEditingWallpaper(null)
    setFormData({
      title: "",
      url: "",
      thumb_url: "",
      source: "",
      sort_order: 0,
      is_active: true,
    })
    setDialogOpen(true)
  }

  const handleEdit = (wallpaper: Wallpaper) => {
    setEditingWallpaper(wallpaper)
    setFormData({
      title: wallpaper.title,
      url: wallpaper.url,
      thumb_url: wallpaper.thumb_url,
      source: wallpaper.source,
      sort_order: wallpaper.sort_order,
      is_active: wallpaper.is_active,
    })
    setDialogOpen(true)
  }

  const handleDelete = async (id: number) => {
    setDeletingId(id)
    setConfirmOpen(true)
  }

  const confirmDelete = async () => {
    if (!deletingId) return
    try {
      await wallpaperApi.delete(deletingId)
      fetchWallpapers()
      toast({
        title: "删除成功",
        description: "壁纸已成功删除",
        variant: "success",
      })
    } catch (error) {
      console.error("删除失败:", error)
      toast({
        title: "删除失败",
        description: error instanceof Error ? error.message : "删除壁纸时发生错误",
        variant: "destructive",
      })
    } finally {
      setDeletingId(null)
    }
  }

  const handleSubmit = async () => {
    try {
      if (editingWallpaper) {
        await wallpaperApi.update(editingWallpaper.id, formData)
      } else {
        await wallpaperApi.create(formData)
      }
      setDialogOpen(false)
      fetchWallpapers()
      toast({
        title: "保存成功",
        description: editingWallpaper ? "壁纸已更新" : "壁纸已创建",
        variant: "success",
      })
    } catch (error) {
      console.error("保存失败:", error)
      toast({
        title: "保存失败",
        description: error instanceof Error ? error.message : "保存壁纸时发生错误",
        variant: "destructive",
      })
    }
  }

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setUploading(true)
    try {
      const data = await wallpaperApi.upload(file)
      setFormData({ ...formData, url: data.url, thumb_url: data.url })
    } catch (error) {
      console.error("上传失败:", error)
    } finally {
      setUploading(false)
    }
  }

  const totalPages = Math.ceil(total / 12)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">壁纸管理</h1>
          <p className="text-gray-500 mt-2">管理壁纸库，供扩展随机获取</p>
        </div>
        <Button onClick={handleCreate} className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white">
          <Plus className="h-4 w-4 mr-2" />
          添加壁纸
        </Button>
      </div>

      {/* 壁纸网格 */}
      <div className="bg-white rounded-lg border border-gray-200 shadow-sm p-6">
        {loading ? (
          <div className="text-center py-12 text-gray-500">加载中...</div>
        ) : wallpapers.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            <ImageIcon className="w-12 h-12 mx-auto mb-4 text-gray-300" />
            <p>暂无壁纸，点击上方按钮添加</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {wallpapers.map((wallpaper) => (
              <div
                key={wallpaper.id}
                className="group relative aspect-video rounded-lg overflow-hidden border border-gray-200 bg-gray-100"
              >
                <img
                  src={wallpaper.thumb_url || wallpaper.url}
                  alt={wallpaper.title}
                  className="w-full h-full object-cover"
                />
                {/* 遮罩层 */}
                <div className="absolute inset-0 bg-black/0 group-hover:bg-black/50 transition-all duration-200">
                  <div className="absolute inset-0 flex items-center justify-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => handleEdit(wallpaper)}
                      className="bg-white/90 hover:bg-white text-gray-700"
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      size="sm"
                      variant="secondary"
                      onClick={() => handleDelete(wallpaper.id)}
                      className="bg-white/90 hover:bg-red-50 text-gray-700 hover:text-red-500"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
                {/* 标题和状态 */}
                <div className="absolute bottom-0 left-0 right-0 p-2 bg-gradient-to-t from-black/70 to-transparent">
                  <div className="flex items-center justify-between">
                    <span className="text-white text-sm truncate">{wallpaper.title}</span>
                    <span
                      className={`px-1.5 py-0.5 rounded text-xs ${
                        wallpaper.is_active
                          ? "bg-green-500/80 text-white"
                          : "bg-gray-500/80 text-white"
                      }`}
                    >
                      {wallpaper.is_active ? "启用" : "禁用"}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* 分页 */}
        {totalPages > 1 && (
          <div className="flex items-center justify-center gap-2 mt-6 pt-6 border-t border-gray-200">
            <Button
              variant="outline"
              size="sm"
              disabled={page === 1}
              onClick={() => setPage(page - 1)}
              className="border-gray-200 text-gray-700 hover:bg-gray-50"
            >
              上一页
            </Button>
            <span className="text-sm text-gray-500">
              {page} / {totalPages}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={page === totalPages}
              onClick={() => setPage(page + 1)}
              className="border-gray-200 text-gray-700 hover:bg-gray-50"
            >
              下一页
            </Button>
          </div>
        )}
      </div>

      {/* 编辑对话框 */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-md bg-white">
          <DialogHeader>
            <DialogTitle className="text-gray-900">
              {editingWallpaper ? "编辑壁纸" : "添加壁纸"}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="text-gray-700">标题 *</Label>
              <Input
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                placeholder="请输入标题"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-gray-700">壁纸图片 *</Label>
              <div className="space-y-2">
                {formData.url && (
                  <div className="aspect-video rounded-lg overflow-hidden border border-gray-200 bg-gray-100">
                    <img
                      src={formData.url}
                      alt="预览"
                      className="w-full h-full object-cover"
                    />
                  </div>
                )}
                <div className="flex gap-2">
                  <Input
                    value={formData.url}
                    onChange={(e) => setFormData({ ...formData, url: e.target.value })}
                    placeholder="图片 URL"
                    className="flex-1 bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
                  />
                  <Button
                    variant="outline"
                    asChild
                    disabled={uploading}
                    className="border-gray-200 text-gray-700 hover:bg-gray-50"
                  >
                    <label className="cursor-pointer">
                      <Upload className="h-4 w-4" />
                      <input
                        type="file"
                        accept="image/*"
                        className="hidden"
                        onChange={handleUpload}
                      />
                    </label>
                  </Button>
                </div>
                {uploading && <p className="text-sm text-gray-500">上传中...</p>}
              </div>
            </div>

            <div className="space-y-2">
              <Label className="text-gray-700">缩略图 URL</Label>
              <Input
                value={formData.thumb_url}
                onChange={(e) => setFormData({ ...formData, thumb_url: e.target.value })}
                placeholder="可选，留空则使用原图"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label className="text-gray-700">来源</Label>
                <Input
                  value={formData.source}
                  onChange={(e) => setFormData({ ...formData, source: e.target.value })}
                  placeholder="如 Unsplash"
                  className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
                />
              </div>
              <div className="space-y-2">
                <Label className="text-gray-700">排序</Label>
                <Input
                  type="number"
                  value={formData.sort_order}
                  onChange={(e) => setFormData({ ...formData, sort_order: Number(e.target.value) })}
                  className="bg-white text-gray-900 border-gray-200"
                />
              </div>
            </div>

            <div className="flex items-center gap-2">
              <Switch
                checked={formData.is_active}
                onCheckedChange={(checked) => setFormData({ ...formData, is_active: checked })}
              />
              <Label className="text-gray-700">启用</Label>
            </div>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDialogOpen(false)}
              className="border-gray-200 text-gray-700 hover:bg-gray-50"
            >
              取消
            </Button>
            <Button
              onClick={handleSubmit}
              className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white"
            >
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 删除确认对话框 */}
      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="删除壁纸"
        description="确定要删除这张壁纸吗？此操作无法撤销。"
        onConfirm={confirmDelete}
        confirmText="删除"
        cancelText="取消"
        variant="destructive"
      />
    </div>
  )
}
