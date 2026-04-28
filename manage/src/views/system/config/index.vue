<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>配置中心</span>
          <el-button v-permission="'config:create'" type="primary" @click="openCreate">新增配置</el-button>
        </div>
      </template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="分组">
          <el-input v-model="query.group" clearable style="width: 140px" />
        </el-form-item>
        <el-form-item label="键">
          <el-input v-model="query.key" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="ID" width="72" />
        <el-table-column prop="config_key" label="键" min-width="160" show-overflow-tooltip />
        <el-table-column prop="config_value" label="值" min-width="200" show-overflow-tooltip />
        <el-table-column prop="config_type" label="类型" width="100" />
        <el-table-column prop="group_name" label="分组" width="120" />
        <el-table-column prop="remark" label="备注" min-width="140" show-overflow-tooltip />
        <el-table-column prop="create_time" label="创建时间" width="170" :formatter="formatUtcForDisplay" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button v-permission="'config:update'" link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button v-permission="'config:delete'" link type="danger" @click="onDelete(row)">删除</el-button>
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

    <el-dialog v-model="dlg" :title="mode === 'create' ? '新增配置' : '编辑配置'" width="520px" destroy-on-close @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="96px">
        <el-form-item label="键" prop="config_key">
          <el-input v-model="form.config_key" :disabled="mode === 'edit'" />
        </el-form-item>
        <el-form-item label="值" prop="config_value">
          <el-input v-model="form.config_value" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="类型">
          <el-input v-model="form.config_type" placeholder="默认 string" />
        </el-form-item>
        <el-form-item label="分组">
          <el-input v-model="form.group_name" placeholder="默认 default" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchConfigList,
  createConfig,
  updateConfig,
  deleteConfig,
  type ConfigItemRow,
} from '@/api/sysConfig'
import { formatUtcForDisplay } from '@/utils/datetime'

const loading = ref(false)
const saving = ref(false)
const rows = ref<ConfigItemRow[]>([])
const total = ref(0)
const dlg = ref(false)
const mode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const editId = ref(0)

const query = reactive({ page: 1, size: 10, group: '', key: '' })
const form = reactive({
  config_key: '',
  config_value: '',
  config_type: 'string',
  group_name: 'default',
  remark: '',
})

const rules: FormRules = {
  config_key: [{ required: true, message: '必填', trigger: 'blur' }],
  config_value: [{ required: true, message: '必填', trigger: 'blur' }],
}

async function loadList() {
  loading.value = true
  try {
    const res = (await fetchConfigList({
      page: query.page,
      size: query.size,
      group: query.group || undefined,
      key: query.key || undefined,
    })) as { data?: { list?: ConfigItemRow[]; total?: number } }
    rows.value = res?.data?.list ?? []
    total.value = res?.data?.total ?? 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.page = 1
  query.size = 10
  query.group = ''
  query.key = ''
  loadList()
}

function resetForm() {
  editId.value = 0
  form.config_key = ''
  form.config_value = ''
  form.config_type = 'string'
  form.group_name = 'default'
  form.remark = ''
  formRef.value?.clearValidate()
}

function openCreate() {
  mode.value = 'create'
  resetForm()
  dlg.value = true
}

function openEdit(row: ConfigItemRow) {
  mode.value = 'edit'
  resetForm()
  editId.value = row.id
  form.config_key = row.config_key
  form.config_value = row.config_value
  form.config_type = row.config_type || 'string'
  form.group_name = row.group_name || 'default'
  form.remark = row.remark
  dlg.value = true
}

async function submit() {
  await formRef.value?.validate(async (ok) => {
    if (!ok) return
    saving.value = true
    try {
      if (mode.value === 'create') {
        await createConfig({
          config_key: form.config_key,
          config_value: form.config_value,
          config_type: form.config_type || undefined,
          group_name: form.group_name || undefined,
          remark: form.remark || undefined,
        })
        ElMessage.success('已创建')
      } else {
        await updateConfig(editId.value, {
          config_value: form.config_value,
          remark: form.remark,
        })
        ElMessage.success('已保存')
      }
      dlg.value = false
      await loadList()
    } catch {
      /* */
    } finally {
      saving.value = false
    }
  })
}

function onDelete(row: ConfigItemRow) {
  ElMessageBox.confirm(`删除配置「${row.config_key}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteConfig(row.id)
      ElMessage.success('已删除')
      await loadList()
    })
    .catch(() => {})
}

onMounted(loadList)
</script>

<style scoped>
.page {
  padding: 8px 0;
}
.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.filter {
  margin-bottom: 12px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
