import request from './request'

export interface LoginLogItem {
  id: number
  log_type: string
  user_id: number
  username: string
  user_type: number
  ip: string
  user_agent: string
  device_info: string
  trace_id: string
  fail_reason: string
  create_time: string
}

export function fetchLoginLogList(params: {
  page?: number
  size?: number
  username?: string
  log_type?: string
  user_type?: number
  start_time?: string
  end_time?: string
}) {
  const q: Record<string, unknown> = {
    page: params.page ?? 1,
    size: params.size ?? 10,
  }
  if (params.username) q.username = params.username
  if (params.log_type) q.log_type = params.log_type
  if (params.user_type != null && params.user_type >= 1) q.user_type = params.user_type
  if (params.start_time) q.start_time = params.start_time
  if (params.end_time) q.end_time = params.end_time
  return request.get<unknown, { data?: { list?: LoginLogItem[]; total?: number } }>('/admin/login-log/list', {
    params: q,
  })
}

/** 与视图层 import 命名一致 */
export const getLoginLogList = fetchLoginLogList
