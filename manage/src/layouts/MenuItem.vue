<template>
  <template v-if="menu.type === 2 && !menu.children?.length">
    <el-menu-item :index="menu.path" class="menu-item-leaf">
      <MenuIcon :name="menu.icon" :is-directory="false" />
      <span>{{ menu.name }}</span>
    </el-menu-item>
  </template>
  <el-sub-menu v-else :index="menu.path || `sub-${menu.id}`" :class="{ 'sub-menu--directory': menu.type === 1 }">
    <template #title>
      <MenuIcon :name="menu.icon" :is-directory="menu.type === 1" />
      <span>{{ menu.name }}</span>
    </template>
    <MenuItem v-for="child in menu.children" :key="child.id" :menu="child" />
  </el-sub-menu>
</template>

<script setup lang="ts">
import type { MenuTreeNode } from '@/api/menu'
import MenuIcon from '@/components/MenuIcon.vue'

defineProps<{
  menu: MenuTreeNode
}>()
</script>
