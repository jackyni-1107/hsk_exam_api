import axios from 'axios'
import request from './request'

export interface MemberItem {
  id: number
  username: string
  nickname: string
  email: string
  mobile: string
  status: number
  create_time: string
}

export function fetchMemberList(params: {
  page?: number
  size?: number
  username?: string
  status?: number
}) {
  const q: Record<string, unknown> = { page: params.page ?? 1, size: params.size ?? 10 }
  if (params.username) q.username = params.username
  if (params.status !== undefined && params.status >= 0) q.status = params.status
  return request.get<unknown, { data?: { list?: MemberItem[]; total?: number } }>('/admin/member/list', {
    params: q,
  })
}

export function createMember(payload: {
  username: string
  password: string
  nickname?: string
  email?: string
  mobile?: string
  status?: number
}) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/member', payload)
}

export function updateMember(
  id: number,
  payload: { password?: string; nickname?: string; email?: string; mobile?: string; status?: number }
) {
  return request.put<unknown, unknown>(`/admin/member/${id}`, payload)
}

export function deleteMember(id: number) {
  return request.delete<unknown, unknown>(`/admin/member/${id}`)
}

export interface MemberImportResult {
  total: number
  success: number
  failed: number
  errors: string[]
}

/** multipart 导入客户（表单字段名 file，CSV UTF-8） */
export function importMembersCsv(file: File) {
  const fd = new FormData()
  fd.append('file', file)
  return request.post<unknown, { data?: MemberImportResult }>('/admin/member/import', fd)
}

/** 下载客户导入 CSV 模板（不经统一 JSON 包装） */
export async function downloadMemberImportTemplate() {
  const token = localStorage.getItem('admin_token')
  const res = await axios.get('/api/admin/member/import-template', {
    responseType: 'blob',
    timeout: 60000,
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
  a.download = '客户导入模板.csv'
  a.rel = 'noopener'
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}
