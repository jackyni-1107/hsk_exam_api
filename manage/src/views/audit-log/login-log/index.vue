<template>
  <div class="page">
    <el-card shadow="never">
      <template #header><span>登录日志</span></template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="用户名">
          <el-input v-model="query.username" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item label="类型">
          <el-input v-model="query.log_type" clearable style="width: 120px" placeholder="如 success" />
        </el-form-item>
        <el-form-item label="用户类型">
          <el-select v-model="userTypeSel" clearable placeholder="全部" style="width: 130px">
            <el-option label="管理端" :value="1" />
            <el-option label="客户端" :value="2" />
          </el-select>
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
        <el-table-column prop="log_type" label="类型" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column label="用户类型" width="96">
          <template #default="{ row }">
            {{ row.user_type === 1 ? '管理端' : row.user_type === 2 ? '客户端' : row.user_type }}
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="fail_reason" label="失败原因" min-width="140" show-overflow-tooltip />
        <el-table-column prop="trace_id" label="Trace" min-width="120" show-overflow-tooltip />
        <el-table-column prop="create_time" label="时间" width="180" :formatter="formatUtcForDisplay" />
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

    <el-drawer v-model="drawer" title="登录日志详情" size="46%">
      <template v-if="current">
        <el-descriptions :column="1" border size="small">
          <el-descriptions-item label="用户名">{{ current.username }} ({{ current.user_id }})</el-descriptions-item>
          <el-descriptions-item label="类型">{{ current.log_type }}</el-descriptions-item>
          <el-descriptions-item label="用户类型">
            {{ current.user_type === 1 ? '管理端' : current.user_type === 2 ? '客户端' : current.user_type }}
          </el-descriptions-item>
          <el-descriptions-item label="IP">{{ current.ip }}</el-descriptions-item>
          <el-descriptions-item label="Trace">{{ current.trace_id }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ formatUtcText(current.create_time) }}</el-descriptions-item>
          <el-descriptions-item label="失败原因">{{ current.fail_reason || '—' }}</el-descriptions-item>
        </el-descriptions>
        <h4 class="sub">User-Agent</h4>
        <pre class="pre">{{ current.user_agent || '—' }}</pre>
        <h4 class="sub">设备信息</h4>
        <pre class="pre">{{ current.device_info || '—' }}</pre>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { fetchLoginLogList, type LoginLogItem } from '@/api/loginLog'
import { formatUtcForDisplay, formatUtcText } from '@/utils/datetime'

const loading = ref(false)
const rows = ref<LoginLogItem[]>([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 10,
  username: '',
  log_type: '',
})
const userTypeSel = ref<number | undefined>(undefined)
const startDt = ref<string | null>(null)
const endDt = ref<string | null>(null)

const drawer = ref(false)
const current = ref<LoginLogItem | null>(null)

async function loadList() {
  loading.value = true
  try {
    const res = (await fetchLoginLogList({
      page: query.page,
      size: query.size,
      username: query.username || undefined,
      log_type: query.log_type || undefined,
      user_type: userTypeSel.value,
      start_time: startDt.value || undefined,
      end_time: endDt.value || undefined,
    })) as { data?: { list?: LoginLogItem[]; total?: number } }
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
  query.log_type = ''
  userTypeSel.value = undefined
  startDt.value = null
  endDt.value = null
  loadList()
}

function openDetail(row: LoginLogItem) {
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
  max-height: 200px;
  overflow: auto;
}
</style>
