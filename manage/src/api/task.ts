import request from './request'

export interface TaskItem {
  id: number
  name: string
  code: string
  type: number
  cron_expr: string
  delay_seconds: number
  handler: string
  params: string
  retry_times: number
  retry_interval: number
  concurrency: number
  alert_on_fail: number
  alert_receivers: string
  status: number
  remark: string
  create_time: string
}

export interface TaskLogItem {
  id: number
  task_id: number
  run_id: string
  trigger_type: number
  status: number
  start_time: string
  end_time: string
  duration_ms: number
  retry_count: number
  error_msg: string
  result: string
  node: string
  create_time: string
}

export function fetchTaskList(params: {
  page?: number
  size?: number
  name?: string
  code?: string
  type?: number
  status?: number | null
  handler?: string
}) {
  const q: Record<string, unknown> = {
    page: params.page ?? 1,
    size: params.size ?? 10,
  }
  if (params.name) q.name = params.name
  if (params.code) q.code = params.code
  if (params.type != null && params.type > 0) q.type = params.type
  if (params.status != null && params.status >= 0) q.status = params.status
  if (params.handler) q.handler = params.handler
  return request.get<unknown, { data?: { list?: TaskItem[]; total?: number } }>('/admin/task/list', { params: q })
}

export function createTask(payload: Record<string, unknown>) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/task', payload)
}

export function updateTask(id: number, payload: Record<string, unknown>) {
  return request.put<unknown, unknown>(`/admin/task/${id}`, payload)
}

export function deleteTask(id: number) {
  return request.delete<unknown, unknown>(`/admin/task/${id}`)
}

export function runTask(id: number) {
  return request.post<unknown, { data?: { run_id?: string } }>('/admin/task/run', { id })
}

export function fetchTaskLogs(params: {
  page?: number
  size?: number
  task_id?: number
  run_id?: string
  status?: number | null
}) {
  const q: Record<string, unknown> = {
    page: params.page ?? 1,
    size: params.size ?? 10,
  }
  if (params.task_id != null && params.task_id > 0) q.task_id = params.task_id
  if (params.run_id) q.run_id = params.run_id
  if (params.status != null && params.status >= 0) q.status = params.status
  return request.get<unknown, { data?: { list?: TaskLogItem[]; total?: number } }>('/admin/task/log', { params: q })
}
