<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>菜单管理</span>
          <div>
            <el-button @click="loadTree">刷新</el-button>
            <el-button type="primary" @click="openCreate(0)">新增根菜单</el-button>
          </div>
        </div>
      </template>
      <el-table v-loading="loading" :data="flatRows" row-key="id" border default-expand-all :tree-props="{ children: 'children' }">
        <el-table-column prop="name" label="名称" min-width="200" />
        <el-table-column label="类型" width="88">
          <template #default="{ row }">
            <span v-if="row.type === MENU_TYPE_DIR">目录</span>
            <span v-else-if="row.type === MENU_TYPE_MENU">菜单</span>
            <span v-else>按钮</span>
          </template>
        </el-table-column>
        <el-table-column prop="permission" label="权限标识" min-width="140" show-overflow-tooltip />
        <el-table-column prop="path" label="路由" min-width="140" show-overflow-tooltip />
        <el-table-column prop="sort" label="排序" width="72" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">{{ row.status === 0 ? '正常' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.type !== MENU_TYPE_BUTTON" link type="primary" @click="openCreate(row.id)">子项</el-button>
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="visible" :title="mode === 'create' ? '新增菜单' : '编辑菜单'" width="560px" destroy-on-close @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="上级" prop="parent_id">
          <el-tree-select
            v-model="form.parent_id"
            :data="parentTreeData"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            check-strictly
            clearable
            :placeholder="parentPlaceholder"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :label="MENU_TYPE_DIR">目录</el-radio>
            <el-radio :label="MENU_TYPE_MENU">菜单</el-radio>
            <el-radio :label="MENU_TYPE_BUTTON">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="权限标识" prop="permission">
          <el-input
            v-model="form.permission"
            :disabled="isCurrentDir"
            :placeholder="isCurrentButton ? '如 user:create' : isCurrentDir ? '目录无需权限标识' : '如 user:list'"
          />
        </el-form-item>
        <el-form-item v-if="!isCurrentButton" label="路由 path" prop="path">
          <el-input v-model="form.path" :placeholder="isCurrentMenu ? '如 /system/user' : '目录可留空'" />
        </el-form-item>
        <el-form-item v-if="!isCurrentButton" label="图标" prop="icon">
          <el-input v-model="form.icon" placeholder="Element 图标名" />
        </el-form-item>
        <el-form-item v-if="isCurrentMenu" label="组件" prop="component">
          <el-input v-model="form.component" placeholder="如 views/system/user/index" />
        </el-form-item>
        <el-form-item v-if="isCurrentMenu" label="组件名" prop="component_name">
          <el-input v-model="form.component_name" placeholder="开启缓存时必填" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="0">正常</el-radio>
            <el-radio :label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="显示">
          <el-switch v-model="form.visible" />
        </el-form-item>
        <el-form-item v-if="isCurrentMenu" label="缓存">
          <el-switch v-model="form.keep_alive" />
        </el-form-item>
        <el-form-item v-if="!isCurrentButton" label="总显示">
          <el-switch v-model="form.always_show" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import type { FormInstance, FormItemRule, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchAdminMenuTree,
  createAdminMenu,
  updateAdminMenu,
  deleteAdminMenu,
  type MenuTreeNode,
} from '@/api/menu'

const MENU_TYPE_DIR = 1
const MENU_TYPE_MENU = 2
const MENU_TYPE_BUTTON = 3

type MenuForm = {
  parent_id: number | null
  name: string
  permission: string
  type: number
  sort: number
  path: string
  icon: string
  component: string
  component_name: string
  status: number
  visible: boolean
  keep_alive: boolean
  always_show: boolean
}

const loading = ref(false)
const saving = ref(false)
const treeData = ref<MenuTreeNode[]>([])
const visible = ref(false)
const mode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const editingId = ref(0)
const forbiddenParentIds = ref<Set<number>>(new Set())

const form = reactive<MenuForm>({
  parent_id: 0,
  name: '',
  permission: '',
  type: MENU_TYPE_DIR,
  sort: 0,
  path: '',
  icon: '',
  component: '',
  component_name: '',
  status: 0,
  visible: true,
  keep_alive: false,
  always_show: false,
})

const isCurrentDir = computed(() => form.type === MENU_TYPE_DIR)
const isCurrentMenu = computed(() => form.type === MENU_TYPE_MENU)
const isCurrentButton = computed(() => form.type === MENU_TYPE_BUTTON)
const flatRows = computed(() => treeData.value)
const parentPlaceholder = computed(() => (isCurrentButton.value ? '按钮必须选择上级目录或菜单' : '不选则为根（实际选“根”）'))

const validateName: FormItemRule['validator'] = (_rule, value, callback) => {
  if (!normalizeText(String(value ?? ''))) {
    callback(new Error('请输入名称'))
    return
  }
  callback()
}

const validateParentId: FormItemRule['validator'] = (_rule, value, callback) => {
  const parentId = value == null ? 0 : Number(value)
  if (form.type === MENU_TYPE_BUTTON && parentId === 0) {
    callback(new Error('按钮必须选择上级目录或菜单'))
    return
  }
  callback()
}

const validatePermission: FormItemRule['validator'] = (_rule, value, callback) => {
  const permission = normalizeText(String(value ?? ''))
  if (form.type === MENU_TYPE_DIR && permission) {
    callback(new Error('目录无需填写权限标识'))
    return
  }
  if (form.type === MENU_TYPE_BUTTON && !permission) {
    callback(new Error('按钮必须填写权限标识'))
    return
  }
  callback()
}

const validatePath: FormItemRule['validator'] = (_rule, value, callback) => {
  if (form.type === MENU_TYPE_MENU && !normalizeText(String(value ?? ''))) {
    callback(new Error('菜单必须填写路由 path'))
    return
  }
  callback()
}

const validateComponent: FormItemRule['validator'] = (_rule, value, callback) => {
  if (form.type === MENU_TYPE_MENU && !normalizeText(String(value ?? ''))) {
    callback(new Error('菜单必须填写组件路径'))
    return
  }
  callback()
}

const validateComponentName: FormItemRule['validator'] = (_rule, value, callback) => {
  if (form.type === MENU_TYPE_MENU && form.keep_alive && !normalizeText(String(value ?? ''))) {
    callback(new Error('开启缓存时必须填写组件名'))
    return
  }
  callback()
}

const rules = computed<FormRules>(() => ({
  parent_id: [{ validator: validateParentId, trigger: 'change' }],
  name: [{ validator: validateName, trigger: 'blur' }],
  permission: [{ validator: validatePermission, trigger: 'blur' }],
  path: [{ validator: validatePath, trigger: 'blur' }],
  component: [{ validator: validateComponent, trigger: 'blur' }],
  component_name: [{ validator: validateComponentName, trigger: 'blur' }],
}))

function normalizeText(value: string) {
  return value.trim()
}

function collectDescendantIds(nodes: MenuTreeNode[], rootId: number): Set<number> {
  const target = findNode(nodes, rootId)
  const set = new Set<number>()
  if (!target) return set
  const walk = (node: MenuTreeNode) => {
    set.add(node.id)
    node.children?.forEach(walk)
  }
  walk(target)
  return set
}

function findNode(nodes: MenuTreeNode[], id: number): MenuTreeNode | null {
  for (const node of nodes) {
    if (node.id === id) return node
    if (node.children?.length) {
      const found = findNode(node.children, id)
      if (found) return found
    }
  }
  return null
}

function stripParentCandidates(nodes: MenuTreeNode[]): MenuTreeNode[] {
  return nodes
    .filter((node) => node.type !== MENU_TYPE_BUTTON && !forbiddenParentIds.value.has(node.id))
    .map((node) => ({
      ...node,
      children: node.children?.length ? stripParentCandidates(node.children) : undefined,
    }))
}

const parentTreeData = computed<MenuTreeNode[]>(() => {
  const root: MenuTreeNode = {
    id: 0,
    name: '根',
    parent_id: 0,
    type: MENU_TYPE_DIR,
    sort: 0,
    path: '',
    icon: '',
    component: '',
    component_name: '',
    permission: '',
    status: 0,
    visible: true,
    children: stripParentCandidates(treeData.value),
  }
  return [root]
})

function syncFormByType(nextType: number) {
  if (nextType === MENU_TYPE_DIR) {
    form.permission = ''
    form.component = ''
    form.component_name = ''
    form.keep_alive = false
    return
  }
  if (nextType === MENU_TYPE_MENU) {
    return
  }
  form.path = ''
  form.icon = ''
  form.component = ''
  form.component_name = ''
  form.keep_alive = false
  form.always_show = false
}

watch(
  () => form.type,
  (nextType, prevType) => {
    if (nextType === prevType) return
    syncFormByType(nextType)
    formRef.value?.clearValidate(['parent_id', 'permission', 'path', 'component', 'component_name'])
  },
)

watch(
  () => form.keep_alive,
  () => {
    formRef.value?.clearValidate(['component_name'])
  },
)

async function loadTree() {
  loading.value = true
  try {
    const res = (await fetchAdminMenuTree()) as { data?: { list?: MenuTreeNode[] } }
    treeData.value = res?.data?.list ?? []
  } finally {
    loading.value = false
  }
}

function resetForm() {
  editingId.value = 0
  forbiddenParentIds.value = new Set()
  form.parent_id = 0
  form.name = ''
  form.permission = ''
  form.type = MENU_TYPE_DIR
  form.sort = 0
  form.path = ''
  form.icon = ''
  form.component = ''
  form.component_name = ''
  form.status = 0
  form.visible = true
  form.keep_alive = false
  form.always_show = false
  formRef.value?.clearValidate()
}

function openCreate(parentId: number) {
  mode.value = 'create'
  resetForm()
  form.parent_id = parentId
  visible.value = true
}

function openEdit(row: MenuTreeNode) {
  mode.value = 'edit'
  resetForm()
  editingId.value = row.id
  forbiddenParentIds.value = collectDescendantIds(treeData.value, row.id)
  form.parent_id = row.parent_id
  form.name = row.name
  form.permission = row.permission
  form.type = row.type
  form.sort = row.sort
  form.path = row.path
  form.icon = row.icon
  form.component = row.component
  form.component_name = row.component_name
  form.status = row.status
  form.visible = row.visible
  form.keep_alive = !!row.keep_alive
  form.always_show = !!row.always_show
  syncFormByType(form.type)
  visible.value = true
}

function payloadFromForm() {
  const parentId = form.parent_id == null ? 0 : Number(form.parent_id)
  const payload = {
    name: normalizeText(form.name),
    permission: normalizeText(form.permission),
    type: form.type,
    sort: form.sort,
    parent_id: parentId,
    path: normalizeText(form.path),
    icon: normalizeText(form.icon),
    component: normalizeText(form.component),
    component_name: normalizeText(form.component_name),
    status: form.status,
    visible: form.visible,
    keep_alive: form.keep_alive,
    always_show: form.always_show,
  }

  if (payload.type === MENU_TYPE_DIR) {
    payload.permission = ''
    payload.component = ''
    payload.component_name = ''
    payload.keep_alive = false
    return payload
  }

  if (payload.type === MENU_TYPE_BUTTON) {
    payload.path = ''
    payload.icon = ''
    payload.component = ''
    payload.component_name = ''
    payload.keep_alive = false
    payload.always_show = false
  }

  return payload
}

async function submit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    if (mode.value === 'create') {
      await createAdminMenu(payloadFromForm())
      ElMessage.success('已创建')
    } else {
      await updateAdminMenu(editingId.value, payloadFromForm())
      ElMessage.success('已保存')
    }
    visible.value = false
    await loadTree()
  } catch {
    /* 拦截器已提示 */
  } finally {
    saving.value = false
  }
}

function onDelete(row: MenuTreeNode) {
  ElMessageBox.confirm(`删除菜单「${row.name}」？子菜单会一并删除。`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteAdminMenu(row.id)
      ElMessage.success('已删除')
      await loadTree()
    })
    .catch(() => {})
}

onMounted(loadTree)
</script>

<style scoped>
.page {
  padding: 8px 0;
}

.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}
</style>
