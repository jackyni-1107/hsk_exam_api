<template>
  <el-icon v-if="iconComponent" class="menu-icon" :class="{ 'menu-icon--directory': isDirectory }">
    <component :is="iconComponent" />
  </el-icon>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Setting,
  User,
  UserFilled,
  Menu,
  Document,
  Key,
  Warning,
  Lock,
  List,
  Folder,
  Bell,
  Files,
  Grid,
  Operation,
  Message,
  Collection,
  Timer,
  Reading,
  Notebook,
  Finished,
  VideoPlay,
  Calendar,
  Odometer,
} from '@element-plus/icons-vue'

const iconMap: Record<string, unknown> = {
  Setting,
  User,
  UserFilled,
  Menu,
  Document,
  Key,
  Warning,
  Lock,
  List,
  Folder,
  Bell,
  Files,
  Grid,
  Operation,
  Message,
  Collection,
  Timer,
  Reading,
  Notebook,
  Finished,
  VideoPlay,
  Calendar,
  Odometer,
}

const props = defineProps<{
  name?: string
  isDirectory?: boolean
}>()

const iconComponent = computed(() => {
  if (!props.name) return null
  return (iconMap[props.name] as object) || null
})
</script>

<style scoped>
.menu-icon {
  margin-right: 8px;
  font-size: 18px;
  vertical-align: middle;
}
.menu-icon--directory {
  font-size: 20px;
  color: #e6a23c;
}
</style>
