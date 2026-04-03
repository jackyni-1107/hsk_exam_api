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
