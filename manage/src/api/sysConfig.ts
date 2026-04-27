import request from './request'

export interface ConfigItemRow {
  id: number
  config_key: string
  config_value: string
  config_type: string
  group_name: string
  remark: string
  create_time: string
}

export function fetchConfigList(params: {
  page?: number
  size?: number
  group?: string
  key?: string
}) {
  return request.get<unknown, { data?: { list?: ConfigItemRow[]; total?: number } }>('/admin/config/list', {
    params: { page: params.page ?? 1, size: params.size ?? 10, group: params.group, key: params.key },
  })
}

export function createConfig(payload: {
  config_key: string
  config_value: string
  config_type?: string
  group_name?: string
  remark?: string
}) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/config', payload)
}

export function updateConfig(id: number, payload: { config_value?: string; remark?: string }) {
  return request.put<unknown, unknown>(`/admin/config/${id}`, payload)
}

export function deleteConfig(id: number) {
  return request.delete<unknown, unknown>(`/admin/config/${id}`)
}
