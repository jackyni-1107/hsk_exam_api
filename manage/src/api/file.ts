import axios from 'axios'
import request from './request'

export interface FileItemRow {
  id: number
  filename: string
  path: string
  size: number
  mime_type: string
  is_private: number
  create_time: string
}

export interface StorageConfigItem {
  id: number
  storage_type: string
  name: string
  is_active: number
  config_json: string
  cleanup_before_days: number
  create_time: string
}

export function fetchFileList(params: { page?: number; size?: number; filename?: string }) {
  return request.get<unknown, { data?: { list?: FileItemRow[]; total?: number } }>('/admin/file/list', {
    params,
  })
}

export function deleteFile(id: number) {
  return request.delete<unknown, unknown>(`/admin/file/${id}`)
}

/** multipart 上传至文件中心（表单字段名 file） */
export function uploadSysFile(file: File, isPrivate = 0) {
  const fd = new FormData()
  fd.append('file', file)
  fd.append('is_private', String(isPrivate))
  return request.post<unknown, { data?: { id?: number; path?: string; filename?: string; size?: number; mime_type?: string } }>(
    '/admin/file/upload',
    fd
  )
}

/** 流式下载（不经统一 JSON 包装），成功时触发浏览器保存 */
export async function downloadSysFile(id: number, filename: string) {
  const token = localStorage.getItem('admin_token')
  const res = await axios.get(`/api/admin/file/${id}/download`, {
    responseType: 'blob',
    timeout: 120000,
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    validateStatus: () => true,
  })
  if (res.status === 401) {
    throw new Error('未登录或登录已过期')
  }
  const ct = res.headers['content-type'] || ''
  if (ct.includes('application/json')) {
    const text = await (res.data as Blob).text()
    let msg = '下载失败'
    try {
      const j = JSON.parse(text) as { message?: string }
      if (j.message) msg = j.message
    } catch {
      /* */
    }
    throw new Error(msg)
  }
  if (res.status !== 200) {
    throw new Error(`下载失败 (${res.status})`)
  }
  const blob = res.data as Blob
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename || 'download'
  a.rel = 'noopener'
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

export function fetchStorageConfigs() {
  return request.get<unknown, { data?: { list?: StorageConfigItem[] } }>('/admin/file/storage-config/list')
}

export function createStorageConfig(payload: {
  storage_type: string
  name: string
  config_json: string
  cleanup_before_days?: number
}) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/file/storage-config', payload)
}

export function updateStorageConfig(
  id: number,
  payload: { name?: string; config_json?: string; cleanup_before_days?: number }
) {
  return request.put<unknown, unknown>(`/admin/file/storage-config/${id}`, payload)
}

export function deleteStorageConfig(id: number) {
  return request.delete<unknown, unknown>(`/admin/file/storage-config/${id}`)
}

export function setActiveStorageConfig(id: number) {
  return request.post<unknown, unknown>(`/admin/file/storage-config/${id}/set-active`)
}
