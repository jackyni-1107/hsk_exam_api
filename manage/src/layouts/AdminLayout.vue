<template>
  <el-container class="layout">
    <el-aside width="220px" class="aside">
      <div class="logo">
        <span class="logo-mark">
          <el-icon :size="22"><Grid /></el-icon>
        </span>
        <span class="logo-text" :title="siteConfig.systemName">{{ siteConfig.systemName }}</span>
      </div>
      <el-menu
        v-loading="menuLoading"
        :default-active="route.path"
        router
        class="side-menu"
        background-color="transparent"
        text-color="#5c6b7a"
        active-text-color="#1d4ed8"
      >
        <template v-for="item in menuList" :key="item.id">
          <el-menu-item v-if="!item.children?.length" :index="item.path" class="menu-item-leaf">
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
        <footer class="layout-footer">{{ siteConfig.copyright }}</footer>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Grid, ArrowDown } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { logout } from '@/api/auth'
import MenuItem from './MenuItem.vue'
import MenuIcon from '@/components/MenuIcon.vue'
import { siteConfig } from '@/config/site'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const menuList = computed(() => userStore.sidebarMenus)
const menuLoading = computed(() => userStore.menusLoading)

const displayName = computed(() => {
  const u = userStore.userInfo
  const n = u?.nickname as string | undefined
  const un = u?.username as string | undefined
  return n || un || 'Admin'
})

const userInitial = computed(() => displayName.value.slice(0, 1).toUpperCase())

const pageTitle = computed(() => (route.meta?.title as string) || '')

function handleCommand(cmd: string) {
  if (cmd === 'logout') {
    logout()
    userStore.logout()
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
  /* 浅灰蓝侧栏，长时间浏览相对柔和 */
  background: linear-gradient(180deg, #f8fafc 0%, #f1f5f9 42%, #e8eef4 100%);
  overflow-y: auto;
  border-right: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 2px 0 18px rgba(15, 23, 42, 0.05);
  /* 隐藏滚动条，仍可用滚轮/触摸板滚动 */
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.aside::-webkit-scrollbar {
  width: 0;
  height: 0;
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
  background: rgba(37, 99, 235, 0.07) !important;
  color: #0f172a !important;
}

.aside :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.14), rgba(59, 130, 246, 0.1)) !important;
  color: #1d4ed8 !important;
  font-weight: 600;
  box-shadow: 0 1px 0 rgba(37, 99, 235, 0.12);
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
  color: #334155;
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
  justify-content: flex-start;
  gap: 10px;
  padding: 0 14px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  flex-shrink: 0;
}

.logo-mark {
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  border-radius: 10px;
  background: linear-gradient(135deg, #3b82f6, #6366f1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #f8fafc;
  box-shadow: 0 2px 10px rgba(59, 130, 246, 0.28);
}

.logo-text {
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.01em;
  color: #0f172a;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
