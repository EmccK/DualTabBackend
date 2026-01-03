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
import { categoryApi, iconApi, type Category, type CreateCategoryRequest, type Icon, type CreateIconRequest } from "@/lib/api"
import { Plus, Pencil, Trash2, ChevronDown, ChevronRight, ExternalLink, GripVertical, Upload } from "lucide-react"
import { useToast } from "@/hooks/use-toast"
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core"
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable"
import { CSS } from "@dnd-kit/utilities"

// 可拖拽的分类行组件
function SortableCategoryRow({
  category,
  isExpanded,
  icons,
  isLoadingIcons,
  onToggleExpand,
  onEdit,
  onDelete,
  onAddIcon,
  onEditIcon,
  onDeleteIcon,
}: {
  category: Category
  isExpanded: boolean
  icons: Icon[]
  isLoadingIcons: boolean
  onToggleExpand: (id: number) => void
  onEdit: (category: Category) => void
  onDelete: (id: number) => void
  onAddIcon: (categoryId: number) => void
  onEditIcon: (icon: Icon) => void
  onDeleteIcon: (iconId: number, categoryId: number) => void
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: category.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    zIndex: isDragging ? 1 : 0,
  }

  return (
    <div ref={setNodeRef} style={style}>
      {/* 分类行 */}
      <div className="flex items-center p-4 hover:bg-gray-50 bg-white">
        {/* 拖拽手柄 */}
        <button
          {...attributes}
          {...listeners}
          className="p-1 mr-2 cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600"
        >
          <GripVertical className="h-5 w-5" />
        </button>

        {/* 展开按钮 */}
        <button
          onClick={() => onToggleExpand(category.id)}
          className="p-1 mr-2 text-gray-400 hover:text-gray-600 rounded"
        >
          {isExpanded ? (
            <ChevronDown className="h-5 w-5" />
          ) : (
            <ChevronRight className="h-5 w-5" />
          )}
        </button>

        {/* 分类信息 */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3">
            <span className="font-medium text-gray-900">{category.name}</span>
            {category.name_en && (
              <span className="text-sm text-gray-400">({category.name_en})</span>
            )}
            <span
              className={`px-2 py-0.5 rounded-full text-xs ${
                category.is_active
                  ? "bg-green-100 text-green-700"
                  : "bg-gray-100 text-gray-500"
              }`}
            >
              {category.is_active ? "启用" : "禁用"}
            </span>
            {icons.length > 0 && (
              <span className="text-sm text-gray-400">
                {icons.length} 个书签
              </span>
            )}
          </div>
        </div>

        {/* 操作按钮 */}
        <div className="flex gap-1">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onEdit(category)}
            className="text-gray-500 hover:text-orange-500 hover:bg-orange-50"
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDelete(category.id)}
            className="text-gray-500 hover:text-red-500 hover:bg-red-50"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* 展开的图标列表 */}
      {isExpanded && (
        <div className="bg-gray-50 border-t border-gray-100">
          {isLoadingIcons ? (
            <div className="p-4 text-center text-gray-500 text-sm">
              加载中...
            </div>
          ) : (
            <div className="p-4">
              {/* 添加书签按钮 */}
              <div className="mb-3 flex justify-end">
                <Button
                  size="sm"
                  onClick={() => onAddIcon(category.id)}
                  className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white"
                >
                  <Plus className="h-4 w-4 mr-1" />
                  添加书签
                </Button>
              </div>

              {icons.length === 0 ? (
                <div className="text-center text-gray-400 text-sm py-4">
                  该分类下暂无书签，点击上方按钮添加
                </div>
              ) : (
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
                  {icons.map((icon) => (
                    <div
                      key={icon.id}
                      className="group relative flex items-center gap-3 p-3 bg-white rounded-lg border border-gray-200 hover:border-orange-300 hover:shadow-sm transition-all"
                    >
                      <div
                        className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0"
                        style={{ backgroundColor: icon.bg_color || '#f3f4f6' }}
                      >
                        {icon.img_url ? (
                          <img
                            src={icon.img_url}
                            alt={icon.title}
                            className="w-6 h-6 object-contain"
                          />
                        ) : (
                          <span className="text-sm font-medium text-gray-600">
                            {icon.title[0]}
                          </span>
                        )}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="font-medium text-gray-900 text-sm truncate">
                          {icon.title}
                        </div>
                        <div className="text-xs text-gray-400 truncate">
                          {(() => {
                            try {
                              return new URL(icon.url).hostname
                            } catch {
                              return icon.url
                            }
                          })()}
                        </div>
                      </div>
                      {/* 操作按钮 - 悬停显示 */}
                      <div className="absolute right-2 top-2 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button
                          onClick={() => onEditIcon(icon)}
                          className="p-1 rounded bg-white shadow-sm border border-gray-200 text-gray-500 hover:text-orange-500 hover:border-orange-300"
                        >
                          <Pencil className="h-3 w-3" />
                        </button>
                        <button
                          onClick={() => onDeleteIcon(icon.id, category.id)}
                          className="p-1 rounded bg-white shadow-sm border border-gray-200 text-gray-500 hover:text-red-500 hover:border-red-300"
                        >
                          <Trash2 className="h-3 w-3" />
                        </button>
                        <a
                          href={icon.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="p-1 rounded bg-white shadow-sm border border-gray-200 text-gray-500 hover:text-blue-500 hover:border-blue-300"
                        >
                          <ExternalLink className="h-3 w-3" />
                        </a>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default function CategoriesPage() {
  const { toast } = useToast()
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [deletingId, setDeletingId] = useState<number | null>(null)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [formData, setFormData] = useState<CreateCategoryRequest>({
    name: "",
    name_en: "",
    sort_order: 0,
    is_active: true,
  })

  // 展开状态和图标数据
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set())
  const [categoryIcons, setCategoryIcons] = useState<Record<number, Icon[]>>({})
  const [loadingIcons, setLoadingIcons] = useState<Set<number>>(new Set())

  // 书签编辑相关状态
  const [iconDialogOpen, setIconDialogOpen] = useState(false)
  const [iconConfirmOpen, setIconConfirmOpen] = useState(false)
  const [deletingIconId, setDeletingIconId] = useState<number | null>(null)
  const [deletingIconCategoryId, setDeletingIconCategoryId] = useState<number | null>(null)
  const [editingIcon, setEditingIcon] = useState<Icon | null>(null)
  const [iconFormData, setIconFormData] = useState<CreateIconRequest>({
    title: "",
    description: "",
    url: "",
    img_url: "",
    bg_color: "#ffffff",
    category_id: 0,
    sort_order: 0,
    is_active: true,
  })
  const [uploading, setUploading] = useState(false)

  // 拖拽传感器配置
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  )

  const fetchCategories = useCallback(async () => {
    setLoading(true)
    try {
      const data = await categoryApi.list()
      const sorted = [...data.list].sort((a, b) => a.sort_order - b.sort_order)
      setCategories(sorted)
    } catch (error) {
      console.error("获取分类列表失败:", error)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchCategories()
  }, [fetchCategories])

  // 处理拖拽结束
  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event

    if (over && active.id !== over.id) {
      const oldIndex = categories.findIndex((c) => c.id === active.id)
      const newIndex = categories.findIndex((c) => c.id === over.id)

      const newCategories = arrayMove(categories, oldIndex, newIndex)
      setCategories(newCategories)

      // 更新排序到服务器
      try {
        await Promise.all(
          newCategories.map((category, index) =>
            categoryApi.update(category.id, { sort_order: index })
          )
        )
      } catch (error) {
        console.error("更新排序失败:", error)
        fetchCategories()
      }
    }
  }

  // 切换展开/收起
  const toggleExpand = async (categoryId: number) => {
    const newExpanded = new Set(expandedIds)

    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId)
    } else {
      newExpanded.add(categoryId)
      // 如果还没加载过该分类的图标，则加载
      if (!categoryIcons[categoryId]) {
        await fetchCategoryIcons(categoryId)
      }
    }

    setExpandedIds(newExpanded)
  }

  // 获取分类下的图标
  const fetchCategoryIcons = async (categoryId: number) => {
    setLoadingIcons(prev => new Set(prev).add(categoryId))
    try {
      const data = await iconApi.list({ category_id: categoryId, size: 100 })
      setCategoryIcons(prev => ({ ...prev, [categoryId]: data.list }))
    } catch (error) {
      console.error("获取图标失败:", error)
    } finally {
      setLoadingIcons(prev => {
        const next = new Set(prev)
        next.delete(categoryId)
        return next
      })
    }
  }

  const handleCreate = () => {
    setEditingCategory(null)
    setFormData({
      name: "",
      name_en: "",
      sort_order: categories.length,
      is_active: true,
    })
    setDialogOpen(true)
  }

  const handleEdit = (category: Category) => {
    setEditingCategory(category)
    setFormData({
      name: category.name,
      name_en: category.name_en,
      sort_order: category.sort_order,
      is_active: category.is_active,
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
      await categoryApi.delete(deletingId)
      fetchCategories()
      toast({
        title: "删除成功",
        description: "分类已成功删除",
        variant: "success",
      })
    } catch (error) {
      console.error("删除失败:", error)
      toast({
        title: "删除失败",
        description: error instanceof Error ? error.message : "删除分类时发生错误",
        variant: "destructive",
      })
    } finally {
      setDeletingId(null)
    }
  }

  const handleSubmit = async () => {
    try {
      if (editingCategory) {
        await categoryApi.update(editingCategory.id, formData)
      } else {
        await categoryApi.create(formData)
      }
      setDialogOpen(false)
      fetchCategories()
      toast({
        title: "保存成功",
        description: editingCategory ? "分类已更新" : "分类已创建",
        variant: "success",
      })
    } catch (error) {
      console.error("保存失败:", error)
      toast({
        title: "保存失败",
        description: error instanceof Error ? error.message : "保存分类时发生错误",
        variant: "destructive",
      })
    }
  }

  // ========== 书签相关操作 ==========
  const handleAddIcon = (categoryId: number) => {
    setEditingIcon(null)
    setIconFormData({
      title: "",
      description: "",
      url: "",
      img_url: "",
      bg_color: "#ffffff",
      category_id: categoryId,
      sort_order: 0,
      is_active: true,
    })
    setIconDialogOpen(true)
  }

  const handleEditIcon = (icon: Icon) => {
    setEditingIcon(icon)
    setIconFormData({
      title: icon.title,
      description: icon.description,
      url: icon.url,
      img_url: icon.img_url,
      bg_color: icon.bg_color,
      category_id: icon.category_id,
      sort_order: icon.sort_order,
      is_active: icon.is_active,
    })
    setIconDialogOpen(true)
  }

  const handleDeleteIcon = async (iconId: number, categoryId: number) => {
    setDeletingIconId(iconId)
    setDeletingIconCategoryId(categoryId)
    setIconConfirmOpen(true)
  }

  const confirmDeleteIcon = async () => {
    if (!deletingIconId || !deletingIconCategoryId) return
    try {
      await iconApi.delete(deletingIconId)
      await fetchCategoryIcons(deletingIconCategoryId)
      toast({
        title: "删除成功",
        description: "书签已成功删除",
        variant: "success",
      })
    } catch (error) {
      console.error("删除失败:", error)
      toast({
        title: "删除失败",
        description: error instanceof Error ? error.message : "删除书签时发生错误",
        variant: "destructive",
      })
    } finally {
      setDeletingIconId(null)
      setDeletingIconCategoryId(null)
    }
  }

  const handleIconSubmit = async () => {
    try {
      if (editingIcon) {
        await iconApi.update(editingIcon.id, iconFormData)
      } else {
        await iconApi.create(iconFormData)
      }
      setIconDialogOpen(false)
      if (iconFormData.category_id) {
        await fetchCategoryIcons(iconFormData.category_id)
      }
      toast({
        title: "保存成功",
        description: editingIcon ? "书签已更新" : "书签已创建",
        variant: "success",
      })
    } catch (error) {
      console.error("保存失败:", error)
      toast({
        title: "保存失败",
        description: error instanceof Error ? error.message : "保存书签时发生错误",
        variant: "destructive",
      })
    }
  }

  const handleIconUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setUploading(true)
    try {
      const data = await iconApi.upload(file)
      setIconFormData({ ...iconFormData, img_url: data.url })
    } catch (error) {
      console.error("上传失败:", error)
    } finally {
      setUploading(false)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">分类管理</h1>
          <p className="text-gray-500 mt-2">拖动排序，点击箭头展开查看书签</p>
        </div>
        <Button onClick={handleCreate} className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white">
          <Plus className="h-4 w-4 mr-2" />
          添加分类
        </Button>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
        {loading ? (
          <div className="text-center py-12 text-gray-500">加载中...</div>
        ) : categories.length === 0 ? (
          <div className="text-center py-12 text-gray-500">暂无数据</div>
        ) : (
          <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            onDragEnd={handleDragEnd}
          >
            <SortableContext
              items={categories.map((c) => c.id)}
              strategy={verticalListSortingStrategy}
            >
              <div className="divide-y divide-gray-100">
                {categories.map((category) => (
                  <SortableCategoryRow
                    key={category.id}
                    category={category}
                    isExpanded={expandedIds.has(category.id)}
                    icons={categoryIcons[category.id] || []}
                    isLoadingIcons={loadingIcons.has(category.id)}
                    onToggleExpand={toggleExpand}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                    onAddIcon={handleAddIcon}
                    onEditIcon={handleEditIcon}
                    onDeleteIcon={handleDeleteIcon}
                  />
                ))}
              </div>
            </SortableContext>
          </DndContext>
        )}
      </div>

      {/* 编辑对话框 */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-md bg-white">
          <DialogHeader>
            <DialogTitle className="text-gray-900">
              {editingCategory ? "编辑分类" : "添加分类"}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="text-gray-700">名称 *</Label>
              <Input
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                placeholder="请输入分类名称"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">英文名</Label>
              <Input
                value={formData.name_en}
                onChange={(e) =>
                  setFormData({ ...formData, name_en: e.target.value })
                }
                placeholder="请输入英文名称"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
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

      {/* 书签编辑对话框 */}
      <Dialog open={iconDialogOpen} onOpenChange={setIconDialogOpen}>
        <DialogContent className="max-w-md bg-white">
          <DialogHeader>
            <DialogTitle className="text-gray-900">
              {editingIcon ? "编辑书签" : "添加书签"}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="text-gray-700">标题 *</Label>
              <Input
                value={iconFormData.title}
                onChange={(e) =>
                  setIconFormData({ ...iconFormData, title: e.target.value })
                }
                placeholder="请输入标题"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">描述</Label>
              <Input
                value={iconFormData.description}
                onChange={(e) =>
                  setIconFormData({ ...iconFormData, description: e.target.value })
                }
                placeholder="请输入描述"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">URL *</Label>
              <Input
                value={iconFormData.url}
                onChange={(e) =>
                  setIconFormData({ ...iconFormData, url: e.target.value })
                }
                placeholder="https://example.com"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">图标</Label>
              <div className="flex gap-2">
                <Input
                  value={iconFormData.img_url}
                  onChange={(e) =>
                    setIconFormData({ ...iconFormData, img_url: e.target.value })
                  }
                  placeholder="图标 URL"
                  className="flex-1 bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
                />
                <Button variant="outline" asChild disabled={uploading} className="border-gray-200 text-gray-700 hover:bg-gray-50">
                  <label className="cursor-pointer">
                    <Upload className="h-4 w-4" />
                    <input
                      type="file"
                      accept="image/*"
                      className="hidden"
                      onChange={handleIconUpload}
                    />
                  </label>
                </Button>
              </div>
              {uploading && <p className="text-sm text-gray-500">上传中...</p>}
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label className="text-gray-700">背景色</Label>
                <div className="flex gap-2">
                  <Input
                    type="color"
                    value={iconFormData.bg_color}
                    onChange={(e) =>
                      setIconFormData({ ...iconFormData, bg_color: e.target.value })
                    }
                    className="w-12 h-9 p-1 border-gray-200"
                  />
                  <Input
                    value={iconFormData.bg_color}
                    onChange={(e) =>
                      setIconFormData({ ...iconFormData, bg_color: e.target.value })
                    }
                    className="flex-1 bg-white text-gray-900 border-gray-200"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label className="text-gray-700">分类</Label>
                <select
                  value={iconFormData.category_id}
                  onChange={(e) =>
                    setIconFormData({
                      ...iconFormData,
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
                checked={iconFormData.is_active}
                onCheckedChange={(checked) =>
                  setIconFormData({ ...iconFormData, is_active: checked })
                }
              />
              <Label className="text-gray-700">启用</Label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIconDialogOpen(false)} className="border-gray-200 text-gray-700 hover:bg-gray-50">
              取消
            </Button>
            <Button onClick={handleIconSubmit} className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white">保存</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 删除分类确认对话框 */}
      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="删除分类"
        description="确定要删除这个分类吗？删除后该分类下的书签将变为未分类状态。"
        onConfirm={confirmDelete}
        confirmText="删除"
        cancelText="取消"
        variant="destructive"
      />

      {/* 删除书签确认对话框 */}
      <ConfirmDialog
        open={iconConfirmOpen}
        onOpenChange={setIconConfirmOpen}
        title="删除书签"
        description="确定要删除这个书签吗？此操作无法撤销。"
        onConfirm={confirmDeleteIcon}
        confirmText="删除"
        cancelText="取消"
        variant="destructive"
      />
    </div>
  )
}
