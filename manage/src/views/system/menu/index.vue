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
            <span v-if="row.type === 1">目录</span>
            <span v-else-if="row.type === 2">菜单</span>
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
            <el-button link type="primary" @click="openCreate(row.id)">子项</el-button>
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
            placeholder="不选则为根（实际选「根」）"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :label="1">目录</el-radio>
            <el-radio :label="2">菜单</el-radio>
            <el-radio :label="3">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="权限标识" prop="permission">
          <el-input v-model="form.permission" placeholder="如 user:list" />
        </el-form-item>
        <el-form-item label="路由 path" prop="path">
          <el-input v-model="form.path" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="form.icon" placeholder="Element 图标名" />
        </el-form-item>
        <el-form-item label="组件" prop="component">
          <el-input v-model="form.component" />
        </el-form-item>
        <el-form-item label="组件名" prop="component_name">
          <el-input v-model="form.component_name" />
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
        <el-form-item label="缓存">
          <el-switch v-model="form.keep_alive" />
        </el-form-item>
        <el-form-item label="总显示">
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
import { ref, reactive, computed, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchAdminMenuTree,
  createAdminMenu,
  updateAdminMenu,
  deleteAdminMenu,
  type MenuTreeNode,
} from '@/api/menu'

const loading = ref(false)
const saving = ref(false)
const treeData = ref<MenuTreeNode[]>([])
const visible = ref(false)
const mode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const editingId = ref(0)
const forbiddenParentIds = ref<Set<number>>(new Set())

const form = reactive({
  parent_id: 0 as number | null,
  name: '',
  permission: '',
  type: 1,
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

const rules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

const flatRows = computed(() => treeData.value)

function collectDescendantIds(nodes: MenuTreeNode[], rootId: number): Set<number> {
  const target = findNode(nodes, rootId)
  const set = new Set<number>()
  if (!target) return set
  const walk = (n: MenuTreeNode) => {
    set.add(n.id)
    n.children?.forEach(walk)
  }
  walk(target)
  return set
}

function findNode(nodes: MenuTreeNode[], id: number): MenuTreeNode | null {
  for (const n of nodes) {
    if (n.id === id) return n
    if (n.children?.length) {
      const f = findNode(n.children, id)
      if (f) return f
    }
  }
  return null
}

const parentTreeData = computed(() => {
  const root: MenuTreeNode = { id: 0, name: '根', parent_id: 0, type: 1, sort: 0, path: '', icon: '', component: '', component_name: '', permission: '', status: 0, visible: true, children: [...treeData.value] }
  const strip = (nodes: MenuTreeNode[]): MenuTreeNode[] =>
    nodes
      .filter((n) => !forbiddenParentIds.value.has(n.id))
      .map((n) => ({
        ...n,
        children: n.children?.length ? strip(n.children) : undefined,
      }))
  return [{ ...root, children: strip(treeData.value) }]
})

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
  form.type = 1
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
  visible.value = true
}

function payloadFromForm() {
  const pid = form.parent_id == null ? 0 : Number(form.parent_id)
  return {
    name: form.name,
    permission: form.permission,
    type: form.type,
    sort: form.sort,
    parent_id: pid,
    path: form.path,
    icon: form.icon,
    component: form.component,
    component_name: form.component_name,
    status: form.status,
    visible: form.visible,
    keep_alive: form.keep_alive,
    always_show: form.always_show,
  }
}

async function submit() {
  await formRef.value?.validate(async (ok) => {
    if (!ok) return
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
  })
}

function onDelete(row: MenuTreeNode) {
  ElMessageBox.confirm(`删除菜单「${row.name}」？子菜单需先删除。`, '确认', { type: 'warning' })
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
