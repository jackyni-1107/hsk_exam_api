<template>
  <div>
    <el-card>
      <template #header>
        <div class="card-header">
          <el-form :inline="true" :model="query">
            <el-form-item label="用户名">
              <el-input v-model="query.username" placeholder="" clearable style="width: 120px" />
            </el-form-item>
            <el-form-item label="类型">
              <el-select v-model="query.log_type" placeholder="全部" style="width: 120px" clearable>
                <el-option label="全部" value="" />
                <el-option label="登录成功" value="login_success" />
                <el-option label="登录失败" value="login_fail" />
                <el-option label="登出" value="logout" />
              </el-select>
            </el-form-item>
            <el-form-item label="用户类型">
              <el-select v-model="query.user_type" placeholder="全部" style="width: 120px" clearable>
                <el-option label="全部" :value="-1" />
                <el-option label="后台" :value="1" />
                <el-option label="前台" :value="2" />
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
        <el-table-column prop="log_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.log_type === 'login_success'" type="success">成功</el-tag>
            <el-tag v-else-if="row.log_type === 'login_fail'" type="danger">失败</el-tag>
            <el-tag v-else type="info">{{ row.log_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" width="100" />
        <el-table-column prop="user_type" label="用户类型" width="90">
          <template #default="{ row }">{{ row.user_type === 1 ? '后台' : '前台' }}</template>
        </el-table-column>
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="fail_reason" label="失败原因" min-width="150" show-overflow-tooltip />
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
import { getLoginLogList, type LoginLogItem } from '@/api/loginLog'

const list = ref<LoginLogItem[]>([])
const total = ref(0)
const query = reactive({
  username: '',
  log_type: '',
  user_type: -1,
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
  if (query.username) params.username = query.username
  if (query.log_type) params.log_type = query.log_type
  if (query.user_type >= 1) params.user_type = query.user_type
  if (query.start_time) params.start_time = query.start_time
  if (query.end_time) params.end_time = query.end_time

  const res = (await getLoginLogList(params)) as { data?: { list: LoginLogItem[]; total: number } }
  list.value = res?.data?.list || []
  total.value = res?.data?.total || 0
}

function search() {
  query.page = 1
  loadList()
}

function resetQuery() {
  query.username = ''
  query.log_type = ''
  query.user_type = -1
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
