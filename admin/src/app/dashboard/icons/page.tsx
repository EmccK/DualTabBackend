"use client"

import { useEffect, useState, useCallback } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Switch } from "@/components/ui/switch"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog"
import { iconApi, categoryApi, type Icon, type Category, type CreateIconRequest } from "@/lib/api"
import { Plus, Pencil, Trash2, Upload } from "lucide-react"

export default function IconsPage() {
  const [icons, setIcons] = useState<Icon[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingIcon, setEditingIcon] = useState<Icon | null>(null)
  const [formData, setFormData] = useState<CreateIconRequest>({
    title: "",
    description: "",
    url: "",
    img_url: "",
    bg_color: "#ffffff",
    category_id: 0,
    sort_order: 0,
    is_active: true,
  })

  const fetchIcons = useCallback(async () => {
    setLoading(true)
    try {
      const data = await iconApi.list({ page, size: 20 })
      setIcons(data.list)
      setTotal(data.total)
    } catch (error) {
      console.error("获取图标列表失败:", error)
    } finally {
      setLoading(false)
    }
  }, [page])

  const fetchCategories = async () => {
    try {
      const data = await categoryApi.list()
      setCategories(data.list)
    } catch (error) {
      console.error("获取分类列表失败:", error)
    }
  }

  useEffect(() => {
    fetchIcons()
    fetchCategories()
  }, [fetchIcons])

  const handleCreate = () => {
    setEditingIcon(null)
    setFormData({
      title: "",
      description: "",
      url: "",
      img_url: "",
      bg_color: "#ffffff",
      category_id: categories[0]?.id || 0,
      sort_order: 0,
      is_active: true,
    })
    setDialogOpen(true)
  }

  const handleEdit = (icon: Icon) => {
    setEditingIcon(icon)
    setFormData({
      title: icon.title,
      description: icon.description,
      url: icon.url,
      img_url: icon.img_url,
      bg_color: icon.bg_color,
      category_id: icon.category_id,
      sort_order: icon.sort_order,
      is_active: icon.is_active,
    })
    setDialogOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm("确定要删除这个图标吗？")) return
    try {
      await iconApi.delete(id)
      fetchIcons()
    } catch (error) {
      console.error("删除失败:", error)
    }
  }

  const handleSubmit = async () => {
    try {
      if (editingIcon) {
        await iconApi.update(editingIcon.id, formData)
      } else {
        await iconApi.create(formData)
      }
      setDialogOpen(false)
      fetchIcons()
    } catch (error) {
      console.error("保存失败:", error)
    }
  }

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    try {
      const data = await iconApi.upload(file)
      setFormData({ ...formData, img_url: data.url })
    } catch (error) {
      console.error("上传失败:", error)
    }
  }

  const totalPages = Math.ceil(total / 20)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">图标管理</h1>
          <p className="text-gray-500 mt-2">管理推荐书签图标</p>
        </div>
        <Button onClick={handleCreate} className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white">
          <Plus className="h-4 w-4 mr-2" />
          添加图标
        </Button>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow className="bg-gray-50">
              <TableHead className="w-16 text-gray-700">图标</TableHead>
              <TableHead className="text-gray-700">标题</TableHead>
              <TableHead className="text-gray-700">URL</TableHead>
              <TableHead className="text-gray-700">分类</TableHead>
              <TableHead className="w-20 text-gray-700">状态</TableHead>
              <TableHead className="w-24 text-gray-700">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-gray-500">
                  加载中...
                </TableCell>
              </TableRow>
            ) : icons.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-gray-500">
                  暂无数据
                </TableCell>
              </TableRow>
            ) : (
              icons.map((icon) => (
                <TableRow key={icon.id}>
                  <TableCell>
                    <div
                      className="w-10 h-10 rounded-lg flex items-center justify-center"
                      style={{ backgroundColor: icon.bg_color }}
                    >
                      {icon.img_url ? (
                        <img
                          src={icon.img_url}
                          alt={icon.title}
                          className="w-6 h-6 object-contain"
                        />
                      ) : (
                        <span className="text-xs">{icon.title[0]}</span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell className="font-medium text-gray-900">{icon.title}</TableCell>
                  <TableCell className="text-gray-500 max-w-xs truncate">
                    {icon.url}
                  </TableCell>
                  <TableCell className="text-gray-700">{icon.category?.name || "-"}</TableCell>
                  <TableCell>
                    <span
                      className={`px-2 py-1 rounded-full text-xs ${
                        icon.is_active
                          ? "bg-green-100 text-green-700"
                          : "bg-gray-100 text-gray-700"
                      }`}
                    >
                      {icon.is_active ? "启用" : "禁用"}
                    </span>
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleEdit(icon)}
                        className="text-gray-600 hover:text-orange-500 hover:bg-orange-50"
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDelete(icon.id)}
                        className="text-gray-600 hover:text-red-500 hover:bg-red-50"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {/* 分页 */}
        {totalPages > 1 && (
          <div className="flex items-center justify-center gap-2 p-4 border-t border-gray-200">
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
            <DialogTitle className="text-gray-900">{editingIcon ? "编辑图标" : "添加图标"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="text-gray-700">标题 *</Label>
              <Input
                value={formData.title}
                onChange={(e) =>
                  setFormData({ ...formData, title: e.target.value })
                }
                placeholder="请输入标题"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">描述</Label>
              <Input
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                placeholder="请输入描述"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">URL *</Label>
              <Input
                value={formData.url}
                onChange={(e) =>
                  setFormData({ ...formData, url: e.target.value })
                }
                placeholder="https://example.com"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">图标</Label>
              <div className="flex gap-2">
                <Input
                  value={formData.img_url}
                  onChange={(e) =>
                    setFormData({ ...formData, img_url: e.target.value })
                  }
                  placeholder="图标 URL"
                  className="flex-1 bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
                />
                <Button variant="outline" asChild className="border-gray-200 text-gray-700 hover:bg-gray-50">
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
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label className="text-gray-700">背景色</Label>
                <div className="flex gap-2">
                  <Input
                    type="color"
                    value={formData.bg_color}
                    onChange={(e) =>
                      setFormData({ ...formData, bg_color: e.target.value })
                    }
                    className="w-12 h-9 p-1 border-gray-200"
                  />
                  <Input
                    value={formData.bg_color}
                    onChange={(e) =>
                      setFormData({ ...formData, bg_color: e.target.value })
                    }
                    className="flex-1 bg-white text-gray-900 border-gray-200"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label className="text-gray-700">分类</Label>
                <select
                  value={formData.category_id}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      category_id: Number(e.target.value),
                    })
                  }
                  className="w-full h-9 rounded-md border border-gray-200 px-3 bg-white text-gray-900"
                >
                  <option value={0}>请选择分类</option>
                  {categories.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Switch
                checked={formData.is_active}
                onCheckedChange={(checked) =>
                  setFormData({ ...formData, is_active: checked })
                }
              />
              <Label className="text-gray-700">启用</Label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)} className="border-gray-200 text-gray-700 hover:bg-gray-50">
              取消
            </Button>
            <Button onClick={handleSubmit} className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white">保存</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
