// 运行时环境变量类型定义
declare global {
  interface Window {
    __ENV__?: {
      API_URL?: string
    }
  }
}

// API 基础配置
// 标准做法: 运行时环境变量注入(window.__ENV__)
const API_BASE_URL = (() => {
  // 1. 优先使用运行时环境变量(Docker 容器启动时注入)
  if (typeof window !== 'undefined' && window.__ENV__?.API_URL) {
    return window.__ENV__.API_URL
  }

  // 2. 使用构建时环境变量(开发环境)
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL
  }

  // 3. 默认值(本地开发)
  return 'http://localhost:8080'
})()

// 获取 Token
function getToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('token')
}

// 设置 Token
export function setToken(token: string): void {
  localStorage.setItem('token', token)
}

// 清除 Token
export function clearToken(): void {
  localStorage.removeItem('token')
}

// 请求配置
interface RequestConfig extends RequestInit {
  params?: Record<string, string | number | undefined>
}

// API 响应类型
interface ApiResponse<T = unknown> {
  msg: string
  data: T
}

// 分页响应类型
export interface PageData<T> {
  list: T[]
  total: number
  page: number
  size: number
}

// 通用请求函数
async function request<T>(endpoint: string, config: RequestConfig = {}): Promise<T> {
  const { params, ...init } = config

  // 构建 URL
  let url = `${API_BASE_URL}${endpoint}`
  if (params) {
    const searchParams = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        searchParams.append(key, String(value))
      }
    })
    const queryString = searchParams.toString()
    if (queryString) {
      url += `?${queryString}`
    }
  }

  // 设置请求头
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...init.headers,
  }

  // 添加认证 Token
  const token = getToken()
  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(url, {
    ...init,
    headers,
  })

  const data: ApiResponse<T> = await response.json()

  if (!response.ok) {
    throw new Error(data.msg || '请求失败')
  }

  return data.data
}

// ========== 认证 API ==========

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    id: number
    username: string
  }
}

export interface UserInfo {
  id: number
  username: string
}

export const authApi = {
  // 登录
  login: (data: LoginRequest) =>
    request<LoginResponse>('/admin/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 获取当前用户
  me: () => request<UserInfo>('/admin/auth/me'),

  // 修改密码
  changePassword: (data: { old_password: string; new_password: string }) =>
    request<null>('/admin/auth/password', {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
}

// ========== 分类 API ==========

export interface Category {
  id: number
  name: string
  name_en: string
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateCategoryRequest {
  name: string
  name_en?: string
  sort_order?: number
  is_active?: boolean
}

export interface UpdateCategoryRequest {
  name?: string
  name_en?: string
  sort_order?: number
  is_active?: boolean
}

export const categoryApi = {
  // 获取列表
  list: () => request<{ list: Category[] }>('/admin/categories'),

  // 获取详情
  get: (id: number) => request<Category>(`/admin/categories/${id}`),

  // 创建
  create: (data: CreateCategoryRequest) =>
    request<Category>('/admin/categories', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 更新
  update: (id: number, data: UpdateCategoryRequest) =>
    request<Category>(`/admin/categories/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除
  delete: (id: number) =>
    request<null>(`/admin/categories/${id}`, {
      method: 'DELETE',
    }),
}

// ========== 图标 API ==========

export interface Icon {
  id: number
  uuid: string
  title: string
  description: string
  url: string
  img_url: string
  bg_color: string
  mime_type: string
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
  categories?: Category[] // 多分类
}

export interface IconListParams {
  page?: number
  size?: number
  category_id?: number
  keyword?: string
  [key: string]: string | number | undefined
}

export interface CreateIconRequest {
  title: string
  description?: string
  url: string
  img_url?: string
  bg_color?: string
  mime_type?: string
  category_ids?: number[] // 多分类
  sort_order?: number
  is_active?: boolean
}

export interface UpdateIconRequest {
  title?: string
  description?: string
  url?: string
  img_url?: string
  bg_color?: string
  mime_type?: string
  category_ids?: number[] // 多分类
  sort_order?: number
  is_active?: boolean
}

export const iconApi = {
  // 获取列表
  list: (params?: IconListParams) =>
    request<PageData<Icon>>('/admin/icons', { params }),

  // 获取详情
  get: (id: number) => request<Icon>(`/admin/icons/${id}`),

  // 创建
  create: (data: CreateIconRequest) =>
    request<Icon>('/admin/icons', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 更新
  update: (id: number, data: UpdateIconRequest) =>
    request<Icon>(`/admin/icons/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除
  delete: (id: number) =>
    request<null>(`/admin/icons/${id}`, {
      method: 'DELETE',
    }),

  // 上传图标
  upload: async (file: File): Promise<{ url: string }> => {
    const formData = new FormData()
    formData.append('file', file)

    const token = getToken()
    const response = await fetch(`${API_BASE_URL}/admin/upload/icon`, {
      method: 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: formData,
    })

    const data: ApiResponse<{ url: string }> = await response.json()
    if (!response.ok) {
      throw new Error(data.msg || '上传失败')
    }
    return data.data
  },
}

// ========== 搜索引擎 API ==========

export interface SearchEngine {
  id: number
  uuid: string
  name: string
  url: string
  icon_url: string
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateSearchEngineRequest {
  name: string
  url: string
  icon_url?: string
  sort_order?: number
  is_active?: boolean
}

export interface UpdateSearchEngineRequest {
  name?: string
  url?: string
  icon_url?: string
  sort_order?: number
  is_active?: boolean
}

export const searchEngineApi = {
  // 获取列表
  list: () => request<{ list: SearchEngine[] }>('/admin/search-engines'),

  // 获取详情
  get: (id: number) => request<SearchEngine>(`/admin/search-engines/${id}`),

  // 创建
  create: (data: CreateSearchEngineRequest) =>
    request<SearchEngine>('/admin/search-engines', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 更新
  update: (id: number, data: UpdateSearchEngineRequest) =>
    request<SearchEngine>(`/admin/search-engines/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除
  delete: (id: number) =>
    request<null>(`/admin/search-engines/${id}`, {
      method: 'DELETE',
    }),
}

// ========== 壁纸 API ==========

export interface Wallpaper {
  id: number
  uuid: string
  title: string
  url: string
  thumb_url: string
  source: string
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface WallpaperListParams {
  page?: number
  size?: number
  [key: string]: string | number | undefined
}

export interface CreateWallpaperRequest {
  title: string
  url: string
  thumb_url?: string
  source?: string
  sort_order?: number
  is_active?: boolean
}

export interface UpdateWallpaperRequest {
  title?: string
  url?: string
  thumb_url?: string
  source?: string
  sort_order?: number
  is_active?: boolean
}

export const wallpaperApi = {
  // 获取列表
  list: (params?: WallpaperListParams) =>
    request<PageData<Wallpaper>>('/admin/wallpapers', { params }),

  // 获取详情
  get: (id: number) => request<Wallpaper>(`/admin/wallpapers/${id}`),

  // 创建
  create: (data: CreateWallpaperRequest) =>
    request<Wallpaper>('/admin/wallpapers', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // 更新
  update: (id: number, data: UpdateWallpaperRequest) =>
    request<Wallpaper>(`/admin/wallpapers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  // 删除
  delete: (id: number) =>
    request<null>(`/admin/wallpapers/${id}`, {
      method: 'DELETE',
    }),

  // 上传壁纸
  upload: async (file: File): Promise<{ url: string }> => {
    const formData = new FormData()
    formData.append('file', file)

    const token = getToken()
    const response = await fetch(`${API_BASE_URL}/admin/upload/wallpaper`, {
      method: 'POST',
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: formData,
    })

    const data: ApiResponse<{ url: string }> = await response.json()
    if (!response.ok) {
      throw new Error(data.msg || '上传失败')
    }
    return data.data
  },
}

// ========== 系统配置 API ==========

export interface SystemConfig {
  id: number
  key: string
  value: string
  remark: string
  created_at: string
  updated_at: string
}

export interface ConfigKeyInfo {
  key: string
  remark: string
  default_value: string
}

export const configApi = {
  // 获取配置列表
  list: () => request<{ list: SystemConfig[] }>('/admin/configs'),

  // 获取可用配置项说明
  getKeys: () => request<{ list: ConfigKeyInfo[] }>('/admin/configs/keys'),

  // 获取单个配置
  get: (key: string) => request<SystemConfig>(`/admin/configs/${key}`),

  // 设置配置
  set: (key: string, value: string) =>
    request<SystemConfig>('/admin/configs', {
      method: 'POST',
      body: JSON.stringify({ key, value }),
    }),

  // 批量设置配置
  batchSet: (configs: { key: string; value: string }[]) =>
    request<null>('/admin/configs/batch', {
      method: 'POST',
      body: JSON.stringify({ configs }),
    }),

  // 删除配置
  delete: (key: string) =>
    request<null>(`/admin/configs/${key}`, {
      method: 'DELETE',
    }),
}
