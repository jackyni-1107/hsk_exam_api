<template>
  <div class="page">
    <el-card shadow="never" class="dash-card">
      <template #header>
        <div class="dash-header">
          <span>{{ t('dashboard.title') }}</span>
          <div class="dash-header-meta" v-if="stats">
            <el-tag v-if="stats.from_cache" size="small" type="info">{{ t('dashboard.dataFromTask') }}</el-tag>
            <el-tag v-else size="small" type="warning">{{ t('dashboard.dataLive') }}</el-tag>
            <span v-if="stats.updated_at" class="dash-updated"
              >{{ t('dashboard.updatedAt') }} {{ stats.updated_at }}</span
            >
          </div>
        </div>
      </template>

      <el-form :inline="true" class="dash-filter" @submit.prevent="loadStats">
        <el-form-item :label="t('dashboard.paperLevel')">
          <el-select
            v-model="mockLevelId"
            clearable
            filterable
            :placeholder="t('common.all')"
            style="width: 220px"
          >
            <el-option
              v-for="lv in mockLevelOptions"
              :key="lv.id"
              :label="mockLevelOptionLabel(lv)"
              :value="lv.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('dashboard.mockPaper')">
          <el-select
            v-model="paperId"
            clearable
            filterable
            :placeholder="t('common.all')"
            style="width: 220px"
          >
            <el-option v-for="p in paperOptions" :key="p.id" :label="`${p.id} · ${p.name}`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('dashboard.batch')">
          <el-select
            v-model="batchId"
            clearable
            filterable
            :placeholder="t('common.all')"
            style="width: 220px"
          >
            <el-option
              v-for="b in batchOptions"
              :key="b.id"
              :label="`${b.id} · ${b.title}`"
              :value="b.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="loadStats">{{ t('dashboard.refresh') }}</el-button>
          <el-button @click="resetFilter">{{ t('dashboard.resetFilter') }}</el-button>
          <el-switch
            v-model="autoRefresh"
            class="dash-poll"
            :active-text="t('dashboard.autoRefresh')"
            style="margin-left: 8px"
          />
        </el-form-item>
      </el-form>

      <el-skeleton v-if="loading && !stats" :rows="6" animated />
      <template v-else-if="stats">
        <el-row :gutter="12" class="dash-kpi">
          <el-col :span="4" v-for="c in kpiCards" :key="c.key" :xs="12" :sm="8" :md="4">
            <el-card
              shadow="hover"
              class="kpi-tile"
              :class="{ 'kpi-tile--link': !!c.drill }"
              @click="c.drill ? goDrill(c.drill) : undefined"
            >
              <div class="kpi-val">{{ c.val }}</div>
              <div class="kpi-label">{{ c.label }}</div>
            </el-card>
          </el-col>
        </el-row>

        <el-row :gutter="12" class="dash-sec">
          <el-col :md="10" :xs="24">
            <el-table :data="trendTable" size="small" border stripe :empty-text="t('dashboard.empty')">
              <el-table-column prop="date" :label="t('dashboard.trendDate')" width="120" />
              <el-table-column prop="count" :label="t('dashboard.trendCount')" align="right" />
            </el-table>
          </el-col>
          <el-col :md="14" :xs="24" class="dash-sec-right">
            <el-table :data="bucketTable" size="small" border stripe :empty-text="t('dashboard.empty')">
              <el-table-column prop="range" :label="t('dashboard.scoreRange')" min-width="140" />
              <el-table-column prop="count" :label="t('dashboard.bucketCount')" width="100" align="right" />
              <el-table-column :label="t('dashboard.ratio')">
                <template #default="{ row }">
                  <el-progress :percentage="row.pct" :stroke-width="10" :show-text="true" />
                </template>
              </el-table-column>
            </el-table>
          </el-col>
        </el-row>
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getAttemptStats, type AttemptStatsData } from '@/api/examAttempt'
import { getMockExaminationPapers, getMockLevelsList, type MockExaminationPaperItem, type MockLevelItem } from '@/api/mockAdmin'
import { getExamBatchList, type ExamBatchListItem } from '@/api/exam'
import { mockLevelOptionLabel } from '@/utils/mockLevel'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const stats = ref<AttemptStatsData | null>(null)
const mockLevelId = ref<number | undefined>(undefined)
const paperId = ref<number | undefined>(undefined)
const batchId = ref<number | undefined>(undefined)
const paperOptions = ref<MockExaminationPaperItem[]>([])
const batchOptions = ref<ExamBatchListItem[]>([])
const mockLevelOptions = ref<MockLevelItem[]>([])

const autoRefresh = ref(false)
let pollTimer: ReturnType<typeof setInterval> | null = null

const trendTable = computed(() => {
  const t7 = stats.value?.trend_7d ?? []
  return t7.map((x) => ({ date: x.date, count: x.count }))
})

const bucketTable = computed(() => {
  const b = stats.value?.score_distribution ?? []
  const sum = b.reduce((s, r) => s + r.count, 0) || 1
  return b.map((r) => {
    const lo = r.bucket_low
    const hi = lo + 10
    const range = `${lo} – ${hi}`
    return {
      range,
      count: r.count,
      pct: Math.round((r.count / sum) * 1000) / 10,
    }
  })
})

const kpiCards = computed(() => {
  if (!stats.value) return []
  const s = stats.value
  return [
    { key: 'n1', label: t('dashboard.s1'), val: s.status_not_started, drill: { status: 1 } },
    { key: 'n2', label: t('dashboard.s2'), val: s.status_in_progress, drill: { status: 2 } },
    { key: 'n3', label: t('dashboard.s3'), val: s.status_submitted, drill: { status: 3 } },
    { key: 'n4', label: t('dashboard.s4'), val: s.status_ended, drill: { status: 4 } },
    { key: 'sp', label: t('dashboard.subjPending'), val: s.subjective_pending, drill: { subjective_pending: 1, status: 4 } },
    { key: 'tn', label: t('dashboard.todayNew'), val: s.today_new, drill: null },
    { key: 'tot', label: t('dashboard.total'), val: s.total, drill: null },
    { key: 'fc', label: t('dashboard.finished'), val: s.finished_count, drill: null },
    { key: 'cp', label: t('dashboard.completion'), val: `${(s.completion_rate * 100).toFixed(1)}%`, drill: null },
    { key: 'ao', label: t('dashboard.avgObj'), val: s.avg_objective.toFixed(2), drill: null },
    { key: 'as', label: t('dashboard.avgSubj'), val: s.avg_subjective.toFixed(2), drill: null },
    { key: 'at', label: t('dashboard.avgTotal'), val: s.avg_total.toFixed(2), drill: null },
  ]
})

async function loadOptions() {
  try {
    const [pr, br, lr] = await Promise.all([
      getMockExaminationPapers({ import_status: 'imported' }) as Promise<{
        data?: { list?: MockExaminationPaperItem[] }
      }>,
      getExamBatchList({ page: 1, size: 200 }) as Promise<{
        data?: { list?: ExamBatchListItem[] }
      }>,
      getMockLevelsList() as Promise<{ data?: { list?: MockLevelItem[] } }>,
    ])
    paperOptions.value = pr?.data?.list ?? []
    batchOptions.value = br?.data?.list ?? []
    mockLevelOptions.value = lr?.data?.list ?? []
  } catch {
    paperOptions.value = []
    batchOptions.value = []
    mockLevelOptions.value = []
  }
}

function statsParams() {
  return {
    mock_level_id: mockLevelId.value && mockLevelId.value > 0 ? mockLevelId.value : undefined,
    examination_paper_id: paperId.value && paperId.value > 0 ? paperId.value : undefined,
    exam_batch_id: batchId.value && batchId.value > 0 ? batchId.value : undefined,
  }
}

async function loadStats() {
  loading.value = true
  try {
    const res = (await getAttemptStats(statsParams())) as { data?: AttemptStatsData }
    stats.value = res?.data ?? null
  } catch {
    stats.value = null
  } finally {
    loading.value = false
  }
}

function resetFilter() {
  mockLevelId.value = undefined
  paperId.value = undefined
  batchId.value = undefined
  loadStats()
}

function baseQuery() {
  const q: Record<string, string> = {}
  if (mockLevelId.value && mockLevelId.value > 0) q.mock_level_id = String(mockLevelId.value)
  if (paperId.value && paperId.value > 0) q.examination_paper_id = String(paperId.value)
  if (batchId.value && batchId.value > 0) q.exam_batch_id = String(batchId.value)
  return q
}

function goDrill(drill: Record<string, number>) {
  const q: Record<string, string> = { ...baseQuery() }
  for (const [k, v] of Object.entries(drill)) {
    q[k] = String(v)
  }
  router.push({ name: 'ExamResult', query: q })
}

watch(
  () => autoRefresh.value,
  (on) => {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
    if (on) {
      pollTimer = setInterval(() => {
        loadStats()
      }, 60_000)
    }
  },
)

onMounted(async () => {
  await loadOptions()
  loadStats()
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<style scoped>
.page {
  padding: 8px 0;
}
.dash-header {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.dash-header-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #606266;
}
.dash-updated {
  white-space: nowrap;
}
.dash-filter {
  margin-bottom: 12px;
}
.dash-poll {
  --el-switch-on-color: #409eff;
}
.dash-kpi {
  margin-bottom: 12px;
}
.kpi-tile {
  text-align: center;
  margin-bottom: 8px;
}
.kpi-tile--link {
  cursor: pointer;
}
.kpi-val {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}
.kpi-label {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
.dash-sec {
  margin-top: 4px;
}
.dash-sec-right {
  margin-top: 8px;
}
@media (min-width: 992px) {
  .dash-sec-right {
    margin-top: 0;
  }
}
</style>
