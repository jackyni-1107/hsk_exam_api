<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>任务管理</span>
          <div>
            <el-button @click="$router.push('/task/log')">执行日志</el-button>
            <el-button v-permission="'task:create'" type="primary" @click="openCreate">新建任务</el-button>
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
      <el-row :gutter="12" class="runtime-stats">
        <el-col v-for="card in runtimeCards" :key="card.key" :xs="24" :sm="12" :md="6">
          <el-card shadow="hover" class="runtime-stat-card">
            <div class="runtime-stat-head">
              <span class="runtime-stat-label">{{ card.label }}</span>
              <el-tag v-if="card.tagText" size="small" effect="light" :type="card.tagType">{{ card.tagText }}</el-tag>
            </div>
            <div class="runtime-stat-value" :class="{ 'runtime-stat-value--time': card.key === 'oldest_due' }">
              {{ card.value }}
            </div>
            <div class="runtime-stat-sub">{{ card.sub }}</div>
          </el-card>
        </el-col>
      </el-row>
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
        <el-table-column prop="create_time" label="创建时间" width="170" :formatter="formatUtcForDisplay" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button v-permission="'task:run'" link type="success" :disabled="row.status !== 0" @click="onRun(row)">执行</el-button>
            <el-button link type="primary" @click="goLog(row)">日志</el-button>
            <el-button v-permission="'task:update'" link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button v-permission="'task:delete'" link type="danger" @click="onDelete(row)">删除</el-button>
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
          <el-input-number v-model="form.concurrency" :min="0" />
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
  fetchTaskRuntimeStats,
  createTask,
  updateTask,
  deleteTask,
  runTask,
  type TaskItem,
  type TaskRuntimeStats,
} from '@/api/task'
import { formatUtcForDisplay } from '@/utils/datetime'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const rows = ref<TaskItem[]>([])
const total = ref(0)
const runtimeStats = ref<TaskRuntimeStats | null>(null)
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

type RuntimeCard = {
  key: string
  label: string
  value: string | number
  sub: string
  tagText?: string
  tagType?: 'success' | 'info' | 'warning' | 'danger'
}

const runtimeCards = computed<RuntimeCard[]>(() => {
  const stats = runtimeStats.value
  return [
    {
      key: 'delay_queue',
      label: '延迟队列',
      value: stats?.delay_queue_size ?? '-',
      sub: '持久化的延迟任务与重试任务。',
    },
    {
      key: 'due_now',
      label: '到期待执行',
      value: stats?.delay_due_count ?? '-',
      sub: '已到执行时间，等待扫描器领取。',
    },
    {
      key: 'scanner',
      label: '扫描器',
      value: stats ? (stats.delay_scanner_active ? '运行中' : '待命') : '-',
      sub: stats
        ? (stats.delay_scanner_active
            ? `锁租约剩余 ${formatDuration(stats.delay_scanner_ttl_millis)}`
            : '正在等待当前扫描器租约释放。')
        : '正在加载运行态信息。',
      tagText: stats ? (stats.delay_scanner_active ? '活跃' : '空闲') : undefined,
      tagType: stats?.delay_scanner_active ? 'success' : 'info',
    },
    {
      key: 'oldest_due',
      label: '最早到期',
      value: stats?.delay_oldest_due_at
        ? formatUtcForDisplay(null, null, stats.delay_oldest_due_at)
        : '--',
      sub: stats?.delay_oldest_due_at ? '当前队列里最早仍未处理的任务。' : '当前没有等待中的延迟任务。',
    },
  ]
})

function formatDuration(ms: number) {
  if (!Number.isFinite(ms) || ms <= 0) return '0 ms'
  if (ms < 1000) return `${ms} ms`
  if (ms < 60_000) return `${(ms / 1000).toFixed(ms >= 10_000 ? 0 : 1)} s`
  return `${(ms / 60_000).toFixed(ms >= 600_000 ? 0 : 1)} min`
}

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

async function loadRuntimeStats() {
  const res = (await fetchTaskRuntimeStats()) as { data?: TaskRuntimeStats }
  runtimeStats.value = res?.data ?? null
}

function resetQuery() {
  query.page = 1
  query.size = 10
  query.name = ''
  query.code = ''
  query.type = undefined
  query.status = undefined
  void Promise.all([loadList(), loadRuntimeStats()])
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
  form.concurrency = typeof row.concurrency === 'number' ? row.concurrency : 1
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
      await Promise.all([loadList(), loadRuntimeStats()])
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
      await Promise.all([loadList(), loadRuntimeStats()])
    })
    .catch(() => {})
}

async function onRun(row: TaskItem) {
  try {
    const res = (await runTask(row.id)) as { data?: { run_id?: string } }
    if (row.type === 2) {
      ElMessage.success(`已加入延迟队列，run_id: ${res?.data?.run_id ?? '-'}`)
      await loadRuntimeStats()
      return
    }
    ElMessage.success(`已触发执行，run_id: ${res?.data?.run_id ?? '-'}`)
    await loadRuntimeStats()
  } catch {
    /* */
  }
}

function goLog(row: TaskItem) {
  router.push({ path: '/task/log', query: { task_id: String(row.id) } })
}

onMounted(async () => {
  await Promise.all([loadList(), loadRuntimeStats()])
})
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
.runtime-stats {
  margin-bottom: 12px;
}
.runtime-stat-card {
  height: 100%;
  border-color: #e4e7ed;
}
.runtime-stat-card :deep(.el-card__body) {
  height: 100%;
  min-height: 132px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 12px;
}
.runtime-stat-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}
.runtime-stat-label {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
}
.runtime-stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  line-height: 1.25;
  letter-spacing: -0.02em;
}
.runtime-stat-value--time {
  font-size: 18px;
  line-height: 1.45;
}
.runtime-stat-sub {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  min-height: 36px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
