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
          <template #default="{ row }">
            <span class="list-total-score">{{ row.total_score?.toFixed?.(2) ?? row.total_score }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="submitted_at" label="交卷时间" width="172" show-overflow-tooltip :formatter="formatUtcForDisplay" />
        <el-table-column prop="submitted_at" label="开始时间" width="172" show-overflow-tooltip :formatter="formatUtcForDisplay" />
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

    <el-drawer
      v-model="drawer"
      title="考试结果详情"
      size="min(960px, 96vw)"
      destroy-on-close
      class="result-drawer-wrap"
    >
      <div v-if="detail" class="result-detail">
        <div class="detail-hero">
          <div class="hero-inner">
            <div class="hero-ring-wrap">
              <el-progress
                type="circle"
                :percentage="ringPercentage"
                :width="148"
                :stroke-width="11"
                :color="ringStrokeColor"
                :show-text="false"
              />
              <div class="hero-score-overlay">
                <span class="hero-score" :style="{ color: ringStrokeColor }">{{ formatScore(displayScore) }}</span>
                <span class="hero-max">/ {{ paperMaxScore > 0 ? formatScore(paperMaxScore) : '—' }}</span>
              </div>
            </div>
            <div class="hero-meta-row">
              <el-tag v-if="gradeBand" :type="gradeBand.tagType" round effect="light" class="grade-tag">
                <component :is="gradeBand.icon" class="grade-ic" />
                {{ gradeBand.label }}
              </el-tag>
              <span class="hero-duration-inline">考试耗时 {{ durationLabel }}</span>
            </div>
          </div>
        </div>

        <div class="detail-mid-grid">
          <el-card shadow="never" class="mid-card">
            <div class="mid-title">学员</div>
            <p class="mid-strong">{{ detail.user.username }}</p>
            <p class="mid-muted">{{ detail.user.nickname || '—' }}</p>
          </el-card>
          <el-card shadow="never" class="mid-card mid-card-chart">
            <div class="mid-title">得分构成</div>
            <div class="pie-wrap">
              <div class="pie-disk" :style="{ background: compositionConic }" />
            </div>
            <div class="pie-legend">
              <span
                ><i class="dot dot-obj" />客观 {{ formatScore(compositionObjective) }}</span
              >
              <span
                ><i class="dot dot-sub" />主观 {{ formatScore(compositionSubjective) }}</span
              >
              <span v-if="pieShowStructuralDisk"
                ><i class="dot dot-sub-cap" />主观题满分占比（示意）</span
              >
            </div>
            <p v-if="compositionSum <= 0" class="mid-placeholder">暂无得分数据</p>
            <div v-if="paperMaxScore > 0" class="pie-analytics">
              <p class="pie-analytics-line">
                <span class="pa-label">客观正确率</span>
                {{ objectiveAccuracyLabel }}
                <span class="pa-sub"
                  >（{{ formatScore(compositionObjective) }} 分）</span
                >
              </p>
              <p v-if="subjectivePaperMax > 0" class="pie-analytics-line">
                <span class="pa-label">主观得分率</span>
                {{ subjectiveScoreRateLabel }}
                <span class="pa-sub"
                  >（{{ formatScore(compositionSubjective) }} /
                  {{ formatScore(subjectivePaperMax) }}）</span
                >
              </p>
            </div>
          </el-card>
          <el-card shadow="never" class="mid-card">
            <div class="mid-title">时间</div>
            <ul class="time-list">
              <li><span class="t-label">开考</span>{{ formatUtcText(detail.attempt.started_at) }}</li>
              <li><span class="t-label">截止</span>{{ formatUtcText(detail.attempt.deadline_at) }}</li>
              <li><span class="t-label">交卷</span>{{ formatUtcText(detail.attempt.submitted_at) }}</li>
              <li><span class="t-label">结束</span>{{ formatUtcText(detail.attempt.ended_at) }}</li>
            </ul>
          </el-card>
        </div>

        <el-card shadow="never" class="info-card muted-card">
          <template #header><span class="card-h-muted">基础信息</span></template>
          <el-descriptions :column="2" size="small" class="desc-plain">
            <el-descriptions-item label="会话 ID">{{ detail.attempt.id }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ statusText(detail.attempt.status) }}</el-descriptions-item>
            <el-descriptions-item label="试卷">{{ paperDisplayName }}</el-descriptions-item>
            <el-descriptions-item label="级别">{{ detail.paper.level }}</el-descriptions-item>
            <el-descriptions-item label="卷编号">{{ detail.paper.paper_id || '—' }}</el-descriptions-item>
            <el-descriptions-item label="含主观题">{{ detail.attempt.has_subjective ? '是' : '否' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card shadow="never" class="answers-card">
          <template #header><span class="card-h-strong">答题明细</span></template>
          <div class="ratio-bar" aria-hidden="true">
            <div
              class="ratio-seg ratio-ok"
              :style="{ width: barWidths.correctPct + '%' }"
            />
            <div
              class="ratio-seg ratio-bad"
              :style="{ width: barWidths.wrongPct + '%' }"
            />
            <div
              class="ratio-seg ratio-neutral"
              :style="{ width: barWidths.neutralPct + '%' }"
            />
          </div>
          <div class="answer-filter-row">
            <el-button-group size="small">
              <el-button :type="answerFilter === 'all' ? 'primary' : 'default'" @click="answerFilter = 'all'"
                >全部</el-button
              >
              <el-button :type="answerFilter === 'wrong' ? 'primary' : 'default'" @click="answerFilter = 'wrong'"
                >只看错题</el-button
              >
            </el-button-group>
          </div>
          <p class="ratio-hint">
            客观题分布：对 {{ objectiveCounts.correct }} · 错 {{ objectiveCounts.wrong }} · 其它
            {{ objectiveCounts.neutral }}
          </p>
          <div class="answer-list">
            <div
              v-for="(a, idx) in filteredAnswers"
              :key="answerRowKey(a, idx)"
              class="answer-item"
              :class="{
                'is-open': isRowExpanded(a, idx),
                'is-wrong-row': isObjectiveWrongRow(a),
              }"
              @click="toggleExpandRow(a, idx)"
            >
              <div class="answer-row-main">
                <span
                  class="q-badge"
                  :class="rowBadgeClass(a)"
                  >{{ a.question_no }}</span
                >
                <span class="q-type">{{
                  a.is_example ? '例题' : a.is_subjective ? '主观' : '客观'
                }}</span>
                <span v-if="a.section_title" class="q-section">{{ a.section_title }}</span>
                <div class="answer-row-trail">
                  <span v-if="a.is_subjective && !a.is_example" class="q-award">
                    {{ a.awarded_score != null ? formatScore(Number(a.awarded_score)) + ' 分' : '未评' }}
                  </span>
                  <span v-else-if="!a.is_example && !a.is_subjective" class="q-award">{{ a.score }} 分</span>
                  <span class="q-icon-wrap">
                    <el-icon v-if="rowIconKind(a) === 'ok'" class="ic-ok"><Check /></el-icon>
                    <el-icon v-else-if="rowIconKind(a) === 'bad'" class="ic-bad"><Close /></el-icon>
                    <el-icon v-else class="ic-muted"><Minus /></el-icon>
                  </span>
                  <el-icon class="chev" :class="{ 'chev-open': isRowExpanded(a, idx) }">
                    <ArrowRight />
                  </el-icon>
                </div>
              </div>
              <div
                v-show="isRowExpanded(a, idx)"
                class="answer-expand"
                @click.stop
              >
                <div class="exp-section">
                  <div class="exp-label">题干</div>
                  <div class="exp-body">{{ stemDisplayText(a) }}</div>
                </div>
                <div class="exp-section">
                  <div class="exp-label">你的答案</div>
                  <div class="exp-body exp-mono">{{ formatAnswerJson(a.answer_json) }}</div>
                </div>
                <template v-if="a.options && a.options.length > 0">
                  <div class="exp-section">
                    <div class="exp-label">选项</div>
                    <ul class="opt-list">
                      <li
                        v-for="o in a.options"
                        :key="o.id"
                        class="opt-row"
                        :class="{ 'opt-correct': o.is_correct === 1 }"
                      >
                        <span class="opt-flag">{{ o.flag }}</span>
                        <span class="opt-content">{{ optionContentLabel(o) }}</span>
                        <el-tag
                          v-if="o.is_correct === 1"
                          size="small"
                          type="success"
                          effect="plain"
                          class="opt-tag"
                          >标答</el-tag
                        >
                      </li>
                    </ul>
                  </div>
                </template>
                <div v-else class="exp-section exp-muted">
                  <div class="exp-label">选项</div>
                  <div class="exp-body">暂无选项数据</div>
                </div>
                <div v-if="a.analysis_text" class="exp-section">
                  <div class="exp-label">解析</div>
                  <div class="exp-body">{{ a.analysis_text }}</div>
                </div>
              </div>
            </div>
          </div>
        </el-card>

        <template v-if="canGradeSubjective">
          <el-card shadow="never" class="grade-card">
            <template #header><span class="card-h-strong">主观题评分</span></template>
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
              <el-button type="primary" :loading="savingScores" @click="saveSubjective"
                >保存主观分</el-button
              >
            </div>
          </el-card>
        </template>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import type { Component } from 'vue'
import { ElMessage } from 'element-plus'
import { Check, Close, Minus, Trophy, Medal, ChatDotRound, ArrowRight } from '@element-plus/icons-vue'
import { formatUtcForDisplay, formatUtcText } from '@/utils/datetime'
import {
  getAttemptList,
  getAttemptDetail,
  saveAttemptSubjectiveScores,
  type AttemptListItem,
  type AttemptDetail,
  type AttemptDetailAnswer,
  type AttemptDetailOption,
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

/** 统一 string key，避免接口返回 question_id 为字符串时与 number 比较导致展开态永远不命中 */
const expandedRowKeys = ref<string[]>([])
const answerFilter = ref<'all' | 'wrong'>('all')
const displayScore = ref(0)
const ringPercentage = ref(0)
let scoreAnimFrame = 0

const paperDisplayName = computed(() => {
  const p = detail.value?.paper
  if (!p) return '—'
  const n = (p.name ?? '').trim()
  if (n) return n
  return (p.title ?? '').trim() || '—'
})

const detailAnswers = computed(() => detail.value?.answers ?? [])

const subjectivePaperMax = computed(() =>
  detailAnswers.value
    .filter((a) => a.is_subjective && !a.is_example)
    .reduce((s, a) => s + (Number(a.score) || 0), 0),
)

const objectivePaperMax = computed(() =>
  detailAnswers.value
    .filter((a) => !a.is_subjective && !a.is_example)
    .reduce((s, a) => s + (Number(a.score) || 0), 0),
)

const filteredAnswers = computed(() => {
  const list = detailAnswers.value
  if (answerFilter.value !== 'wrong') return list
  return list.filter(
    (a) => !a.is_example && !a.is_subjective && a.objective_correct === false,
  )
})

function answerRowKey(a: AttemptDetailAnswer, idx: number) {
  const qid = a.question_id ?? idx
  const qno = a.question_no ?? idx
  return `q-${qid}-${qno}-${idx}`
}

/** 与列表行一一对应的展开键（含 idx 兜底，避免 id 缺失或重复） */
function expandRowKey(a: AttemptDetailAnswer, idx: number): string {
  const raw = a.question_id as unknown
  if (raw !== undefined && raw !== null && String(raw) !== '') {
    return `${String(raw)}#${idx}`
  }
  return `row-${idx}`
}

function isRowExpanded(a: AttemptDetailAnswer, idx: number) {
  return expandedRowKeys.value.includes(expandRowKey(a, idx))
}

function toggleExpandRow(a: AttemptDetailAnswer, idx: number) {
  const key = expandRowKey(a, idx)
  const cur = expandedRowKeys.value
  const i = cur.indexOf(key)
  expandedRowKeys.value = i >= 0 ? cur.filter((x) => x !== key) : [...cur, key]
}

const subjectiveRows = computed(() => {
  if (!detail.value) return []
  return (detail.value.answers ?? []).filter((a) => a.is_subjective && !a.is_example)
})

const canGradeSubjective = computed(() => {
  const a = detail.value?.attempt
  return !!(a && a.status === 4 && a.has_subjective === 1)
})

const paperMaxScore = computed(() => {
  if (!detail.value) return 0
  return (detail.value.answers ?? [])
    .filter((a) => !a.is_example)
    .reduce((s, a) => s + (Number(a.score) || 0), 0)
})

const scoreRatio = computed(() => {
  const max = paperMaxScore.value
  if (max <= 0 || !detail.value) return 0
  return detail.value.attempt.total_score / max
})

/** 与 scoreRatio 分档一致：低分暖色、高分绿色、中段品牌蓝 */
const ringStrokeColor = computed(() => {
  if (paperMaxScore.value <= 0) return '#2c5282'
  const r = scoreRatio.value
  if (r >= 0.9) return '#3d9a65'
  if (r >= 0.8) return '#2c5282'
  if (r >= 0.6) return '#d4943c'
  return '#d66a5c'
})

const gradeBand = computed((): {
  label: string
  icon: Component
  tagType: 'success' | 'warning' | 'info'
} | null => {
  if (!detail.value || paperMaxScore.value <= 0) return null
  const r = scoreRatio.value
  if (r >= 0.9) return { label: '优异', icon: Trophy, tagType: 'success' }
  if (r >= 0.8) return { label: '良好', icon: Medal, tagType: 'success' }
  if (r >= 0.6) return { label: '及格', icon: Medal, tagType: 'warning' }
  return { label: '再接再厉', icon: ChatDotRound, tagType: 'info' }
})

const durationLabel = computed(() => {
  const sec = detail.value?.attempt.duration_seconds
  if (sec == null || sec < 0) return '—'
  if (sec < 60) return `${sec} 秒`
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return s > 0 ? `${m} 分 ${s} 秒` : `${m} 分钟`
})

const compositionObjective = computed(() => detail.value?.attempt.objective_score ?? 0)
const compositionSubjective = computed(() => detail.value?.attempt.subjective_score ?? 0)
const compositionSum = computed(
  () => compositionObjective.value + compositionSubjective.value,
)

const pieShowStructuralDisk = computed(() => {
  const att = detail.value?.attempt
  return !!(
    att &&
    att.has_subjective === 1 &&
    subjectivePaperMax.value > 0 &&
    compositionSubjective.value === 0
  )
})

function clampDeg(deg: number) {
  return Math.max(0, Math.min(360, deg))
}

/** 有主观满分但主观得分为 0 时按试卷满分拆环，避免整圆单色 */
const compositionConic = computed(() => {
  const p = paperMaxScore.value
  if (p <= 0) {
    return 'conic-gradient(var(--result-muted-bg) 0deg, var(--result-muted-bg) 360deg)'
  }
  if (pieShowStructuralDisk.value) {
    const objGot = compositionObjective.value
    const subMax = subjectivePaperMax.value
    let d1 = clampDeg((objGot / p) * 360)
    let d2 = clampDeg((subMax / p) * 360)
    const sum12 = d1 + d2
    if (sum12 > 360) {
      const k = 360 / sum12
      d1 *= k
      d2 *= k
    }
    const e1 = d1
    const e2 = d1 + d2
    return `conic-gradient(var(--result-brand) 0deg ${e1}deg, var(--result-sub-empty) ${e1}deg ${e2}deg, var(--result-pie-rest) ${e2}deg 360deg)`
  }
  const o = compositionObjective.value
  const s = compositionSubjective.value
  const t = o + s
  if (t <= 0) {
    return 'conic-gradient(var(--result-muted-bg) 0deg, var(--result-muted-bg) 360deg)'
  }
  const deg = (o / t) * 360
  return `conic-gradient(var(--result-brand) 0deg ${deg}deg, var(--result-subjective) ${deg}deg 360deg)`
})

const objectiveCounts = computed(() => {
  let correct = 0
  let wrong = 0
  let neutral = 0
  if (!detail.value) return { correct, wrong, neutral }
  for (const a of detail.value.answers ?? []) {
    if (a.is_example) {
      neutral++
      continue
    }
    if (a.is_subjective) {
      neutral++
      continue
    }
    if (a.objective_correct === true) correct++
    else if (a.objective_correct === false) wrong++
    else neutral++
  }
  return { correct, wrong, neutral }
})

const objectiveAccuracyLabel = computed(() => {
  const { correct, wrong } = objectiveCounts.value
  const d = correct + wrong
  if (d <= 0) return '—'
  return `${((correct / d) * 100).toFixed(1)}%（${correct} / ${d} 题）`
})

const subjectiveScoreRateLabel = computed(() => {
  const max = subjectivePaperMax.value
  if (max <= 0) return '—'
  const got = compositionSubjective.value
  return `${((got / max) * 100).toFixed(1)}%`
})

const barWidths = computed(() => {
  const { correct, wrong, neutral } = objectiveCounts.value
  const sum = correct + wrong + neutral
  if (sum <= 0) return { correctPct: 0, wrongPct: 0, neutralPct: 100 }
  return {
    correctPct: (correct / sum) * 100,
    wrongPct: (wrong / sum) * 100,
    neutralPct: (neutral / sum) * 100,
  }
})

function cancelScoreAnimation() {
  if (scoreAnimFrame) {
    cancelAnimationFrame(scoreAnimFrame)
    scoreAnimFrame = 0
  }
}

function startScoreAnimation() {
  cancelScoreAnimation()
  if (!detail.value) return
  const target = Number(detail.value.attempt.total_score) || 0
  const max = paperMaxScore.value
  const targetPct = max > 0 ? Math.min(100, (target / max) * 100) : 0
  const start = performance.now()
  const dur = 780
  displayScore.value = 0
  ringPercentage.value = 0
  function easeOutQuad(t: number) {
    return 1 - (1 - t) * (1 - t)
  }
  function tick(now: number) {
    const t = Math.min(1, (now - start) / dur)
    const k = easeOutQuad(t)
    displayScore.value = target * k
    ringPercentage.value = targetPct * k
    if (t < 1) {
      scoreAnimFrame = requestAnimationFrame(tick)
    } else {
      displayScore.value = target
      ringPercentage.value = targetPct
      scoreAnimFrame = 0
    }
  }
  scoreAnimFrame = requestAnimationFrame(tick)
}

watch(
  () => [drawer.value, detail.value?.attempt.id] as const,
  ([open]) => {
    if (open && detail.value) {
      expandedRowKeys.value = []
      answerFilter.value = 'all'
      nextTick(() => startScoreAnimation())
    } else {
      cancelScoreAnimation()
    }
  },
)

onUnmounted(() => cancelScoreAnimation())

function rowBadgeClass(a: AttemptDetailAnswer) {
  if (a.is_example || a.is_subjective) return 'badge-neutral'
  if (a.objective_correct === true) return 'badge-ok'
  if (a.objective_correct === false) return 'badge-bad'
  return 'badge-neutral'
}

function rowIconKind(a: AttemptDetailAnswer): 'ok' | 'bad' | 'muted' {
  if (a.is_example || a.is_subjective) return 'muted'
  if (a.objective_correct === true) return 'ok'
  if (a.objective_correct === false) return 'bad'
  return 'muted'
}

function isObjectiveWrongRow(a: AttemptDetailAnswer) {
  return !a.is_example && !a.is_subjective && a.objective_correct === false
}

function stemDisplayText(a: AttemptDetailAnswer) {
  const s = (a.stem_text || '').trim()
  if (s) return s
  return '（无文字题干；听力/图片等题型可能仅有资源文件名，见下方选项或试卷资源）'
}

function optionContentLabel(o: AttemptDetailOption) {
  const c = (o.content || '').trim()
  if (!c) return '—'
  const t = (o.option_type || '').toLowerCase()
  if (t === 'image' || /\.(jpe?g|png|gif|webp|svg)(\?|$)/i.test(c)) {
    return `[图片] ${c}`
  }
  if (t === 'audio' || /\.(mp3|wav|m4a|ogg)(\?|$)/i.test(c)) {
    return `[音频] ${c}`
  }
  return c
}

function formatScore(n: number) {
  if (Number.isNaN(n)) return '—'
  return n.toFixed(2)
}

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
    const sel = o.selected_option_ids
    if (Array.isArray(sel) && sel.length > 0) {
      return `选项 ID: ${sel.map((x) => String(x)).join('、')}`
    }
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
        const qid = Number(a.question_id)
        if (!Number.isNaN(qid)) {
          m[qid] = a.awarded_score != null ? Number(a.awarded_score) : 0
        }
      }
    }
    subjectiveScores.value = m
    answerFilter.value = 'all'
    expandedRowKeys.value = []
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
    nextTick(() => startScoreAnimation())
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
  --result-brand: #2c5282;
  --result-subjective: #c5a880;
  --result-ok: #5aa67e;
  --result-bad: #c97b72;
  --result-muted-bg: #c9ccd4;
  --result-ok-soft: #e8f2ec;
  --result-bad-soft: #f5e8e6;
  --result-neutral-soft: #eceef2;
  --result-sub-empty: #c8ccd6;
  --result-pie-rest: #e4e7ee;
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
.list-total-score {
  font-weight: 700;
  color: var(--result-brand);
}
.hint {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin: 0 0 10px;
}
.grade-actions {
  margin-top: 12px;
}

.result-detail {
  margin-top: -8px;
  padding-bottom: 24px;
}

.detail-hero {
  margin: 0 -20px 20px;
  padding: 28px 20px 22px;
  background: linear-gradient(165deg, #eef2f7 0%, #f7f5fb 48%, #eef4f2 100%);
  border-radius: 0 0 12px 12px;
}

.hero-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.hero-ring-wrap {
  position: relative;
  width: 148px;
  height: 148px;
}

.hero-ring-wrap :deep(.el-progress-circle) {
  display: block;
}

.hero-score-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  padding-top: 4px;
}

.hero-score {
  font-size: 2.25rem;
  font-weight: 800;
  line-height: 1.1;
  color: var(--result-brand);
  letter-spacing: -0.02em;
}

.hero-max {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--el-text-color-secondary);
  margin-top: 2px;
}

.hero-meta-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 8px 14px;
  min-height: 28px;
  margin-top: 2px;
}

.grade-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.grade-ic {
  width: 16px;
  height: 16px;
}

.hero-duration-inline {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  white-space: nowrap;
}

.detail-mid-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 14px;
}

@media (max-width: 720px) {
  .detail-mid-grid {
    grid-template-columns: 1fr;
  }
}

.mid-card {
  border-radius: 10px;
}

.mid-card :deep(.el-card__body) {
  padding: 14px 16px;
}

.mid-title {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 8px;
  text-transform: none;
}

.mid-strong {
  margin: 0 0 4px;
  font-size: 15px;
  font-weight: 600;
}

.mid-muted {
  margin: 0;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.mid-placeholder {
  margin: 8px 0 0;
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}

.pie-wrap {
  display: flex;
  justify-content: center;
  margin: 8px 0 10px;
}

.pie-disk {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  box-shadow: inset 0 0 0 10px var(--el-fill-color-blank);
}

.pie-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 16px;
  font-size: 12px;
  color: var(--el-text-color-regular);
}

.pie-legend .dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 6px;
  vertical-align: middle;
}

.dot-obj {
  background: var(--result-brand);
}
.dot-sub {
  background: var(--result-subjective);
}

.dot-sub-cap {
  background: var(--result-sub-empty);
}

.pie-analytics {
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px dashed var(--el-border-color-lighter);
  font-size: 12px;
  line-height: 1.5;
}

.pie-analytics-line {
  margin: 4px 0 0;
  color: var(--el-text-color-regular);
}

.pie-analytics-line:first-child {
  margin-top: 0;
}

.pa-label {
  font-weight: 600;
  color: var(--el-text-color-secondary);
  margin-right: 6px;
}

.pa-sub {
  color: var(--el-text-color-placeholder);
  font-size: 11px;
  margin-left: 4px;
}

.time-list {
  list-style: none;
  margin: 0;
  padding: 0;
  font-size: 13px;
}

.time-list li {
  margin-bottom: 6px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.t-label {
  color: var(--el-text-color-secondary);
  min-width: 36px;
}

.info-card,
.answers-card,
.grade-card {
  margin-bottom: 14px;
  border-radius: 10px;
}

.answers-card :deep(.el-card__body) {
  overflow: visible;
}

.muted-card :deep(.el-card__header) {
  padding: 10px 16px;
  background: var(--el-fill-color-light);
}

.card-h-muted {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  font-weight: 600;
}

.card-h-strong {
  font-size: 14px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

.desc-plain :deep(.el-descriptions__label) {
  color: var(--el-text-color-secondary);
  font-weight: 500;
}

.ratio-bar {
  display: flex;
  height: 8px;
  border-radius: 999px;
  overflow: hidden;
  background: var(--el-fill-color-light);
  margin-bottom: 8px;
}

.ratio-seg {
  height: 100%;
  transition: width 0.65s ease-out;
}

.ratio-ok {
  background: var(--result-ok);
}

.ratio-bad {
  background: var(--result-bad);
}

.ratio-neutral {
  background: var(--result-muted-bg);
}

.answer-filter-row {
  margin: 8px 0 6px;
}

.ratio-hint {
  margin: 0 0 14px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.answer-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 480px;
  overflow-y: auto;
  padding-right: 4px;
}

.answer-list > .answer-item {
  flex-shrink: 0;
}

.answer-item {
  display: flex;
  flex-direction: column;
  isolation: isolate;
  border-radius: 10px;
  background: var(--el-fill-color-blank);
  border: 1px solid var(--el-border-color-lighter);
  overflow: hidden;
  min-height: 48px;
  cursor: pointer;
  transition:
    border-color 0.2s,
    box-shadow 0.2s;
}

.answer-item:hover {
  border-color: var(--el-border-color);
}

.answer-item.is-open {
  border-color: color-mix(in srgb, var(--result-brand) 35%, var(--el-border-color));
  box-shadow: 0 2px 10px rgb(44 82 130 / 0.08);
}

.answer-item.is-wrong-row {
  background: rgb(201 123 114 / 0.08);
}

.answer-item.is-wrong-row.is-open {
  background: rgb(201 123 114 / 0.1);
}

.answer-row-main {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 12px 12px 14px;
  flex-wrap: nowrap;
  flex-shrink: 0;
  min-height: 48px;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  position: relative;
  z-index: 1;
  background: inherit;
}

.answer-row-trail {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: 10px;
  margin-left: auto;
  padding-left: 8px;
}

.q-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  height: 32px;
  padding: 0 6px;
  border-radius: 50%;
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 900;
  letter-spacing: -0.02em;
  color: #fff;
}

.badge-ok {
  background: var(--result-ok);
}

.badge-bad {
  background: var(--result-bad);
}

.badge-neutral {
  background: var(--result-muted-bg);
  color: var(--el-text-color-regular);
}

.q-type {
  font-size: 12px;
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}

.q-section {
  font-size: 12px;
  flex: 1 1 auto;
  min-width: 0;
  color: var(--el-text-color-regular);
  max-width: min(320px, 50vw);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.q-award {
  font-size: 13px;
  flex-shrink: 0;
  font-weight: 600;
  color: var(--el-text-color-primary);
  text-align: right;
  min-width: 3.5em;
}

.q-icon-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 24px;
  height: 24px;
}

.ic-ok,
.ic-bad,
.ic-muted {
  font-size: 20px;
}

.ic-ok {
  color: var(--result-ok);
}

.ic-bad {
  color: var(--result-bad);
}

.ic-muted {
  color: var(--result-muted-bg);
}

.chev {
  font-size: 18px;
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  color: var(--el-text-color-placeholder);
  transition: transform 0.25s ease;
}

.chev-open {
  transform: rotate(90deg);
}

.answer-expand {
  flex: 0 0 auto;
  width: 100%;
  box-sizing: border-box;
  padding: 0 14px 14px 58px;
  border-top: 1px dashed var(--el-border-color-lighter);
  font-size: 13px;
  line-height: 1.55;
  color: var(--el-text-color-regular);
  position: relative;
  z-index: 0;
  word-break: break-word;
  overflow-wrap: anywhere;
}

.exp-section {
  margin-top: 12px;
}

.exp-section:first-child {
  margin-top: 10px;
}

.exp-body {
  word-break: break-word;
  overflow-wrap: anywhere;
}

.exp-mono {
  font-family: ui-monospace, 'Cascadia Code', 'Consolas', monospace;
  font-size: 12px;
}

.exp-muted {
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}

.exp-label {
  display: block;
  font-size: 11px;
  font-weight: 600;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.opt-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.opt-row {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 6px 8px;
  margin-bottom: 8px;
  padding: 6px 8px;
  border-radius: 8px;
  background: color-mix(in srgb, var(--el-fill-color) 92%, transparent);
}

.opt-content {
  flex: 1 1 140px;
  min-width: 0;
  word-break: break-word;
  overflow-wrap: anywhere;
}

.opt-tag {
  flex-shrink: 0;
}

.opt-correct {
  font-weight: 500;
}

.opt-flag {
  font-weight: 700;
  flex-shrink: 0;
  color: var(--result-brand);
}
</style>
