<template>
  <div>
    <el-card>
      <template #header>
        <div class="card-header">
          <el-form :inline="true" :model="query">
            <el-form-item label="用户名">
              <el-input v-model="query.username" placeholder="" clearable style="width: 120px" />
            </el-form-item>
            <el-form-item label="请求路径">
              <el-input v-model="query.path" placeholder="" clearable style="width: 180px" />
            </el-form-item>
            <el-form-item label="操作类型">
              <el-select v-model="query.action" placeholder="全部" style="width: 120px" clearable>
                <el-option label="全部" value="" />
                <el-option label="create" value="create" />
                <el-option label="update" value="update" />
                <el-option label="delete" value="delete" />
                <el-option label="query" value="query" />
              </el-select>
            </el-form-item>
            <el-form-item label="日志类型">
              <el-select v-model="query.log_type" placeholder="全部" style="width: 120px" clearable>
                <el-option label="全部" value="" />
                <el-option label="操作" value="operation" />
                <el-option label="API访问" value="api_access" />
              </el-select>
            </el-form-item>
            <el-form-item label="开始时间">
              <el-date-picker
                v-model="query.start_time"
                type="datetime"
                placeholder="选择开始时间"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 180px"
                clearable
              />
            </el-form-item>
            <el-form-item label="结束时间">
              <el-date-picker
                v-model="query.end_time"
                type="datetime"
                placeholder="选择结束时间"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 180px"
                clearable
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadList">搜索</el-button>
              <el-button @click="resetQuery">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </template>
      <el-table :data="list" border @expand-change="onExpandChange">
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="expand-content">
              <div v-if="row.trace_id" class="expand-section">
                <div class="expand-label">TraceId</div>
                <pre class="expand-pre">{{ row.trace_id }}</pre>
              </div>
              <div v-if="row.device_info" class="expand-section">
                <div class="expand-label">设备信息</div>
                <pre class="expand-pre">{{ row.device_info }}</pre>
              </div>
              <div v-if="row.request_data" class="expand-section">
                <div class="expand-label">请求数据</div>
                <pre class="expand-pre">{{ row.request_data }}</pre>
              </div>
              <div v-if="row.response_data" class="expand-section">
                <div class="expand-label">响应数据</div>
                <pre class="expand-pre">{{ row.response_data }}</pre>
              </div>
              <div v-if="changeDetailsMap[row.id]?.length" class="expand-section">
                <div class="expand-label">变更明细</div>
                <el-table :data="changeDetailsMap[row.id]" size="small" border>
                  <el-table-column prop="field_name" label="字段" width="120" />
                  <el-table-column prop="before_value" label="变更前" min-width="120" show-overflow-tooltip />
                  <el-table-column prop="after_value" label="变更后" min-width="120" show-overflow-tooltip />
                </el-table>
              </div>
              <div v-if="!row.request_data && !row.response_data && !row.trace_id && !row.device_info && !changeDetailsMap[row.id]?.length" class="expand-empty">无</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="username" label="用户名" width="100" />
        <el-table-column prop="module" label="模块" width="80" />
        <el-table-column prop="action" label="操作" width="80" />
        <el-table-column prop="log_type" label="日志类型" width="90" />
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column prop="path" label="路径" min-width="180" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP" width="120" />
        <el-table-column prop="duration_ms" label="耗时(ms)" width="90" />
        <el-table-column prop="create_time" label="时间" width="170" />
      </el-table>
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.size"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @current-change="loadList"
        @size-change="loadList"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getAuditLogList, getAuditLogChangeDetails, type AuditLogItem } from '@/api/auditLog'

const list = ref<AuditLogItem[]>([])
const total = ref(0)
const query = reactive({
  username: '',
  path: '',
  action: '',
  log_type: '',
  start_time: '',
  end_time: '',
  page: 1,
  size: 10,
})

const changeDetailsMap = ref<Record<number, { field_name: string; before_value: string; after_value: string }[]>>({})

async function loadList() {
  const params: Record<string, unknown> = {
    page: query.page,
    size: query.size,
  }
  if (query.username) params.username = query.username
  if (query.path) params.path = query.path
  if (query.action) params.action = query.action
  if (query.log_type) params.log_type = query.log_type
  if (query.start_time) params.start_time = query.start_time
  if (query.end_time) params.end_time = query.end_time

  const res = (await getAuditLogList(params)) as { data?: { list: AuditLogItem[]; total: number } }
  list.value = res?.data?.list || []
  total.value = res?.data?.total || 0
}

async function onExpandChange(_row: AuditLogItem, expandedRows: AuditLogItem[]) {
  for (const row of expandedRows) {
    if (!changeDetailsMap.value[row.id] && row.log_type === 'operation') {
      try {
        const res = (await getAuditLogChangeDetails(row.id)) as { data?: { list: { field_name: string; before_value: string; after_value: string }[] } }
        changeDetailsMap.value[row.id] = res?.data?.list || []
      } catch {
        changeDetailsMap.value[row.id] = []
      }
    }
  }
}

function resetQuery() {
  query.username = ''
  query.path = ''
  query.action = ''
  query.log_type = ''
  query.start_time = ''
  query.end_time = ''
  query.page = 1
  changeDetailsMap.value = {}
  loadList()
}

onMounted(loadList)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.expand-content {
  padding: 12px 24px;
}
.expand-section {
  margin-bottom: 12px;
}
.expand-section:last-child {
  margin-bottom: 0;
}
.expand-label {
  font-weight: 500;
  margin-bottom: 6px;
  color: #606266;
}
.expand-pre {
  margin: 0;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 12px;
  max-height: 200px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
.expand-empty {
  color: #909399;
  font-size: 13px;
}
</style>
