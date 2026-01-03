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
import { searchEngineApi, iconApi, type SearchEngine, type CreateSearchEngineRequest } from "@/lib/api"
import { Plus, Pencil, Trash2, Upload, GripVertical } from "lucide-react"
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

// 可拖拽的行组件
function SortableRow({
  engine,
  onEdit,
  onDelete,
}: {
  engine: SearchEngine
  onEdit: (engine: SearchEngine) => void
  onDelete: (id: number) => void
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: engine.id })

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
      <td className="p-4 w-16">
        {engine.icon_url ? (
          <img
            src={engine.icon_url}
            alt={engine.name}
            className="w-8 h-8 object-contain"
          />
        ) : (
          <div className="w-8 h-8 bg-gray-100 rounded flex items-center justify-center text-xs text-gray-600">
            {engine.name[0]}
          </div>
        )}
      </td>
      <td className="p-4 font-medium text-gray-900">{engine.name}</td>
      <td className="p-4 text-gray-500 max-w-xs truncate">{engine.url}</td>
      <td className="p-4">
        <span
          className={`px-2 py-1 rounded-full text-xs ${
            engine.is_active
              ? "bg-green-100 text-green-700"
              : "bg-gray-100 text-gray-700"
          }`}
        >
          {engine.is_active ? "启用" : "禁用"}
        </span>
      </td>
      <td className="p-4">
        <div className="flex gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onEdit(engine)}
            className="text-gray-600 hover:text-orange-500 hover:bg-orange-50"
          >
            <Pencil className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDelete(engine.id)}
            className="text-gray-600 hover:text-red-500 hover:bg-red-50"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </td>
    </tr>
  )
}

export default function SearchEnginesPage() {
  const { toast } = useToast()
  const [engines, setEngines] = useState<SearchEngine[]>([])
  const [loading, setLoading] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [deletingId, setDeletingId] = useState<number | null>(null)
  const [editingEngine, setEditingEngine] = useState<SearchEngine | null>(null)
  const [formData, setFormData] = useState<CreateSearchEngineRequest>({
    name: "",
    url: "",
    icon_url: "",
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

  const fetchEngines = useCallback(async () => {
    setLoading(true)
    try {
      const data = await searchEngineApi.list()
      // 按 sort_order 排序
      const sorted = [...data.list].sort((a, b) => a.sort_order - b.sort_order)
      setEngines(sorted)
    } catch (error) {
      console.error("获取搜索引擎列表失败:", error)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchEngines()
  }, [fetchEngines])

  // 处理拖拽结束
  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event

    if (over && active.id !== over.id) {
      const oldIndex = engines.findIndex((e) => e.id === active.id)
      const newIndex = engines.findIndex((e) => e.id === over.id)

      const newEngines = arrayMove(engines, oldIndex, newIndex)
      setEngines(newEngines)

      // 更新排序到服务器
      try {
        await Promise.all(
          newEngines.map((engine, index) =>
            searchEngineApi.update(engine.id, { sort_order: index })
          )
        )
      } catch (error) {
        console.error("更新排序失败:", error)
        // 失败时重新获取列表
        fetchEngines()
      }
    }
  }

  const handleCreate = () => {
    setEditingEngine(null)
    setFormData({
      name: "",
      url: "",
      icon_url: "",
      sort_order: engines.length,
      is_active: true,
    })
    setDialogOpen(true)
  }

  const handleEdit = (engine: SearchEngine) => {
    setEditingEngine(engine)
    setFormData({
      name: engine.name,
      url: engine.url,
      icon_url: engine.icon_url,
      sort_order: engine.sort_order,
      is_active: engine.is_active,
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
      await searchEngineApi.delete(deletingId)
      fetchEngines()
      toast({
        title: "删除成功",
        description: "搜索引擎已成功删除",
        variant: "success",
      })
    } catch (error) {
      console.error("删除失败:", error)
      toast({
        title: "删除失败",
        description: error instanceof Error ? error.message : "删除搜索引擎时发生错误",
        variant: "destructive",
      })
    } finally {
      setDeletingId(null)
    }
  }

  const handleSubmit = async () => {
    try {
      if (editingEngine) {
        await searchEngineApi.update(editingEngine.id, formData)
      } else {
        await searchEngineApi.create(formData)
      }
      setDialogOpen(false)
      fetchEngines()
      toast({
        title: "保存成功",
        description: editingEngine ? "搜索引擎已更新" : "搜索引擎已创建",
        variant: "success",
      })
    } catch (error) {
      console.error("保存失败:", error)
      toast({
        title: "保存失败",
        description: error instanceof Error ? error.message : "保存搜索引擎时发生错误",
        variant: "destructive",
      })
    }
  }

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    try {
      const data = await iconApi.upload(file)
      setFormData({ ...formData, icon_url: data.url })
    } catch (error) {
      console.error("上传失败:", error)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">搜索引擎管理</h1>
          <p className="text-gray-500 mt-2">拖动调整顺序，管理可用的搜索引擎</p>
        </div>
        <Button onClick={handleCreate} className="bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 text-white">
          <Plus className="h-4 w-4 mr-2" />
          添加搜索引擎
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
                <th className="p-4 text-left text-gray-700 font-medium w-16">图标</th>
                <th className="p-4 text-left text-gray-700 font-medium">名称</th>
                <th className="p-4 text-left text-gray-700 font-medium">搜索 URL</th>
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
              ) : engines.length === 0 ? (
                <tr>
                  <td colSpan={6} className="text-center py-8 text-gray-500">
                    暂无数据
                  </td>
                </tr>
              ) : (
                <SortableContext
                  items={engines.map((e) => e.id)}
                  strategy={verticalListSortingStrategy}
                >
                  {engines.map((engine) => (
                    <SortableRow
                      key={engine.id}
                      engine={engine}
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
              {editingEngine ? "编辑搜索引擎" : "添加搜索引擎"}
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
                placeholder="请输入名称"
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">搜索 URL *</Label>
              <Input
                value={formData.url}
                onChange={(e) =>
                  setFormData({ ...formData, url: e.target.value })
                }
                placeholder="https://www.google.com/search?q="
                className="bg-white text-gray-900 border-gray-200 placeholder:text-gray-400"
              />
              <p className="text-xs text-gray-500">
                搜索关键词会追加到 URL 末尾
              </p>
            </div>
            <div className="space-y-2">
              <Label className="text-gray-700">图标</Label>
              <div className="flex gap-2">
                <Input
                  value={formData.icon_url}
                  onChange={(e) =>
                    setFormData({ ...formData, icon_url: e.target.value })
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

      {/* 删除确认对话框 */}
      <ConfirmDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="删除搜索引擎"
        description="确定要删除这个搜索引擎吗？此操作无法撤销。"
        onConfirm={confirmDelete}
        confirmText="删除"
        cancelText="取消"
        variant="destructive"
      />
    </div>
  )
}
