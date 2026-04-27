import request from './request'

export interface SecurityEventLogItem {
  id: number
  event_type: string
  user_id: number
  ip: string
  user_agent: string
  detail: string
  trace_id: string
  create_time: string
}

export function fetchSecurityEventLogList(params: {
  page?: number
  size?: number
  event_type?: string
  start_time?: string
  end_time?: string
}) {
  return request.get<unknown, { data?: { list?: SecurityEventLogItem[]; total?: number } }>(
    '/admin/log/security-event-log/list',
    {
      params: {
        page: params.page ?? 1,
        size: params.size ?? 10,
        event_type: params.event_type,
        start_time: params.start_time,
        end_time: params.end_time,
      },
    }
  )
}

/** 与视图层 import 命名一致 */
export const getSecurityEventLogList = fetchSecurityEventLogList
