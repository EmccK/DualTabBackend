"use client"

import * as React from "react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { AlertTriangle, AlertCircle } from "lucide-react"

interface ConfirmDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  title: string
  description: string
  onConfirm: () => void
  confirmText?: string
  cancelText?: string
  variant?: "default" | "destructive"
}

export function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  onConfirm,
  confirmText = "确认",
  cancelText = "取消",
  variant = "default",
}: ConfirmDialogProps) {
  const handleConfirm = () => {
    onConfirm()
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[440px] bg-white border-gray-200">
        <DialogHeader>
          <div className="flex items-start gap-4">
            {/* 图标区域 */}
            <div className={`flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full ${
              variant === "destructive"
                ? "bg-red-50"
                : "bg-orange-50"
            }`}>
              {variant === "destructive" ? (
                <AlertTriangle className="h-6 w-6 text-red-600" />
              ) : (
                <AlertCircle className="h-6 w-6 text-orange-600" />
              )}
            </div>
            {/* 文字区域 */}
            <div className="flex-1 pt-1">
              <DialogTitle className="text-lg font-semibold text-gray-900">
                {title}
              </DialogTitle>
              <DialogDescription className="text-sm text-gray-500 mt-2 leading-relaxed">
                {description}
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>
        <DialogFooter className="mt-6 gap-3 sm:gap-3">
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            className="flex-1 sm:flex-1 border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400"
          >
            {cancelText}
          </Button>
          <Button
            onClick={handleConfirm}
            variant={variant === "destructive" ? "destructive" : "default"}
            className={`flex-1 sm:flex-1 ${
              variant === "destructive"
                ? "bg-red-500 hover:bg-red-600 shadow-sm"
                : "bg-gradient-to-r from-orange-500 to-amber-500 hover:from-orange-600 hover:to-amber-600 shadow-sm"
            }`}
          >
            {confirmText}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
