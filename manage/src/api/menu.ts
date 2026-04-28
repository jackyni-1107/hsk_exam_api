import request from './request'

/** 与后端 system_menu.type 一致：1 目录 2 菜单 3 按钮（侧栏不展示） */
export const MENU_TYPE_DIR = 1
export const MENU_TYPE_MENU = 2
export const MENU_TYPE_BUTTON = 3

export interface MenuTreeNode {
  id: number
  name: string
  permission: string
  type: number
  sort: number
  parent_id: number
  path: string
  icon: string
  component: string
  component_name: string
  status: number
  visible: boolean
  keep_alive?: boolean
  always_show?: boolean
  children?: MenuTreeNode[]
}

/** 侧栏用：去掉按钮节点，并去掉过滤后无子级的目录 */
export function filterSidebarMenus(nodes: MenuTreeNode[]): MenuTreeNode[] {
  const mapped = nodes
    .filter((n) => n.type !== MENU_TYPE_BUTTON)
    .map((n) => {
      const children = n.children?.length ? filterSidebarMenus(n.children) : undefined
      return { ...n, children }
    })
  return mapped.filter((n) => n.type !== MENU_TYPE_DIR || !!n.path || (n.children?.length ?? 0) > 0)
}

export function getUserMenus() {
  return request.get<unknown, { data?: { list?: MenuTreeNode[] } }>('/admin/me/menus')
}

/** 管理端完整菜单树（含按钮），用于角色授权等 */
export function fetchAdminMenuTree() {
  return request.get<unknown, { data?: { list?: MenuTreeNode[] } }>('/admin/menu/tree')
}

export function createAdminMenu(payload: Record<string, unknown>) {
  return request.post<unknown, { data?: { id?: number } }>('/admin/menu', payload)
}

export function updateAdminMenu(id: number, payload: Record<string, unknown>) {
  return request.put<unknown, unknown>(`/admin/menu/${id}`, payload)
}

export function deleteAdminMenu(id: number) {
  return request.delete<unknown, unknown>(`/admin/menu/${id}`)
}
