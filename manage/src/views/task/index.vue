<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>任务管理</span>
          <div>
            <el-button @click="$router.push('/task/log')">执行日志</el-button>
            <el-button type="primary" @click="openCreate">新建任务</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="名称">
          <el-input v-model="query.name" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="编码">
          <el-input v-model="query.code" clearable style="width: 140px" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="query.type" clearable placeholder="全部" style="width: 120px">
            <el-option label="Cron" :value="1" />
            <el-option label="延迟" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable placeholder="全部" style="width: 110px">
            <el-option label="启用" :value="0" />
            <el-option label="停用" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="ID" width="72" />
        <el-table-column prop="name" label="名称" min-width="140" />
        <el-table-column prop="code" label="编码" width="130" show-overflow-tooltip />
        <el-table-column label="类型" width="88">
          <template #default="{ row }">{{ row.type === 1 ? 'Cron' : '延迟' }}</template>
        </el-table-column>
        <el-table-column prop="cron_expr" label="Cron" min-width="120" show-overflow-tooltip />
        <el-table-column prop="handler" label="处理器" min-width="140" show-overflow-tooltip />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'" size="small">{{ row.status === 0 ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="170" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button link type="success" :disabled="row.status !== 0" @click="onRun(row)">执行</el-button>
            <el-button link type="primary" @click="goLog(row)">日志</el-button>
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

    <el-dialog v-model="dlg" :title="mode === 'create' ? '新建任务' : '编辑任务'" width="640px" destroy-on-close @closed="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="编码" prop="code">
          <el-input v-model="form.code" :disabled="mode === 'edit'" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :label="1">Cron</el-radio>
            <el-radio :label="2">延迟</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.type === 1" label="Cron 表达式" prop="cron_expr">
          <el-input v-model="form.cron_expr" placeholder="如 0 0 * * * *" />
        </el-form-item>
        <el-form-item v-if="form.type === 2" label="延迟(秒)" prop="delay_seconds">
          <el-input-number v-model="form.delay_seconds" :min="1" />
        </el-form-item>
        <el-form-item label="处理器" prop="handler">
          <el-input v-model="form.handler" placeholder="注册的处理函数名" />
        </el-form-item>
        <el-form-item label="参数 JSON">
          <el-input v-model="form.params" type="textarea" :rows="3" placeholder="{}" />
        </el-form-item>
        <el-form-item label="重试次数">
          <el-input-number v-model="form.retry_times" :min="0" />
        </el-form-item>
        <el-form-item label="重试间隔(秒)">
          <el-input-number v-model="form.retry_interval" :min="0" />
        </el-form-item>
        <el-form-item label="并发数">
          <el-input-number v-model="form.concurrency" :min="1" />
        </el-form-item>
        <el-form-item label="失败告警">
          <el-radio-group v-model="form.alert_on_fail">
            <el-radio :label="0">否</el-radio>
            <el-radio :label="1">是</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="告警接收">
          <el-input v-model="form.alert_receivers" placeholder="可选 JSON 或邮箱列表" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="0">启用</el-radio>
            <el-radio :label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
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
import { reactive, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  fetchTaskList,
  createTask,
  updateTask,
  deleteTask,
  runTask,
  type TaskItem,
} from '@/api/task'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const rows = ref<TaskItem[]>([])
const total = ref(0)
const dlg = ref(false)
const mode = ref<'create' | 'edit'>('create')
const formRef = ref<FormInstance>()
const editId = ref(0)

const query = reactive({
  page: 1,
  size: 10,
  name: '',
  code: '',
  type: undefined as number | undefined,
  status: undefined as number | undefined,
})

const form = reactive({
  name: '',
  code: '',
  type: 1,
  cron_expr: '',
  delay_seconds: 60,
  handler: '',
  params: '',
  retry_times: 0,
  retry_interval: 60,
  concurrency: 1,
  alert_on_fail: 0,
  alert_receivers: '',
  status: 0,
  remark: '',
})

const rules = computed<FormRules>(() => ({
  name: [{ required: true, message: '必填', trigger: 'blur' }],
  code: [{ required: true, message: '必填', trigger: 'blur' }],
  type: [{ required: true, message: '必选', trigger: 'change' }],
  handler: [{ required: true, message: '必填', trigger: 'blur' }],
  cron_expr:
    form.type === 1 ? [{ required: true, message: 'Cron 必填', trigger: 'blur' }] : [],
  delay_seconds:
    form.type === 2 ? [{ required: true, message: '延迟须大于 0', trigger: 'change' }] : [],
}))

function payload() {
  return {
    name: form.name,
    code: form.code,
    type: form.type,
    cron_expr: form.type === 1 ? form.cron_expr : '',
    delay_seconds: form.type === 2 ? form.delay_seconds : 0,
    handler: form.handler,
    params: form.params || '',
    retry_times: form.retry_times,
    retry_interval: form.retry_interval,
    concurrency: form.concurrency,
    alert_on_fail: form.alert_on_fail,
    alert_receivers: form.alert_receivers || '',
    status: form.status,
    remark: form.remark || '',
  }
}

async function loadList() {
  loading.value = true
  try {
    const res = (await fetchTaskList({
      page: query.page,
      size: query.size,
      name: query.name || undefined,
      code: query.code || undefined,
      type: query.type,
      status: query.status,
    })) as { data?: { list?: TaskItem[]; total?: number } }
    rows.value = res?.data?.list ?? []
    total.value = res?.data?.total ?? 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.page = 1
  query.size = 10
  query.name = ''
  query.code = ''
  query.type = undefined
  query.status = undefined
  loadList()
}

function resetForm() {
  editId.value = 0
  Object.assign(form, {
    name: '',
    code: '',
    type: 1,
    cron_expr: '',
    delay_seconds: 60,
    handler: '',
    params: '',
    retry_times: 0,
    retry_interval: 60,
    concurrency: 1,
    alert_on_fail: 0,
    alert_receivers: '',
    status: 0,
    remark: '',
  })
  formRef.value?.clearValidate()
}

function openCreate() {
  mode.value = 'create'
  resetForm()
  dlg.value = true
}

function openEdit(row: TaskItem) {
  mode.value = 'edit'
  resetForm()
  editId.value = row.id
  form.name = row.name
  form.code = row.code
  form.type = row.type
  form.cron_expr = row.cron_expr || ''
  form.delay_seconds = row.delay_seconds || 60
  form.handler = row.handler
  form.params = row.params || ''
  form.retry_times = row.retry_times
  form.retry_interval = row.retry_interval
  form.concurrency = row.concurrency || 1
  form.alert_on_fail = row.alert_on_fail
  form.alert_receivers = row.alert_receivers || ''
  form.status = row.status
  form.remark = row.remark || ''
  dlg.value = true
}

async function submit() {
  await formRef.value?.validate(async (ok) => {
    if (!ok) return
    saving.value = true
    try {
      if (mode.value === 'create') {
        await createTask(payload())
        ElMessage.success('已创建')
      } else {
        await updateTask(editId.value, payload())
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

function onDelete(row: TaskItem) {
  ElMessageBox.confirm(`删除任务「${row.name}」？`, '确认', { type: 'warning' })
    .then(async () => {
      await deleteTask(row.id)
      ElMessage.success('已删除')
      await loadList()
    })
    .catch(() => {})
}

async function onRun(row: TaskItem) {
  try {
    const res = (await runTask(row.id)) as { data?: { run_id?: string } }
    ElMessage.success(`已触发执行，run_id: ${res?.data?.run_id ?? '-'}`)
  } catch {
    /* */
  }
}

function goLog(row: TaskItem) {
  router.push({ path: '/task/log', query: { task_id: String(row.id) } })
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
  gap: 12px;
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
