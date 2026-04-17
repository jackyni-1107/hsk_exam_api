import request from './request'

export interface AuditLogItem {
  id: number
  user_id: number
  username: string
  user_type: number
  module: string
  action: string
  log_type: string
  method: string
  path: string
  request_data: string
  response_data: string
  ip: string
  user_agent: string
  trace_id: string
  device_info: string
  duration_ms: number
  create_time: string
}

export interface AuditChangeDetailItem {
  id: number
  table_name: string
  record_id: number
  field_name: string
  before_value: string
  after_value: string
  create_time: string
}

export function getAuditLogList(params: {
  page?: number
  size?: number
  username?: string
  path?: string
  action?: string
  log_type?: string
  trace_id?: string
  start_time?: string
  end_time?: string
}) {
  return request.get<unknown, { data?: { list?: AuditLogItem[]; total?: number } }>('/admin/audit-log/list', {
    params: {
      page: params.page ?? 1,
      size: params.size ?? 10,
      username: params.username,
      path: params.path,
      action: params.action,
      log_type: params.log_type,
      trace_id: params.trace_id,
      start_time: params.start_time,
      end_time: params.end_time,
    },
  })
}

export function getAuditLogChangeDetails(id: number) {
  return request.get<unknown, { data?: { list?: AuditChangeDetailItem[] } }>(
    `/admin/audit-log/${id}/change-details`
  )
}
