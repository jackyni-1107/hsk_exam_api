<template>
  <div class="page">
    <el-card shadow="never" class="card">
      <template #header>
        <div class="card-head">
          <span>用户管理</span>
          <el-button type="primary" @click="openCreate">新增用户</el-button>
        </div>
      </template>

      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="用户名">
          <el-input v-model="query.username" clearable placeholder="模糊搜索" style="width: 200px" />
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
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="email" label="邮箱" min-width="160" show-overflow-tooltip />
        <el-table-column prop="mobile" label="手机" width="120" />
        <el-table-column label="角色" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">
            {{ roleNames(row.role_ids) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="88">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">
              {{ row.status === 0 ? '正常' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="170" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="warning" @click="onKick(row)">下线</el-button>
            <el-button
              v-if="row.id !== SUPER_ADMIN_ID"
              link
              type="danger"
              @click="onDelete(row)"
            >
              删除
            </el-button>
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
      :title="dialogMode === 'create' ? '新增用户' : '编辑用户'"
      width="520px"
      destroy-on-close
      @closed="resetForm"
    >
      <el-form ref="formRef" :model="form" :rules="dynamicRules" label-width="88px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="dialogMode === 'edit'" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="dialogMode === 'create' ? '密码' : '新密码'" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            autocomplete="new-password"
            :placeholder="dialogMode === 'edit' ? '不修改请留空' : ''"
          />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="手机" prop="mobile">
          <el-input v-model="form.mobile" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="0">正常</el-radio>
            <el-radio :label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role_ids" multiple collapse-tags collapse-tags-tooltip placeholder="选择角色" style="width: 100%">
            <el-option v-for="r in roleOptions" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchUserList,
  createUser,
  updateUser,
  assignUserRoles,
  deleteUser,
  kickUserSessions,
  type AdminUserItem,
} from '@/api/user'
import { fetchRoleList, type RoleItem } from '@/api/role'

const SUPER_ADMIN_ID = 1

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<AdminUserItem[]>([])
const total = ref(0)
const roleOptions = ref<RoleItem[]>([])

const query = reactive({
  page: 1,
  size: 10,
  username: '',
  status: -1,
})

const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()

const form = reactive({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
  mobile: '',
  status: 0,
  role_ids: [] as number[],
})

const dynamicRules = computed<FormRules>(() => ({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password:
    dialogMode.value === 'create'
      ? [{ required: true, message: '请输入密码', trigger: 'blur' }]
      : [],
}))

async function loadRoles() {
  try {
    const res = (await fetchRoleList()) as { data?: { list?: RoleItem[] } }
    roleOptions.value = res?.data?.list ?? []
  } catch {
    roleOptions.value = []
  }
}

function roleNames(ids: number[]) {
  if (!ids?.length) return '—'
  const map = new Map(roleOptions.value.map((r) => [r.id, r.name]))
  return ids.map((id) => map.get(id) || id).join('、')
}

async function loadList() {
  loading.value = true
  try {
    const res = (await fetchUserList({
      page: query.page,
      size: query.size,
      username: query.username || undefined,
      status: query.status < 0 ? undefined : query.status,
    })) as { data?: { list?: AdminUserItem[]; total?: number } }
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
  query.username = ''
  query.status = -1
  loadList()
}

function resetForm() {
  form.id = 0
  form.username = ''
  form.password = ''
  form.nickname = ''
  form.email = ''
  form.mobile = ''
  form.status = 0
  form.role_ids = []
  formRef.value?.clearValidate()
}

function openCreate() {
  dialogMode.value = 'create'
  resetForm()
  dialogVisible.value = true
}

function openEdit(row: AdminUserItem) {
  dialogMode.value = 'edit'
  resetForm()
  form.id = row.id
  form.username = row.username
  form.nickname = row.nickname
  form.email = row.email
  form.mobile = row.mobile
  form.status = row.status
  form.role_ids = [...(row.role_ids ?? [])]
  dialogVisible.value = true
}

async function submitForm() {
  if (!formRef.value) return
  await formRef.value.validate(async (ok) => {
    if (!ok) return
    submitLoading.value = true
    try {
      if (dialogMode.value === 'create') {
        await createUser({
          username: form.username,
          password: form.password,
          nickname: form.nickname,
          email: form.email,
          mobile: form.mobile,
          status: form.status,
          role_ids: form.role_ids,
        })
        ElMessage.success('创建成功')
      } else {
        await updateUser(form.id, {
          password: form.password || undefined,
          nickname: form.nickname,
          email: form.email,
          mobile: form.mobile,
          status: form.status,
        })
        await assignUserRoles(form.id, form.role_ids)
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

function onDelete(row: AdminUserItem) {
  ElMessageBox.confirm(`确定删除用户「${row.username}」？`, '删除确认', {
    type: 'warning',
  })
    .then(async () => {
      await deleteUser(row.id)
      ElMessage.success('已删除')
      await loadList()
    })
    .catch(() => {})
}

function onKick(row: AdminUserItem) {
  ElMessageBox.confirm(`将用户「${row.username}」的所有会话强制下线？`, '下线确认', {
    type: 'warning',
  })
    .then(async () => {
      await kickUserSessions(row.id)
      ElMessage.success('已执行下线')
    })
    .catch(() => {})
}

onMounted(async () => {
  await loadRoles()
  await loadList()
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
</style>
