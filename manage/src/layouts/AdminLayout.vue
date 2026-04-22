<template>
  <el-container class="layout">
    <el-aside width="220px" class="aside">
      <div class="logo">
        <span class="logo-mark">
          <el-icon :size="22"><Grid /></el-icon>
        </span>
        <span class="logo-text">Admin</span>
      </div>
      <el-menu
        v-loading="menuLoading"
        :default-active="route.path"
        router
        class="side-menu"
        background-color="transparent"
        text-color="#94a3b8"
        active-text-color="#f8fafc"
      >
        <template v-for="item in menuList" :key="item.id">
          <el-menu-item v-if="item.type === 2 && !item.children?.length" :index="item.path" class="menu-item-leaf">
            <MenuIcon :name="item.icon" :is-directory="false" />
            <span>{{ item.name }}</span>
          </el-menu-item>
          <el-sub-menu v-else :index="item.path || `sub-${item.id}`" :class="{ 'sub-menu--directory': item.type === 1 }">
            <template #title>
              <MenuIcon :name="item.icon" :is-directory="item.type === 1" />
              <span>{{ item.name }}</span>
            </template>
            <template v-if="item.children?.length">
              <MenuItem v-for="child in item.children" :key="child.id" :menu="child" />
            </template>
          </el-sub-menu>
        </template>
      </el-menu>
    </el-aside>
    <el-container class="right-container">
      <el-header class="header">
        <div class="header-title">
          <span class="header-crumb">{{ pageTitle }}</span>
        </div>
        <el-dropdown trigger="click" @command="handleCommand">
          <span class="user-trigger">
            <el-avatar :size="36" class="user-avatar">{{ userInitial }}</el-avatar>
            <span class="user-name">{{ displayName }}</span>
            <el-icon class="user-caret"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main class="main">
        <div class="main-content">
          <router-view />
        </div>
        <footer class="layout-footer">© 2025 后台管理系统</footer>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Grid, ArrowDown } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { logout } from '@/api/auth'
import { getUserMenus, filterSidebarMenus } from '@/api/menu'
import type { MenuTreeNode } from '@/api/menu'
import MenuItem from './MenuItem.vue'
import MenuIcon from '@/components/MenuIcon.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const menuList = ref<MenuTreeNode[]>([])
const menuLoading = ref(false)

const displayName = computed(() => {
  const u = userStore.userInfo
  const n = u?.nickname as string | undefined
  const un = u?.username as string | undefined
  return n || un || 'Admin'
})

const userInitial = computed(() => displayName.value.slice(0, 1).toUpperCase())

const pageTitle = computed(() => (route.meta?.title as string) || '')

onMounted(async () => {
  menuLoading.value = true
  try {
    const res = (await getUserMenus()) as { data?: { list?: MenuTreeNode[] } }
    const list = filterSidebarMenus(res?.data?.list ?? [])
    const hasDashboardPath = (nodes: MenuTreeNode[]): boolean => {
      for (const n of nodes) {
        if (n.path === '/dashboard') return true
        if (n.children?.length && hasDashboardPath(n.children)) return true
      }
      return false
    }
    // 若数据库未配置「工作台」侧栏，则补一条（与 HLS 联调同级）
    if (!hasDashboardPath(list)) {
      list.unshift({
        id: -998,
        name: '工作台',
        permission: '',
        type: 2,
        sort: 0,
        parent_id: 0,
        path: '/dashboard',
        icon: 'Odometer',
        component: '',
        component_name: '',
        status: 1,
        visible: true,
        children: [],
      })
    }
    // 临时：HLS 联调入口（不写入菜单表，上线前可删）
    list.push({
      id: -999,
      name: 'HLS 联调',
      permission: '',
      type: 2,
      sort: 9999,
      parent_id: 0,
      path: '/exam/hls-debug',
      icon: 'VideoPlay',
      component: '',
      component_name: '',
      status: 1,
      visible: true,
      children: [],
    })

    menuList.value = list
  } finally {
    menuLoading.value = false
  }
})

function handleCommand(cmd: string) {
  if (cmd === 'logout') {
    logout()
    userStore.logout()
    menuList.value = []
    router.push('/login')
  }
}
</script>

<style scoped>
.layout {
  height: 100vh;
  overflow: hidden;
}

.aside {
  background: linear-gradient(180deg, #0f172a 0%, #1e293b 55%, #0f172a 100%);
  overflow-y: auto;
  border-right: 1px solid rgba(148, 163, 184, 0.12);
  box-shadow: 4px 0 24px rgba(15, 23, 42, 0.35);
}

.aside :deep(.el-menu) {
  border-right: none;
}

.side-menu {
  padding: 8px 10px 24px;
}

.aside :deep(.el-sub-menu .el-menu) {
  padding-left: 0;
  background-color: transparent !important;
}

.aside :deep(.el-menu-item),
.aside :deep(.el-sub-menu__title) {
  padding-left: 14px !important;
  height: 48px;
  line-height: 48px;
  border-radius: 10px;
  margin-bottom: 4px;
  transition:
    background 0.2s,
    color 0.2s;
}

.aside :deep(.el-menu-item:hover),
.aside :deep(.el-sub-menu__title:hover) {
  background: rgba(148, 163, 184, 0.12) !important;
  color: #e2e8f0 !important;
}

.aside :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.35), rgba(139, 92, 246, 0.22)) !important;
  color: #f8fafc !important;
  font-weight: 600;
  box-shadow: 0 4px 14px rgba(59, 130, 246, 0.2);
}

.aside :deep(.el-sub-menu .el-menu-item) {
  padding-left: 14px !important;
  min-width: auto;
  height: 44px;
  line-height: 44px;
}

.aside :deep(.el-sub-menu .el-sub-menu__title) {
  padding-left: 14px !important;
}

.aside :deep(.el-sub-menu .el-sub-menu .el-menu-item),
.aside :deep(.el-sub-menu .el-sub-menu .el-sub-menu__title) {
  padding-left: 14px !important;
}

.aside :deep(.el-sub-menu__icon-arrow) {
  right: 14px;
  color: #64748b;
}

.aside :deep(.sub-menu--directory > .el-sub-menu__title) {
  font-weight: 600;
  font-size: 14px;
  color: #cbd5e1;
}

.aside :deep(.sub-menu--directory > .el-sub-menu__title .menu-icon) {
  font-size: 18px;
}

.aside :deep(.menu-item-leaf .menu-icon) {
  font-size: 17px;
}

.right-container {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 0 16px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.12);
  flex-shrink: 0;
}

.logo-mark {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.45), rgba(139, 92, 246, 0.35));
  display: flex;
  align-items: center;
  justify-content: center;
  color: #f8fafc;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.25);
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.02em;
  color: #f1f5f9;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 56px !important;
  padding: 0 24px;
  background: #ffffff;
  border-bottom: 1px solid #e2e8f0;
  flex-shrink: 0;
  box-shadow: 0 1px 0 rgba(15, 23, 42, 0.04);
}

.header-title {
  min-width: 0;
}

.header-crumb {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
  letter-spacing: -0.01em;
}

.user-trigger {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 10px;
  margin-right: -10px;
  border-radius: 12px;
  transition: background 0.2s;
}

.user-trigger:hover {
  background: #f1f5f9;
}

.user-avatar {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  color: #fff;
  font-weight: 600;
  font-size: 15px;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #334155;
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-caret {
  font-size: 12px;
  color: #94a3b8;
}

.main {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #f1f5f9 0%, #e8edf3 100%);
  padding: 0;
}

.main-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
}

.main-content > * {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.layout-footer {
  flex-shrink: 0;
  height: 40px;
  line-height: 40px;
  padding: 0 24px;
  text-align: right;
  color: #94a3b8;
  font-size: 12px;
  background: #fff;
  border-top: 1px solid #e2e8f0;
}

.main-content :deep(.el-card) {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.06);
}

.main-content :deep(.el-card .el-card__header) {
  font-weight: 600;
  color: #0f172a;
  border-bottom-color: #e2e8f0;
}

.main-content :deep(.el-card .el-card__body) {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}
</style>
