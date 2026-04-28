import { defineComponent, h } from 'vue'
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { MENU_TYPE_BUTTON, type MenuTreeNode } from '@/api/menu'

const viewModules = import.meta.glob('../views/**/*.vue')
let dynamicRoutesAdded = false
let dynamicRouteNames: Array<string | symbol> = []

const NoPermission = defineComponent({
  name: 'NoPermission',
  setup() {
    return () =>
      h('div', { class: 'empty-permission' }, [
        h('h3', '暂无可访问菜单'),
        h('p', '当前账号未分配任何管理端菜单权限，请联系管理员。'),
      ])
  },
})

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/index.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      name: 'AdminRoot',
      component: () => import('@/layouts/AdminLayout.vue'),
      children: [
        {
          path: 'no-permission',
          name: 'NoPermission',
          component: NoPermission,
          meta: { title: '暂无权限' },
        },
      ],
    },
  ],
})

function routePath(path: string) {
  return path.replace(/^\/+/, '')
}

function normalizeComponent(component: string, path: string) {
  const raw = component.trim().replace(/^\/+/, '').replace(/\.vue$/, '')
  const withoutViews = raw.replace(/^views\//, '')
  if (withoutViews && withoutViews !== 'Layout') return withoutViews
  return path.replace(/^\/+/, '').replace(/\/$/, '') + '/index'
}

function componentLoader(node: MenuTreeNode) {
  const normalized = normalizeComponent(node.component || '', node.path || '')
  const key = `../views/${normalized}.vue`
  const loader = viewModules[key]
  if (!loader) {
    console.warn(`[router] menu component not found: ${key}`)
    return null
  }
  return loader
}

function buildRoutesFromMenus(nodes: MenuTreeNode[]): RouteRecordRaw[] {
  const routes: RouteRecordRaw[] = []
  const walk = (items: MenuTreeNode[]) => {
    for (const item of items) {
      if (item.type === MENU_TYPE_BUTTON) continue
      const isLayoutDirectory =
        item.children?.length && (!item.component || item.component.toLowerCase() === 'layout')
      if (item.path && item.visible !== false && !isLayoutDirectory) {
        const loader = componentLoader(item)
        if (loader) {
          routes.push({
            path: routePath(item.path),
            name: `MenuRoute${item.id}`,
            component: loader,
            meta: {
              title: item.name,
              permission: item.permission,
              keepAlive: item.keep_alive,
            },
          })
        }
      }
      if (item.children?.length) walk(item.children)
    }
  }
  walk(nodes)
  return routes
}

function firstMenuPath(nodes: MenuTreeNode[]): string {
  for (const node of nodes) {
    if (node.type !== MENU_TYPE_BUTTON && node.visible !== false && node.path) {
      const isLayoutDirectory =
        node.children?.length && (!node.component || node.component.toLowerCase() === 'layout')
      if (!isLayoutDirectory) {
        const loader = componentLoader(node)
        if (loader) return node.path
      }
    }
    if (node.children?.length) {
      const childPath = firstMenuPath(node.children)
      if (childPath) return childPath
    }
  }
  return ''
}

async function ensureDynamicRoutes() {
  const userStore = useUserStore()
  if (dynamicRoutesAdded && !userStore.menusLoaded) {
    for (const name of dynamicRouteNames) {
      if (router.hasRoute(name)) router.removeRoute(name)
    }
    dynamicRouteNames = []
    dynamicRoutesAdded = false
  }
  const menus = await userStore.loadMenus()
  if (!dynamicRoutesAdded) {
    for (const route of buildRoutesFromMenus(menus)) {
      router.addRoute('AdminRoot', route)
      if (typeof route.name === 'string' || typeof route.name === 'symbol') {
        dynamicRouteNames.push(route.name)
      }
    }
    dynamicRoutesAdded = true
    return true
  }
  return false
}

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()
  if (to.meta.public) {
    if (userStore.token) {
      await ensureDynamicRoutes()
      next({ path: firstMenuPath(userStore.rawMenus) || '/no-permission' })
    } else {
      next()
    }
    return
  }

  if (!userStore.token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
    return
  }

  const justAddedRoutes = await ensureDynamicRoutes()
  if (justAddedRoutes && to.path !== '/') {
    next({ ...to, replace: true })
    return
  }

  if (to.path === '/') {
    next({ path: firstMenuPath(userStore.rawMenus) || '/no-permission' })
    return
  }

  if (to.name == null && to.matched.length <= 1) {
    next({ path: firstMenuPath(userStore.rawMenus) || '/no-permission' })
    return
  }

  const permission = to.meta.permission as string | undefined
  if (permission && !userStore.hasPermission(permission)) {
    next({ path: firstMenuPath(userStore.rawMenus) || '/no-permission' })
    return
  }

  next()
})

export default router
