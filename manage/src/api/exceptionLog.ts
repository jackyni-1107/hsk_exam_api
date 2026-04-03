import request from './request'

export interface ExceptionLogItem {
  id: number
  trace_id: string
  path: string
  method: string
  error_msg: string
  stack: string
  user_id: number
  ip: string
  create_time: string
}

export function fetchExceptionLogList(params: {
  page?: number
  size?: number
  trace_id?: string
  path?: string
  start_time?: string
  end_time?: string
}) {
  return request.get<unknown, { data?: { list?: ExceptionLogItem[]; total?: number } }>(
    '/admin/exception-log/list',
    {
      params: {
        page: params.page ?? 1,
        size: params.size ?? 10,
        trace_id: params.trace_id,
        path: params.path,
        start_time: params.start_time,
        end_time: params.end_time,
      },
    }
  )
}
