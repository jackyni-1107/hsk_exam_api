<template>
  <div>
    <el-card>
      <template #header>
        <div class="card-header">
          <el-form :inline="true" :model="query">
            <el-form-item label="TraceId">
              <el-input v-model="query.trace_id" placeholder="" clearable style="width: 200px" />
            </el-form-item>
            <el-form-item label="路径">
              <el-input v-model="query.path" placeholder="" clearable style="width: 180px" />
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
              <el-button type="primary" @click="search">搜索</el-button>
              <el-button @click="resetQuery">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </template>
      <el-table :data="list" border>
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="expand-content">
              <div v-if="row.stack" class="expand-section">
                <div class="expand-label">堆栈</div>
                <pre class="expand-pre">{{ row.stack }}</pre>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="trace_id" label="TraceId" width="280" show-overflow-tooltip />
        <el-table-column prop="path" label="路径" min-width="180" show-overflow-tooltip />
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column prop="error_msg" label="错误信息" min-width="200" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP" width="120" />
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
import { getExceptionLogList, type ExceptionLogItem } from '@/api/exceptionLog'

const list = ref<ExceptionLogItem[]>([])
const total = ref(0)
const query = reactive({
  trace_id: '',
  path: '',
  start_time: '',
  end_time: '',
  page: 1,
  size: 10,
})

async function loadList() {
  const params: Record<string, unknown> = {
    page: query.page,
    size: query.size,
  }
  if (query.trace_id) params.trace_id = query.trace_id
  if (query.path) params.path = query.path
  if (query.start_time) params.start_time = query.start_time
  if (query.end_time) params.end_time = query.end_time

  const res = (await getExceptionLogList(params)) as { data?: { list: ExceptionLogItem[]; total: number } }
  list.value = res?.data?.list || []
  total.value = res?.data?.total || 0
}

function search() {
  query.page = 1
  loadList()
}

function resetQuery() {
  query.trace_id = ''
  query.path = ''
  query.start_time = ''
  query.end_time = ''
  query.page = 1
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
</style>
