import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { LoginUser } from '@/api/auth'
import { getUserMenus, filterSidebarMenus, type MenuTreeNode } from '@/api/menu'

const TOKEN_KEY = 'admin_token'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem(TOKEN_KEY) || '')
  const userInfo = ref<LoginUser | null>(null)
  const rawMenus = ref<MenuTreeNode[]>([])
  const sidebarMenus = ref<MenuTreeNode[]>([])
  const menusLoaded = ref(false)
  const menusLoading = ref(false)

  const permissions = computed(() => {
    const set = new Set<string>()
    const walk = (nodes: MenuTreeNode[]) => {
      for (const node of nodes) {
        if (node.permission) set.add(node.permission)
        if (node.children?.length) walk(node.children)
      }
    }
    walk(rawMenus.value)
    return set
  })

  function setSession(t: string, info: LoginUser | null) {
    token.value = t
    userInfo.value = info
    rawMenus.value = []
    sidebarMenus.value = []
    menusLoaded.value = false
    if (t) {
      localStorage.setItem(TOKEN_KEY, t)
    } else {
      localStorage.removeItem(TOKEN_KEY)
    }
  }

  async function loadMenus(force = false) {
    if (!token.value) return []
    if (menusLoaded.value && !force) return rawMenus.value
    menusLoading.value = true
    try {
      const res = (await getUserMenus()) as { data?: { list?: MenuTreeNode[] } }
      const list = res?.data?.list ?? []
      rawMenus.value = list
      sidebarMenus.value = filterSidebarMenus(list)
      menusLoaded.value = true
      return list
    } catch {
      rawMenus.value = []
      sidebarMenus.value = []
      menusLoaded.value = true
      return []
    } finally {
      menusLoading.value = false
    }
  }

  function hasPermission(permission?: string | string[]) {
    if (!permission) return true
    const required = Array.isArray(permission) ? permission : [permission]
    if (!required.length) return true
    const set = permissions.value
    return required.some((p) => set.has(p))
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    rawMenus.value = []
    sidebarMenus.value = []
    menusLoaded.value = false
    localStorage.removeItem(TOKEN_KEY)
  }

  return {
    token,
    userInfo,
    rawMenus,
    sidebarMenus,
    menusLoaded,
    menusLoading,
    permissions,
    setSession,
    loadMenus,
    hasPermission,
    logout,
  }
})
