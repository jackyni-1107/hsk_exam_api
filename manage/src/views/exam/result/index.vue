<template>
  <div class="page">
    <el-card shadow="never">
      <template #header><span>考试结果</span></template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="学员账号">
          <el-input v-model="query.username" clearable placeholder="模糊" style="width: 160px" />
        </el-form-item>
        <el-form-item label="试卷级别">
          <el-input v-model="query.level" clearable placeholder="如 hsk1" style="width: 120px" />
        </el-form-item>
        <el-form-item label="Mock 卷">
          <el-select
            v-model="paperSel"
            clearable
            filterable
            placeholder="全部"
            style="width: 220px"
          >
            <el-option
              v-for="p in paperOptions"
              :key="p.id"
              :label="`${p.id} · ${p.name}`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="批次ID">
          <el-input-number
            v-model="query.exam_batch_id"
            :min="0"
            :controls="false"
            placeholder="0=不限"
            style="width: 120px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable placeholder="全部" style="width: 140px">
            <el-option label="全部" :value="0" />
            <el-option label="未开始" :value="1" />
            <el-option label="进行中" :value="2" />
            <el-option label="已交卷" :value="3" />
            <el-option label="已结束" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="会话ID" width="96" />
        <el-table-column prop="username" label="学员" width="120" show-overflow-tooltip />
        <el-table-column prop="nickname" label="昵称" width="100" show-overflow-tooltip />
        <el-table-column prop="paper_title" label="试卷" min-width="160" show-overflow-tooltip />
        <el-table-column prop="paper_level" label="级别" width="80" />
        <el-table-column prop="exam_batch_id" label="批次" width="88" align="right" />
        <el-table-column prop="mock_level_id" label="等级" width="88" align="right" />
        <el-table-column label="状态" width="88">
          <template #default="{ row }">
            {{ statusText(row.status) }}
          </template>
        </el-table-column>
        <el-table-column label="客观分" width="88" align="right">
          <template #default="{ row }">{{ row.objective_score?.toFixed?.(2) ?? row.objective_score }}</template>
        </el-table-column>
        <el-table-column label="主观分" width="88" align="right">
          <template #default="{ row }">{{ row.subjective_score?.toFixed?.(2) ?? row.subjective_score }}</template>
        </el-table-column>
        <el-table-column label="总分" width="88" align="right">
          <template #default="{ row }">{{ row.total_score?.toFixed?.(2) ?? row.total_score }}</template>
        </el-table-column>
        <el-table-column prop="submitted_at" label="交卷时间" width="172" show-overflow-tooltip />
        <el-table-column prop="ended_at" label="结束时间" width="172" show-overflow-tooltip />
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

    <el-drawer v-model="drawer" title="考试结果详情" size="min(720px, 92vw)" destroy-on-close>
      <template v-if="detail">
        <el-descriptions :column="2" border size="small" class="blk">
          <el-descriptions-item label="会话ID">{{ detail.attempt.id }}</el-descriptions-item>
          <el-descriptions-item label="状态">{{ statusText(detail.attempt.status) }}</el-descriptions-item>
          <el-descriptions-item label="学员">
            {{ detail.user.username }}（{{ detail.user.nickname || '—' }}）
          </el-descriptions-item>
          <el-descriptions-item label="试卷">{{ detail.paper.title }}</el-descriptions-item>
          <el-descriptions-item label="客观分">{{ detail.attempt.objective_score }}</el-descriptions-item>
          <el-descriptions-item label="主观分">{{ detail.attempt.subjective_score }}</el-descriptions-item>
          <el-descriptions-item label="总分">{{ detail.attempt.total_score }}</el-descriptions-item>
          <el-descriptions-item label="含主观题">{{ detail.attempt.has_subjective ? '是' : '否' }}</el-descriptions-item>
          <el-descriptions-item label="开考">{{ detail.attempt.started_at || '—' }}</el-descriptions-item>
          <el-descriptions-item label="截止">{{ detail.attempt.deadline_at || '—' }}</el-descriptions-item>
          <el-descriptions-item label="交卷">{{ detail.attempt.submitted_at || '—' }}</el-descriptions-item>
          <el-descriptions-item label="结束">{{ detail.attempt.ended_at || '—' }}</el-descriptions-item>
        </el-descriptions>

        <h4 class="sub">答题明细</h4>
        <el-table :data="detail.answers" border size="small" max-height="360">
          <el-table-column prop="question_no" label="题号" width="72" />
          <el-table-column label="类型" width="80">
            <template #default="{ row }">
              <span v-if="row.is_example">例题</span>
              <span v-else-if="row.is_subjective">主观</span>
              <span v-else>客观</span>
            </template>
          </el-table-column>
          <el-table-column label="题干" min-width="140" show-overflow-tooltip>
            <template #default="{ row }">{{ row.stem_text || '—' }}</template>
          </el-table-column>
          <el-table-column label="答案" min-width="120" show-overflow-tooltip>
            <template #default="{ row }">{{ formatAnswerJson(row.answer_json) }}</template>
          </el-table-column>
          <el-table-column label="客观对错" width="88">
            <template #default="{ row }">
              <template v-if="row.objective_correct === true">对</template>
              <template v-else-if="row.objective_correct === false">错</template>
              <template v-else>—</template>
            </template>
          </el-table-column>
          <el-table-column label="满分" width="72" align="right">
            <template #default="{ row }">{{ row.score ?? '—' }}</template>
          </el-table-column>
          <el-table-column label="得分" width="72" align="right">
            <template #default="{ row }">
              <template v-if="row.is_subjective && !row.is_example">
                {{ row.awarded_score != null ? row.awarded_score : '未评' }}
              </template>
              <template v-else>—</template>
            </template>
          </el-table-column>
        </el-table>

        <template v-if="canGradeSubjective">
          <h4 class="sub">主观题评分</h4>
          <p class="hint">每题得分不超过该题满分；保存后更新主观分与总分。</p>
          <el-table :data="subjectiveRows" border size="small">
            <el-table-column prop="question_no" label="题号" width="72" />
            <el-table-column label="题干" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">{{ row.stem_text || '—' }}</template>
            </el-table-column>
            <el-table-column label="满分" width="80" align="right">
              <template #default="{ row }">{{ row.score }}</template>
            </el-table-column>
            <el-table-column label="得分" width="140">
              <template #default="{ row }">
                <el-input-number
                  v-model="subjectiveScores[row.question_id]"
                  :min="0"
                  :max="row.score"
                  :precision="2"
                  :step="0.5"
                  size="small"
                  controls-position="right"
                  style="width: 120px"
                />
              </template>
            </el-table-column>
          </el-table>
          <div class="grade-actions">
            <el-button type="primary" :loading="savingScores" @click="saveSubjective">保存主观分</el-button>
          </div>
        </template>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getAttemptList,
  getAttemptDetail,
  saveAttemptSubjectiveScores,
  type AttemptListItem,
  type AttemptDetail,
} from '@/api/examAttempt'
import { getMockExaminationPapers, type MockExaminationPaperItem } from '@/api/mockAdmin'

const loading = ref(false)
const rows = ref<AttemptListItem[]>([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 10,
  username: '',
  level: '',
  status: 0 as number,
  examination_paper_id: 0 as number,
  exam_batch_id: 0 as number,
})

const paperSel = ref<number | undefined>(undefined)
const paperOptions = ref<MockExaminationPaperItem[]>([])

const drawer = ref(false)
const detail = ref<AttemptDetail | null>(null)
const detailId = ref(0)
const subjectiveScores = ref<Record<number, number>>({})
const savingScores = ref(false)

const subjectiveRows = computed(() => {
  if (!detail.value) return []
  return detail.value.answers.filter((a) => a.is_subjective && !a.is_example)
})

const canGradeSubjective = computed(() => {
  const a = detail.value?.attempt
  return !!(a && a.status === 4 && a.has_subjective === 1)
})

function statusText(s: number) {
  switch (s) {
    case 1:
      return '未开始'
    case 2:
      return '进行中'
    case 3:
      return '已交卷'
    case 4:
      return '已结束'
    default:
      return String(s)
  }
}

function formatAnswerJson(raw: string) {
  if (!raw) return '—'
  try {
    const o = JSON.parse(raw) as Record<string, unknown>
    if (typeof o.option_id === 'number') return `选项 ID: ${o.option_id}`
    if (typeof o.text === 'string') return o.text.length > 80 ? `${o.text.slice(0, 80)}…` : o.text
    return JSON.stringify(o)
  } catch {
    return raw.length > 100 ? `${raw.slice(0, 100)}…` : raw
  }
}

async function loadPapers() {
  try {
    const res = (await getMockExaminationPapers({ import_status: 'imported' })) as {
      data?: { list?: MockExaminationPaperItem[] }
    }
    paperOptions.value = res?.data?.list ?? []
  } catch {
    paperOptions.value = []
  }
}

async function loadList() {
  loading.value = true
  try {
    query.examination_paper_id = paperSel.value ?? 0
    const res = (await getAttemptList({
      page: query.page,
      size: query.size,
      username: query.username || undefined,
      level: query.level || undefined,
      examination_paper_id: query.examination_paper_id || undefined,
      exam_batch_id: query.exam_batch_id > 0 ? query.exam_batch_id : undefined,
      status: query.status || undefined,
    })) as { data?: { list?: AttemptListItem[]; total?: number } }
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
  query.level = ''
  query.status = 0
  query.exam_batch_id = 0
  paperSel.value = undefined
  loadList()
}

async function openDetail(row: AttemptListItem) {
  detailId.value = row.id
  drawer.value = true
  try {
    const res = (await getAttemptDetail(row.id)) as { data?: AttemptDetail }
    detail.value = res?.data ?? null
    const m: Record<number, number> = {}
    for (const a of detail.value?.answers ?? []) {
      if (a.is_subjective && !a.is_example) {
        m[a.question_id] = a.awarded_score != null ? Number(a.awarded_score) : 0
      }
    }
    subjectiveScores.value = m
  } catch {
    detail.value = null
  }
}

async function saveSubjective() {
  if (!detail.value) return
  const items = subjectiveRows.value.map((r) => ({
    question_id: r.question_id,
    score: subjectiveScores.value[r.question_id] ?? 0,
  }))
  if (items.length === 0) {
    ElMessage.warning('没有可保存的主观题')
    return
  }
  savingScores.value = true
  try {
    const res = (await saveAttemptSubjectiveScores(detailId.value, items)) as {
      data?: { subjective_score: number; total_score: number }
    }
    ElMessage.success('保存成功')
    if (res?.data && detail.value) {
      detail.value.attempt.subjective_score = res.data.subjective_score
      detail.value.attempt.total_score = res.data.total_score
    }
    loadList()
  } finally {
    savingScores.value = false
  }
}

onMounted(() => {
  loadPapers()
  loadList()
})
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
.blk {
  margin-bottom: 16px;
}
.sub {
  margin: 14px 0 8px;
  font-size: 14px;
  font-weight: 600;
}
.hint {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin: 0 0 10px;
}
.grade-actions {
  margin-top: 12px;
}
</style>
