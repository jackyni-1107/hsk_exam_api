<template>
  <div class="page">
    <el-card shadow="never">
      <template #header><span>异常日志</span></template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="TraceId">
          <el-input v-model="query.trace_id" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="路径">
          <el-input v-model="query.path" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="开始">
          <el-date-picker
            v-model="startDt"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="可选"
          />
        </el-form-item>
        <el-form-item label="结束">
          <el-date-picker
            v-model="endDt"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="可选"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="ID" width="72" />
        <el-table-column prop="method" label="方法" width="72" />
        <el-table-column prop="path" label="路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="error_msg" label="错误信息" min-width="200" show-overflow-tooltip />
        <el-table-column prop="user_id" label="用户ID" width="88" />
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="trace_id" label="Trace" min-width="120" show-overflow-tooltip />
        <el-table-column prop="create_time" label="时间" width="180" />
        <el-table-column label="操作" width="88" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
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

    <el-drawer v-model="drawer" title="异常详情" size="55%">
      <template v-if="current">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="路径">{{ current.method }} {{ current.path }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{ current.user_id }}</el-descriptions-item>
          <el-descriptions-item label="IP">{{ current.ip }}</el-descriptions-item>
          <el-descriptions-item label="Trace">{{ current.trace_id }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ current.create_time }}</el-descriptions-item>
        </el-descriptions>
        <h4 class="sub">错误信息</h4>
        <pre class="pre err">{{ current.error_msg || '—' }}</pre>
        <h4 class="sub">堆栈</h4>
        <pre class="pre stack">{{ current.stack || '—' }}</pre>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { fetchExceptionLogList, type ExceptionLogItem } from '@/api/exceptionLog'

const loading = ref(false)
const rows = ref<ExceptionLogItem[]>([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 10,
  trace_id: '',
  path: '',
})
const startDt = ref<string | null>(null)
const endDt = ref<string | null>(null)

const drawer = ref(false)
const current = ref<ExceptionLogItem | null>(null)

async function loadList() {
  loading.value = true
  try {
    const res = (await fetchExceptionLogList({
      page: query.page,
      size: query.size,
      trace_id: query.trace_id || undefined,
      path: query.path || undefined,
      start_time: startDt.value || undefined,
      end_time: endDt.value || undefined,
    })) as { data?: { list?: ExceptionLogItem[]; total?: number } }
    rows.value = res?.data?.list ?? []
    total.value = res?.data?.total ?? 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.page = 1
  query.size = 10
  query.trace_id = ''
  query.path = ''
  startDt.value = null
  endDt.value = null
  loadList()
}

function openDetail(row: ExceptionLogItem) {
  current.value = row
  drawer.value = true
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
.sub {
  margin: 14px 0 8px;
  font-size: 14px;
  font-weight: 600;
}
.pre {
  background: #f1f5f9;
  padding: 10px;
  border-radius: 8px;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 280px;
  overflow: auto;
}
.pre.err {
  background: #fef2f2;
  color: #991b1b;
}
.pre.stack {
  background: #0f172a;
  color: #e2e8f0;
  max-height: 400px;
}
</style>
