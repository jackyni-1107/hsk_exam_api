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

export interface DictTypeItem {
  id: number
  dict_name: string
  dict_type: string
  status: number
  create_time: string
}

export interface DictDataItem {
  id: number
  dict_type: string
  dict_label: string
  dict_value: string
  sort: number
  status: number
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

export function fetchDictTypeList(params: {
  page?: number
  size?: number
  dict_type?: string
}) {
  const q: Record<string, unknown> = { page: params.page ?? 1, size: params.size ?? 10 }
  if (params.dict_type) q.dict_type = params.dict_type
  return request.get<unknown, { data?: { list?: DictTypeItem[]; total?: number } }>('/admin/dict/type/list', {
    params: q,
  })
}

export function createDictType(payload: {
  dict_name: string
  dict_type: string
  status?: number
  remark?: string
}) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/dict/type', payload)
}

export function updateDictType(
  id: number,
  payload: { dict_name?: string; status?: number; remark?: string }
) {
  return request.put<unknown, unknown>(`/admin/dict/type/${id}`, payload)
}

export function deleteDictType(id: number) {
  return request.delete<unknown, unknown>(`/admin/dict/type/${id}`)
}

export function fetchDictDataList(params: { page?: number; size?: number; dict_type: string }) {
  return request.get<unknown, { data?: { list?: DictDataItem[]; total?: number } }>('/admin/dict/data/list', {
    params: { page: params.page ?? 1, size: params.size ?? 10, dict_type: params.dict_type },
  })
}

export function createDictData(payload: Record<string, unknown>) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/dict/data', payload)
}

export function updateDictData(id: number, payload: Record<string, unknown>) {
  return request.put<unknown, unknown>(`/admin/dict/data/${id}`, payload)
}

export function deleteDictData(id: number) {
  return request.delete<unknown, unknown>(`/admin/dict/data/${id}`)
}
