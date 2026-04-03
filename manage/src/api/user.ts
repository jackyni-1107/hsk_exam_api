import request from './request'

export interface AdminUserItem {
  id: number
  username: string
  nickname: string
  email: string
  mobile: string
  status: number
  role_ids: number[]
  create_time: string
}

export function fetchUserList(params: {
  page?: number
  size?: number
  username?: string
  status?: number
}) {
  return request.get<unknown, { data?: { list?: AdminUserItem[]; total?: number } }>(
    '/admin/user/list',
    { params }
  )
}

export function createUser(payload: {
  username: string
  password: string
  nickname?: string
  email?: string
  mobile?: string
  status?: number
  role_ids?: number[]
}) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/user', payload)
}

export function updateUser(
  id: number,
  payload: {
    password?: string
    nickname?: string
    email?: string
    mobile?: string
    status?: number
  }
) {
  return request.put<unknown, unknown>(`/admin/user/${id}`, payload)
}

export function assignUserRoles(id: number, role_ids: number[]) {
  return request.post<unknown, unknown>(`/admin/user/${id}/roles`, { role_ids })
}

export function deleteUser(id: number) {
  return request.delete<unknown, unknown>(`/admin/user/${id}`)
}

export function kickUserSessions(id: number) {
  return request.post<unknown, unknown>(`/admin/user/${id}/kick-sessions`)
}
