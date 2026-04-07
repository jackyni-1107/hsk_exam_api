<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>任务执行日志</span>
          <el-button @click="$router.push('/task')">返回任务</el-button>
        </div>
      </template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="任务 ID">
          <el-input v-model="taskIdStr" clearable style="width: 120px" placeholder="可选" />
        </el-form-item>
        <el-form-item label="RunId">
          <el-input v-model="query.run_id" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="statusSel" clearable placeholder="全部" style="width: 120px">
            <el-option label="执行中" :value="0" />
            <el-option label="成功" :value="1" />
            <el-option label="失败" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="ID" width="72" />
        <el-table-column prop="task_id" label="任务ID" width="88" />
        <el-table-column prop="run_id" label="RunId" min-width="160" show-overflow-tooltip />
        <el-table-column label="触发" width="88">
          <template #default="{ row }">
            <span v-if="row.trigger_type === 1">定时</span>
            <span v-else-if="row.trigger_type === 2">延迟</span>
            <span v-else-if="row.trigger_type === 3">手动</span>
            <span v-else>{{ row.trigger_type }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="88">
          <template #default="{ row }">
            <el-tag v-if="row.status === 0" type="warning" size="small">执行中</el-tag>
            <el-tag v-else-if="row.status === 1" type="success" size="small">成功</el-tag>
            <el-tag v-else type="danger" size="small">失败</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration_ms" label="耗时ms" width="96" />
        <el-table-column prop="start_time" label="开始" width="170" :formatter="formatUtcForDisplay" />
        <el-table-column prop="end_time" label="结束" width="170" :formatter="formatUtcForDisplay" />
        <el-table-column prop="error_msg" label="错误" min-width="160" show-overflow-tooltip />
        <el-table-column prop="node" label="节点" width="120" show-overflow-tooltip />
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
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { fetchTaskLogs, type TaskLogItem } from '@/api/task'
import { formatUtcForDisplay } from '@/utils/datetime'

const route = useRoute()
const loading = ref(false)
const rows = ref<TaskLogItem[]>([])
const total = ref(0)
const taskIdStr = ref('')

const query = reactive({
  page: 1,
  size: 10,
  run_id: '',
})

const statusSel = ref<number | undefined>(undefined)

const statusParam = computed(() => (statusSel.value === undefined ? null : statusSel.value))

async function loadList() {
  loading.value = true
  try {
    const tid = parseInt(taskIdStr.value, 10)
    const res = (await fetchTaskLogs({
      page: query.page,
      size: query.size,
      task_id: Number.isFinite(tid) && tid > 0 ? tid : undefined,
      run_id: query.run_id || undefined,
      status: statusParam.value,
    })) as { data?: { list?: TaskLogItem[]; total?: number } }
    rows.value = res?.data?.list ?? []
    total.value = res?.data?.total ?? 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.page = 1
  query.size = 10
  query.run_id = ''
  taskIdStr.value = ''
  statusSel.value = undefined
  loadList()
}

onMounted(() => {
  const q = route.query.task_id
  if (q != null) taskIdStr.value = String(q)
  loadList()
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
