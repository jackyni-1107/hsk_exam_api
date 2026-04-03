import request from './request'

export interface NotificationLogItem {
  id: number
  template_code: string
  channel: string
  recipient: string
  status: number
  error_msg: string
  create_time: string
}

export interface ChannelConfigItem {
  id: number
  channel: string
  provider: string
  name: string
  is_active: number
  config_json: string
  create_time: string
}

export interface TemplateItem {
  id: number
  code: string
  name: string
  channel: string
  content: string
  variables: string
  status: number
  create_time: string
}

export function fetchNotificationLogs(params: {
  page?: number
  size?: number
  channel?: string
  recipient?: string
}) {
  return request.get<unknown, { data?: { list?: NotificationLogItem[]; total?: number } }>(
    '/admin/notification/log/list',
    { params }
  )
}

export function fetchChannelConfigs(params?: { channel?: string }) {
  return request.get<unknown, { data?: { list?: ChannelConfigItem[] } }>('/admin/notification/channel/list', {
    params,
  })
}

export function createChannelConfig(payload: {
  channel: string
  provider: string
  name: string
  config_json: string
}) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/notification/channel', payload)
}

export function updateChannelConfig(
  id: number,
  payload: { name?: string; config_json?: string }
) {
  return request.put<unknown, unknown>(`/admin/notification/channel/${id}`, payload)
}

export function deleteChannelConfig(id: number) {
  return request.delete<unknown, unknown>(`/admin/notification/channel/${id}`)
}

export function setActiveChannelConfig(id: number) {
  return request.post<unknown, unknown>(`/admin/notification/channel/${id}/set-active`)
}

export function fetchTemplates(params: { page?: number; size?: number; code?: string; channel?: string }) {
  return request.get<unknown, { data?: { list?: TemplateItem[]; total?: number } }>(
    '/admin/notification/template/list',
    { params }
  )
}

export function createTemplate(payload: Record<string, unknown>) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/notification/template', payload)
}

export function updateTemplate(id: number, payload: Record<string, unknown>) {
  return request.put<unknown, unknown>(`/admin/notification/template/${id}`, payload)
}

export function deleteTemplate(id: number) {
  return request.delete<unknown, unknown>(`/admin/notification/template/${id}`)
}

export function sendNotification(payload: {
  template_code: string
  channel: string
  recipient: string
  variables?: string
}) {
  return request.post<unknown, { data?: { ok?: boolean } }>('/admin/notification/send', payload)
}
