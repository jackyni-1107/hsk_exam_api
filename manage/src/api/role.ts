import request from './request'

export interface RoleItem {
  id: number
  name: string
  code: string
  status: number
  sort: number
  type: number
  remark: string
  menu_ids: number[]
  create_time: string
}

/** 用户管理页：下拉用，仅正常状态、一页拉全 */
export function fetchRoleList(params?: { page?: number; size?: number; name?: string; status?: number }) {
  return request.get<unknown, { data?: { list?: RoleItem[]; total?: number } }>('/admin/role/list', {
    params: { page: 1, size: 500, status: 0, ...params },
  })
}

/** 角色管理页：分页列表 */
export function fetchRolePage(params: {
  page?: number
  size?: number
  name?: string
  /** -1 表示不传 status，查全部 */
  status?: number
}) {
  const q: Record<string, unknown> = {
    page: params.page ?? 1,
    size: params.size ?? 10,
  }
  if (params.name) q.name = params.name
  if (params.status !== undefined && params.status >= 0) q.status = params.status
  return request.get<unknown, { data?: { list?: RoleItem[]; total?: number } }>('/admin/role/list', { params: q })
}

export function createRole(payload: {
  name: string
  code: string
  status?: number
  sort?: number
  type?: number
  remark?: string
}) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/role', payload)
}

export function updateRole(
  id: number,
  payload: {
    name?: string
    code?: string
    status?: number
    sort?: number
    type?: number
    remark?: string
  }
) {
  return request.put<unknown, unknown>(`/admin/role/${id}`, payload)
}

export function assignRoleMenus(id: number, menu_ids: number[]) {
  return request.post<unknown, unknown>(`/admin/role/${id}/menus`, { menu_ids })
}

export function deleteRole(id: number) {
  return request.delete<unknown, unknown>(`/admin/role/${id}`)
}
