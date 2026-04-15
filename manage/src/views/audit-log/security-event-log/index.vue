<template>
  <div>
    <el-card>
      <template #header>
        <div class="card-header">
          <el-form :inline="true" :model="query">
            <el-form-item label="事件类型">
              <el-select v-model="query.event_type" placeholder="全部" style="width: 150px" clearable>
                <el-option label="全部" value="" />
                <el-option label="Token无效" value="token_invalid" />
                <el-option label="权限拒绝" value="permission_denied" />
                <el-option label="暴力破解" value="brute_force" />
                <el-option label="可疑IP" value="suspicious_ip" />
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
              <el-button type="primary" @click="search">搜索</el-button>
              <el-button @click="resetQuery">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </template>
      <el-table :data="list" border>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="event_type" label="事件类型" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.event_type === 'token_invalid'" type="warning">Token无效</el-tag>
            <el-tag v-else-if="row.event_type === 'permission_denied'" type="danger">权限拒绝</el-tag>
            <el-tag v-else-if="row.event_type === 'brute_force'" type="danger">暴力破解</el-tag>
            <el-tag v-else-if="row.event_type === 'suspicious_ip'" type="warning">可疑IP</el-tag>
            <el-tag v-else type="info">{{ row.event_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_id" label="用户ID" width="90" />
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="detail" label="详情" min-width="200" show-overflow-tooltip />
        <el-table-column prop="trace_id" label="TraceId" width="200" show-overflow-tooltip />
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
import { getSecurityEventLogList, type SecurityEventLogItem } from '@/api/securityEventLog'

const list = ref<SecurityEventLogItem[]>([])
const total = ref(0)
const query = reactive({
  event_type: '',
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
  if (query.event_type) params.event_type = query.event_type
  if (query.start_time) params.start_time = query.start_time
  if (query.end_time) params.end_time = query.end_time

  const res = (await getSecurityEventLogList(params)) as { data?: { list: SecurityEventLogItem[]; total: number } }
  list.value = res?.data?.list || []
  total.value = res?.data?.total || 0
}

function search() {
  query.page = 1
  loadList()
}

function resetQuery() {
  query.event_type = ''
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
</style>
