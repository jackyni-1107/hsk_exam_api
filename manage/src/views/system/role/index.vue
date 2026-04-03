<template>
  <div class="page">
    <el-card shadow="never" class="card">
      <template #header>
        <div class="card-head">
          <span>角色管理</span>
          <el-button type="primary" @click="openCreate">新增角色</el-button>
        </div>
      </template>

      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="角色名">
          <el-input v-model="query.name" clearable placeholder="模糊搜索" style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" style="width: 120px">
            <el-option label="全部" :value="-1" />
            <el-option label="正常" :value="0" />
            <el-option label="停用" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="72" />
        <el-table-column prop="name" label="角色名称" min-width="140" />
        <el-table-column prop="code" label="权限字符" min-width="140" show-overflow-tooltip />
        <el-table-column label="状态" width="88">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">
              {{ row.status === 0 ? '正常' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="菜单数" width="88">
          <template #default="{ row }">
            {{ row.menu_ids?.length ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="onDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pager">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.size"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          background
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogMode === 'create' ? '新增角色' : '编辑角色'"
      width="640px"
      destroy-on-close
      @closed="onDialogClosed"
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="base">
          <el-form ref="formRef" :model="form" :rules="formRules" label-width="96px" class="tab-form">
            <el-form-item label="角色名称" prop="name">
              <el-input v-model="form.name" placeholder="显示名称" />
            </el-form-item>
            <el-form-item label="权限字符" prop="code">
              <el-input v-model="form.code" :disabled="dialogMode === 'edit'" placeholder="唯一编码，创建后不可改" />
            </el-form-item>
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio :label="0">正常</el-radio>
                <el-radio :label="1">停用</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" :max="9999" controls-position="right" />
            </el-form-item>
            <el-form-item label="类型" prop="type">
              <el-input-number v-model="form.type" :min="0" :max="99" controls-position="right" />
            </el-form-item>
            <el-form-item label="备注" prop="remark">
              <el-input v-model="form.remark" type="textarea" :rows="2" placeholder="可选" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="菜单权限" name="menus">
          <div class="tree-wrap">
            <el-tree
              v-if="menuTree.length"
              :key="menuTreeKey"
              ref="menuTreeRef"
              :data="menuTree"
              show-checkbox
              node-key="id"
              :props="{ label: 'name', children: 'children' }"
              :default-checked-keys="checkedMenuIdsForTree"
              default-expand-all
            >
              <template #default="{ data }">
                <span class="tree-node">
                  {{ data.name }}
                  <el-tag v-if="data.type === 1" size="small" type="info" class="tag">目录</el-tag>
                  <el-tag v-else-if="data.type === 2" size="small" type="success" class="tag">菜单</el-tag>
                  <el-tag v-else-if="data.type === 3" size="small" class="tag">按钮</el-tag>
                </span>
              </template>
            </el-tree>
            <el-empty v-else description="正在加载菜单树…" />
          </div>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitForm">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, nextTick } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { ElTree } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchRolePage,
  createRole,
  updateRole,
  assignRoleMenus,
  deleteRole,
  type RoleItem,
} from '@/api/role'
import { fetchAdminMenuTree, type MenuTreeNode } from '@/api/menu'

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<RoleItem[]>([])
const total = ref(0)
const menuTree = ref<MenuTreeNode[]>([])
const menuTreeRef = ref<InstanceType<typeof ElTree>>()
const menuTreeKey = ref(0)

const query = reactive({
  page: 1,
  size: 10,
  name: '',
  status: -1,
})

const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const activeTab = ref('base')
const formRef = ref<FormInstance>()
const checkedMenuIdsForTree = ref<number[]>([])

const form = reactive({
  id: 0,
  name: '',
  code: '',
  status: 0,
  sort: 0,
  type: 0,
  remark: '',
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限字符', trigger: 'blur' }],
}

async function loadMenuTree() {
  if (menuTree.value.length) return
  try {
    const res = (await fetchAdminMenuTree()) as { data?: { list?: MenuTreeNode[] } }
    menuTree.value = res?.data?.list ?? []
  } catch {
    menuTree.value = []
  }
}

function collectMenuIds(): number[] {
  const t = menuTreeRef.value
  if (!t) return [...checkedMenuIdsForTree.value]
  const checked = t.getCheckedKeys(false) as number[]
  const half = t.getHalfCheckedKeys() as number[]
  return [...new Set([...checked, ...half])]
}

async function loadList() {
  loading.value = true
  try {
    const res = (await fetchRolePage({
      page: query.page,
      size: query.size,
      name: query.name || undefined,
      status: query.status,
    })) as { data?: { list?: RoleItem[]; total?: number } }
    tableData.value = res?.data?.list ?? []
    total.value = res?.data?.total ?? 0
  } catch {
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.page = 1
  query.size = 10
  query.name = ''
  query.status = -1
  loadList()
}

function resetForm() {
  form.id = 0
  form.name = ''
  form.code = ''
  form.status = 0
  form.sort = 0
  form.type = 0
  form.remark = ''
  checkedMenuIdsForTree.value = []
  activeTab.value = 'base'
  formRef.value?.clearValidate()
}

function onDialogClosed() {
  resetForm()
}

async function openCreate() {
  dialogMode.value = 'create'
  resetForm()
  await loadMenuTree()
  menuTreeKey.value += 1
  checkedMenuIdsForTree.value = []
  dialogVisible.value = true
  await nextTick()
}

async function openEdit(row: RoleItem) {
  dialogMode.value = 'edit'
  resetForm()
  await loadMenuTree()
  form.id = row.id
  form.name = row.name
  form.code = row.code
  form.status = row.status
  form.sort = row.sort ?? 0
  form.type = row.type ?? 0
  form.remark = row.remark ?? ''
  checkedMenuIdsForTree.value = [...(row.menu_ids ?? [])]
  menuTreeKey.value += 1
  dialogVisible.value = true
  await nextTick()
}

async function submitForm() {
  if (!formRef.value) return
  await formRef.value.validate(async (ok) => {
    if (!ok) return
    submitLoading.value = true
    try {
      const menuIds = collectMenuIds()
      if (dialogMode.value === 'create') {
        const res = (await createRole({
          name: form.name,
          code: form.code,
          status: form.status,
          sort: form.sort,
          type: form.type,
          remark: form.remark,
        })) as { data?: { id?: number } }
        const newId = res?.data?.id
        if (newId != null) {
          await assignRoleMenus(newId, menuIds)
        }
        ElMessage.success('创建成功')
      } else {
        await updateRole(form.id, {
          name: form.name,
          status: form.status,
          sort: form.sort,
          type: form.type,
          remark: form.remark,
        })
        await assignRoleMenus(form.id, menuIds)
        ElMessage.success('已保存')
      }
      dialogVisible.value = false
      await loadList()
    } catch {
      /* 错误已提示 */
    } finally {
      submitLoading.value = false
    }
  })
}

function onDelete(row: RoleItem) {
  ElMessageBox.confirm(`确定删除角色「${row.name}」？`, '删除确认', { type: 'warning' })
    .then(async () => {
      await deleteRole(row.id)
      ElMessage.success('已删除')
      await loadList()
    })
    .catch(() => {})
}

onMounted(() => {
  loadList()
})
</script>

<style scoped>
.page {
  padding: 8px 0;
}

.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.filter {
  margin-bottom: 16px;
}

.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.tab-form {
  padding-top: 8px;
}

.tree-wrap {
  max-height: 420px;
  overflow: auto;
  padding: 8px 0;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.tree-node {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.tree-node .tag {
  transform: scale(0.85);
}
</style>
