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
import { categoryApi, type Category, type CreateCategoryRequest } from "@/lib/api"
import { Plus, Pencil, Trash2, GripVertical } from "lucide-react"
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

// 可拖拽的行组件
function SortableRow({
  category,
  onEdit,
  onDelete,
}: {
  category: Category
  onEdit: (category: Category) => void
  onDelete: (id: number) => void
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
  }

  return (
    <tr
      ref={setNodeRef}
      style={style}
      className="border-b border-gray-100 hover:bg-orange-50/50"
    >
      <td className="p-4 w-12">
        <button
          {...attributes}
          {...listeners}
          className="cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600"
        >
          <GripVertical className="h-5 w-5" />
        </button>
      </td>
      <td className="p-4 text-gray-500 w-16">{category.id}</td>
      <td className="p-4 font-medium text-gray-900">{category.name}</td>
      <td className="p-4 text-gray-500">{category.name_en || "-"}</td>
      <td className="p-4">
        <span
          className={`px-2 py-1 rounded-full text-xs ${
            category.is_active
              ? "bg-green-100 text-green-700"
              : "bg-gray-100 text-gray-700"
          }`}
        >
          {category.is_active ? "启用" : "禁用"}
        </span>
      </td>
      <td className="p-4">
        <div className="flex gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onEdit(category)}
            className="text-gray-600 hover:text-orange-500 hover:bg-orange-50"
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDelete(category.id)}
            className="text-gray-600 hover:text-red-500 hover:bg-red-50"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </td>
    </tr>
  )
}

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [formData, setFormData] = useState<CreateCategoryRequest>({
    name: "",
    name_en: "",
    sort_order: 0,
    is_active: true,
  })

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
      // 按 sort_order 排序
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
        // 失败时重新获取列表
        fetchCategories()
      }
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
    if (!confirm("确定要删除这个分类吗？")) return
    try {
      await categoryApi.delete(id)
      fetchCategories()
    } catch (error) {
      console.error("删除失败:", error)
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
    } catch (error) {
      console.error("保存失败:", error)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">分类管理</h1>
          <p className="text-gray-500 mt-2">拖动调整顺序，管理图标分类</p>
        </div>
        <Button onClick={handleCreate} className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white">
          <Plus className="h-4 w-4 mr-2" />
          添加分类
        </Button>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
        >
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="p-4 text-left text-gray-700 font-medium w-12"></th>
                <th className="p-4 text-left text-gray-700 font-medium w-16">ID</th>
                <th className="p-4 text-left text-gray-700 font-medium">名称</th>
                <th className="p-4 text-left text-gray-700 font-medium">英文名</th>
                <th className="p-4 text-left text-gray-700 font-medium w-20">状态</th>
                <th className="p-4 text-left text-gray-700 font-medium w-24">操作</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={6} className="text-center py-8 text-gray-500">
                    加载中...
                  </td>
                </tr>
              ) : categories.length === 0 ? (
                <tr>
                  <td colSpan={6} className="text-center py-8 text-gray-500">
                    暂无数据
                  </td>
                </tr>
              ) : (
                <SortableContext
                  items={categories.map((c) => c.id)}
                  strategy={verticalListSortingStrategy}
                >
                  {categories.map((category) => (
                    <SortableRow
                      key={category.id}
                      category={category}
                      onEdit={handleEdit}
                      onDelete={handleDelete}
                    />
                  ))}
                </SortableContext>
              )}
            </tbody>
          </table>
        </DndContext>
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
    </div>
  )
}
