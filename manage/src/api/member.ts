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
