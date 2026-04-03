<template>
  <div class="page">
    <el-card shadow="never">
      <template #header><span>操作审计</span></template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="用户">
          <el-input v-model="query.username" clearable style="width: 140px" />
        </el-form-item>
        <el-form-item label="路径">
          <el-input v-model="query.path" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="动作">
          <el-input v-model="query.action" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item label="类型">
          <el-input v-model="query.log_type" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item label="TraceId">
          <el-input v-model="query.trace_id" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="开始">
          <el-date-picker v-model="startDt" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="可选" />
        </el-form-item>
        <el-form-item label="结束">
          <el-date-picker v-model="endDt" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="可选" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="ID" width="72" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="method" label="方法" width="72" />
        <el-table-column prop="path" label="路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="action" label="动作" width="100" />
        <el-table-column prop="log_type" label="类型" width="100" />
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="duration_ms" label="耗时ms" width="88" />
        <el-table-column prop="trace_id" label="Trace" min-width="120" show-overflow-tooltip />
        <el-table-column prop="create_time" label="时间" width="180" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
            <el-button link type="primary" @click="openChanges(row)">变更</el-button>
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

    <el-drawer v-model="detailOpen" title="审计详情" size="50%">
      <template v-if="current">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="用户">{{ current.username }} ({{ current.user_id }})</el-descriptions-item>
          <el-descriptions-item label="请求">{{ current.method }} {{ current.path }}</el-descriptions-item>
          <el-descriptions-item label="动作">{{ current.action }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ current.log_type }}</el-descriptions-item>
          <el-descriptions-item label="模块">{{ current.module }}</el-descriptions-item>
          <el-descriptions-item label="IP">{{ current.ip }}</el-descriptions-item>
          <el-descriptions-item label="UA">{{ current.user_agent }}</el-descriptions-item>
          <el-descriptions-item label="Trace">{{ current.trace_id }}</el-descriptions-item>
          <el-descriptions-item label="耗时">{{ current.duration_ms }} ms</el-descriptions-item>
          <el-descriptions-item label="时间">{{ current.create_time }}</el-descriptions-item>
        </el-descriptions>
        <h4 class="block-title">请求体</h4>
        <pre class="json-box">{{ pretty(current.request_data) }}</pre>
        <h4 class="block-title">响应体</h4>
        <pre class="json-box">{{ pretty(current.response_data) }}</pre>
      </template>
    </el-drawer>

    <el-dialog v-model="changeOpen" title="变更明细" width="720px" destroy-on-close>
      <el-table v-loading="changeLoading" :data="changeRows" border max-height="420">
        <el-table-column prop="table_name" label="表" width="160" />
        <el-table-column prop="record_id" label="记录ID" width="88" />
        <el-table-column prop="field_name" label="字段" width="140" />
        <el-table-column prop="before_value" label="变更前" min-width="160" show-overflow-tooltip />
        <el-table-column prop="after_value" label="变更后" min-width="160" show-overflow-tooltip />
        <el-table-column prop="create_time" label="时间" width="180" />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { fetchAuditLogList, fetchAuditChangeDetails, type AuditLogItem } from '@/api/auditLog'

const loading = ref(false)
const rows = ref<AuditLogItem[]>([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 10,
  username: '',
  path: '',
  action: '',
  log_type: '',
  trace_id: '',
})

const startDt = ref<string | null>(null)
const endDt = ref<string | null>(null)

const detailOpen = ref(false)
const current = ref<AuditLogItem | null>(null)

const changeOpen = ref(false)
const changeLoading = ref(false)
const changeRows = ref<unknown[]>([])

function pretty(s: string) {
  if (!s) return '—'
  try {
    return JSON.stringify(JSON.parse(s), null, 2)
  } catch {
    return s
  }
}

async function loadList() {
  loading.value = true
  try {
    const res = (await fetchAuditLogList({
      page: query.page,
      size: query.size,
      username: query.username || undefined,
      path: query.path || undefined,
      action: query.action || undefined,
      log_type: query.log_type || undefined,
      trace_id: query.trace_id || undefined,
      start_time: startDt.value || undefined,
      end_time: endDt.value || undefined,
    })) as { data?: { list?: AuditLogItem[]; total?: number } }
    rows.value = res?.data?.list ?? []
    total.value = res?.data?.total ?? 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.page = 1
  query.size = 10
  query.username = ''
  query.path = ''
  query.action = ''
  query.log_type = ''
  query.trace_id = ''
  startDt.value = null
  endDt.value = null
  loadList()
}

function openDetail(row: AuditLogItem) {
  current.value = row
  detailOpen.value = true
}

async function openChanges(row: AuditLogItem) {
  changeOpen.value = true
  changeLoading.value = true
  changeRows.value = []
  try {
    const res = (await fetchAuditChangeDetails(row.id)) as { data?: { list?: unknown[] } }
    changeRows.value = res?.data?.list ?? []
  } finally {
    changeLoading.value = false
  }
}

onMounted(loadList)
</script>

<style scoped>
.page {
  padding: 8px 0;
}
.filter {
  margin-bottom: 12px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
.block-title {
  margin: 16px 0 8px;
  font-size: 14px;
  font-weight: 600;
}
.json-box {
  background: #0f172a;
  color: #e2e8f0;
  padding: 12px;
  border-radius: 8px;
  font-size: 12px;
  overflow: auto;
  max-height: 240px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
