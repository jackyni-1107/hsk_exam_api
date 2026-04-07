<template>
  <div class="page">
    <el-card shadow="never">
      <template #header><span>考试批次</span></template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="Mock 卷">
          <el-select
            v-model="query.mock_examination_paper_id"
            clearable
            filterable
            placeholder="全部"
            style="width: 260px"
          >
            <el-option
              v-for="p in paperOptions"
              :key="p.id"
              :label="`${p.id} · ${p.name}`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
          <el-button type="success" @click="openCreate">新建批次</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="批次ID" width="88" />
        <el-table-column prop="mock_examination_paper_id" label="Mock卷ID" width="100" />
        <el-table-column prop="title" label="批次名称" min-width="140" show-overflow-tooltip />
        <el-table-column prop="exam_start_at" label="开始时间" width="172" show-overflow-tooltip :formatter="formatUtcForDisplay" />
        <el-table-column prop="exam_end_at" label="结束时间" width="172" show-overflow-tooltip :formatter="formatUtcForDisplay" />
        <el-table-column label="等级IDs" min-width="120" show-overflow-tooltip>
          <template #default="{ row }">{{ (row.mock_level_ids || []).join(', ') || '—' }}</template>
        </el-table-column>
        <el-table-column prop="member_count" label="学员数" width="80" align="right" />
        <el-table-column prop="create_time" label="创建时间" width="172" show-overflow-tooltip :formatter="formatUtcForDisplay" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="primary" @click="openMembers(row)">成员</el-button>
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

    <el-dialog v-model="formVisible" :title="formMode === 'create' ? '新建批次' : '编辑批次'" width="520px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="120px">
        <el-form-item label="Mock 卷" prop="mock_examination_paper_id">
          <el-select
            v-model="form.mock_examination_paper_id"
            filterable
            placeholder="选择已导入的卷"
            style="width: 100%"
            :disabled="formMode === 'edit'"
          >
            <el-option
              v-for="p in importedPapers"
              :key="p.id"
              :label="`${p.id} · ${p.name}`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="批次名称" prop="title">
          <el-input v-model="form.title" clearable placeholder="可选" />
        </el-form-item>
        <el-form-item label="考试开始" prop="exam_start_at">
          <el-date-picker
            v-model="form.exam_start_at"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="选择时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="考试结束" prop="exam_end_at">
          <el-date-picker
            v-model="form.exam_end_at"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="选择时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="考试等级" prop="mock_level_ids">
          <el-select v-model="form.mock_level_ids" multiple filterable placeholder="mock_levels.id 多选" style="width: 100%">
            <el-option
              v-for="lv in levelOptions"
              :key="lv.id"
              :label="`${lv.id} · ${lv.level_name || lv.app_level_name}`"
              :value="lv.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" :loading="formSaving" @click="submitForm">保存</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="memberDrawer" title="批次成员" size="min(640px, 92vw)" destroy-on-close @opened="loadMemberList">
      <template v-if="currentBatch">
        <p class="hint">批次 #{{ currentBatch.id }} · {{ currentBatch.title || '（无标题）' }}</p>
        <el-form :inline="true" class="mb">
          <el-form-item label="导入等级" required>
            <el-select
              v-model="importMockLevelId"
              placeholder="mock_levels.id"
              style="width: 220px"
              filterable
            >
              <el-option
                v-for="id in currentBatch.mock_level_ids || []"
                :key="id"
                :label="`${id}`"
                :value="id"
              />
            </el-select>
          </el-form-item>
        </el-form>
        <el-space wrap class="mb">
          <el-input
            v-model="importIdsText"
            type="textarea"
            :rows="2"
            placeholder="输入要导入的会员 ID，逗号或换行分隔"
            style="width: 320px"
          />
          <el-button type="primary" :loading="importing" @click="doImport">导入</el-button>
        </el-space>
        <el-form :inline="true" class="filter" @submit.prevent="loadMemberList">
          <el-form-item label="账号">
            <el-input v-model="memberSearch.username" clearable placeholder="模糊" style="width: 140px" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="searchMembersToPick">搜索会员</el-button>
          </el-form-item>
        </el-form>
        <el-table v-if="pickMembers.length" :data="pickMembers" size="small" max-height="200" border class="mb">
          <el-table-column prop="id" label="ID" width="72" />
          <el-table-column prop="username" label="账号" />
          <el-table-column prop="nickname" label="昵称" />
          <el-table-column label="" width="100">
            <template #default="{ row }">
              <el-button link type="primary" @click="addPickId(row.id)">加入导入</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-table v-loading="memberLoading" :data="memberRows" border size="small">
          <el-table-column prop="member_id" label="会员ID" width="88" />
          <el-table-column prop="mock_level_id" label="等级ID" width="88" />
          <el-table-column prop="username" label="账号" />
          <el-table-column prop="nickname" label="昵称" />
          <el-table-column prop="import_time" label="导入时间" width="172" :formatter="formatUtcForDisplay" />
          <el-table-column label="" width="88">
            <template #default="{ row }">
              <el-button link type="danger" @click="removeOne(row)">移除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="pager">
          <el-pagination
            v-model:current-page="memberQuery.page"
            v-model:page-size="memberQuery.size"
            :total="memberTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            small
            background
            @size-change="loadMemberList"
            @current-change="loadMemberList"
          />
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  createExamBatch,
  deleteExamBatch,
  getExamBatchList,
  getExamBatchMemberList,
  importExamBatchMembers,
  removeExamBatchMembers,
  updateExamBatch,
  type ExamBatchListItem,
  type ExamBatchMemberItem,
} from '@/api/exam'
import { getMockExaminationPapers, getMockLevelsList, type MockExaminationPaperItem } from '@/api/mockAdmin'
import { fetchMemberList, type MemberItem } from '@/api/member'
import { formatUtcForDisplay, formatUtcText } from '@/utils/datetime'

const loading = ref(false)
const rows = ref<ExamBatchListItem[]>([])
const total = ref(0)
const paperOptions = ref<MockExaminationPaperItem[]>([])
const importedPapers = ref<MockExaminationPaperItem[]>([])
const levelOptions = ref<{ id: number; level_name: string; app_level_name: string }[]>([])

const query = reactive({
  mock_examination_paper_id: undefined as number | undefined,
  page: 1,
  size: 10,
})

async function loadPapers() {
  const res = await getMockExaminationPapers()
  paperOptions.value = res.data?.list ?? []
  importedPapers.value = (res.data?.list ?? []).filter((p) => p.imported)
}

async function loadLevels() {
  const res = await getMockLevelsList()
  levelOptions.value = res.data?.list ?? []
}

async function loadList() {
  loading.value = true
  try {
    const res = await getExamBatchList({
      mock_examination_paper_id: query.mock_examination_paper_id || 0,
      page: query.page,
      size: query.size,
    })
    rows.value = res.data?.list ?? []
    total.value = res.data?.total ?? 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.mock_examination_paper_id = undefined
  query.page = 1
  loadList()
}

const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formSaving = ref(false)
const formRef = ref<FormInstance>()
const editingId = ref(0)
const form = reactive({
  mock_examination_paper_id: undefined as number | undefined,
  title: '',
  exam_start_at: '',
  exam_end_at: '',
  mock_level_ids: [] as number[],
})

const formRules: FormRules = {
  mock_examination_paper_id: [{ required: true, message: '请选择 Mock 卷', trigger: 'change' }],
  exam_start_at: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  exam_end_at: [{ required: true, message: '请选择结束时间', trigger: 'change' }],
  mock_level_ids: [
    { type: 'array', required: true, message: '请选择等级', trigger: 'change' },
    {
      type: 'array',
      min: 1,
      message: '至少选择一个等级',
      trigger: 'change',
    },
  ],
}

function openCreate() {
  formMode.value = 'create'
  editingId.value = 0
  form.mock_examination_paper_id = importedPapers.value[0]?.id
  form.title = ''
  form.exam_start_at = ''
  form.exam_end_at = ''
  form.mock_level_ids = []
  formVisible.value = true
}

function openEdit(row: ExamBatchListItem) {
  formMode.value = 'edit'
  editingId.value = row.id
  form.mock_examination_paper_id = row.mock_examination_paper_id
  form.title = row.title
  form.exam_start_at = formatUtcText(row.exam_start_at)
  form.exam_end_at = formatUtcText(row.exam_end_at)
  form.mock_level_ids = [...(row.mock_level_ids || [])]
  formVisible.value = true
}

async function submitForm() {
  await formRef.value?.validate().catch(() => Promise.reject())
  if (!form.mock_examination_paper_id) return
  formSaving.value = true
  try {
    const start = form.exam_start_at.trim()
    const end = form.exam_end_at.trim()
    if (formMode.value === 'create') {
      await createExamBatch({
        mock_examination_paper_id: form.mock_examination_paper_id,
        title: form.title,
        exam_start_at: start,
        exam_end_at: end,
        mock_level_ids: form.mock_level_ids,
      })
      ElMessage.success('已创建')
    } else {
      await updateExamBatch(editingId.value, {
        title: form.title,
        exam_start_at: start,
        exam_end_at: end,
        mock_level_ids: form.mock_level_ids,
      })
      ElMessage.success('已保存')
    }
    formVisible.value = false
    loadList()
  } finally {
    formSaving.value = false
  }
}

async function onDelete(row: ExamBatchListItem) {
  await ElMessageBox.confirm(`确定删除批次 #${row.id}？`, '确认', { type: 'warning' })
  await deleteExamBatch(row.id)
  ElMessage.success('已删除')
  loadList()
}

const memberDrawer = ref(false)
const currentBatch = ref<ExamBatchListItem | null>(null)
const memberLoading = ref(false)
const memberRows = ref<ExamBatchMemberItem[]>([])
const memberTotal = ref(0)
const memberQuery = reactive({ page: 1, size: 10 })
const importIdsText = ref('')
const importMockLevelId = ref(0)
const importing = ref(false)
const pickMembers = ref<MemberItem[]>([])

function openMembers(row: ExamBatchListItem) {
  currentBatch.value = row
  memberQuery.page = 1
  importIdsText.value = ''
  importMockLevelId.value = row.mock_level_ids?.[0] ?? 0
  pickMembers.value = []
  memberDrawer.value = true
}

async function loadMemberList() {
  if (!currentBatch.value) return
  memberLoading.value = true
  try {
    const res = await getExamBatchMemberList(currentBatch.value.id, {
      page: memberQuery.page,
      size: memberQuery.size,
    })
    memberRows.value = res.data?.list ?? []
    memberTotal.value = res.data?.total ?? 0
  } finally {
    memberLoading.value = false
  }
}

function parseIdList(text: string): number[] {
  const parts = text.split(/[\s,;，；]+/).map((s) => s.trim()).filter(Boolean)
  const ids: number[] = []
  for (const p of parts) {
    const n = Number(p)
    if (Number.isFinite(n) && n > 0) ids.push(n)
  }
  return [...new Set(ids)]
}

async function doImport() {
  if (!currentBatch.value) return
  if (!importMockLevelId.value) {
    ElMessage.warning('请选择导入等级')
    return
  }
  const ids = parseIdList(importIdsText.value)
  if (!ids.length) {
    ElMessage.warning('请输入有效的会员 ID')
    return
  }
  importing.value = true
  try {
    const res = await importExamBatchMembers(currentBatch.value.id, {
      mock_level_id: importMockLevelId.value,
      member_ids: ids,
    })
    ElMessage.success(`已导入 ${res.data?.inserted ?? 0} 人`)
    importIdsText.value = ''
    loadMemberList()
    loadList()
  } finally {
    importing.value = false
  }
}

async function removeOne(row: ExamBatchMemberItem) {
  if (!currentBatch.value) return
  await ElMessageBox.confirm('从批次中移除该学员？', '确认', { type: 'warning' })
  await removeExamBatchMembers(currentBatch.value.id, {
    mock_level_id: row.mock_level_id,
    member_ids: [row.member_id],
  })
  ElMessage.success('已移除')
  loadMemberList()
  loadList()
}

const memberSearch = reactive({ username: '' })

async function searchMembersToPick() {
  const res = await fetchMemberList({
    page: 1,
    size: 20,
    username: memberSearch.username || undefined,
    status: 0,
  })
  pickMembers.value = res.data?.list ?? []
}

function addPickId(id: number) {
  const cur = parseIdList(importIdsText.value)
  if (!cur.includes(id)) cur.push(id)
  importIdsText.value = cur.join(', ')
}

onMounted(async () => {
  await loadPapers()
  await loadLevels()
  loadList()
})
</script>

<style scoped>
.page {
  padding: 16px;
}
.filter {
  margin-bottom: 12px;
}
.pager {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}
.hint {
  margin: 0 0 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
.mb {
  margin-bottom: 12px;
}
</style>
