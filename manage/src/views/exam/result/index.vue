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
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <div class="result-list-ops">
              <el-button link type="primary" @click="openDetail(row)">详情</el-button>
              <el-tooltip
                v-if="showSubjectiveGradeButton(row)"
                content="主观题已评分，不可再次修改"
                placement="top"
                :disabled="!isSubjectiveGradeButtonDisabled(row)"
                :show-after="200"
              >
                <span class="result-list-op__sub-clip">
                  <el-button
                    class="result-list-op--subjective"
                    link
                    type="primary"
                    :disabled="isSubjectiveGradeButtonDisabled(row)"
                    @click="openSubjectiveGrade(row)"
                  >
                    主观题评分
                  </el-button>
                </span>
              </el-tooltip>
            </div>
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
          <div v-if="groupedResultSections.length" class="paper-answer-by-section">
            <div v-if="groupedResultSections.length === 1" class="paper-answer-sec-body">
              <ExamQuestionReviewCard
                v-for="(row, ri) in groupedResultSections[0]!.rows"
                :key="answerRowKey(row.answer, ri)"
                mode="review"
                :question-no="row.answer.question_no"
                :score="Number(row.answer.score) || 0"
                :is-example="row.answer.is_example"
                :is-subjective="row.answer.is_subjective"
                :stem-text="row.answer.stem_text || ''"
                :screen-text-json="row.answer.screen_text_json || ''"
                :topic-json="row.meta?.topicJson || ''"
                :block-index="row.meta?.blockIndex ?? 0"
                :question-index="row.meta?.questionIndex ?? 0"
                :block-passage-text="
                  row.meta
                    ? blockReadingPassageFromTopic(row.meta.topicJson, row.meta.blockIndex)
                    : ''
                "
                :show-block-passage="row.showBlockPassage"
                :source-base-url="resolvedSourceBaseUrl"
                :audio-file="questionAudioById.get(row.answer.question_id) || ''"
                :options="row.optionsForCard"
                :show-correct-options="true"
                :answer-json="row.answer.answer_json || ''"
                :objective-correct="row.answer.objective_correct"
                :awarded-score="row.answer.awarded_score"
                :analysis-text="row.answer.analysis_text || ''"
              />
            </div>
            <el-tabs v-else v-model="resultAnswersActiveTab" class="result-sec-tabs">
              <el-tab-pane
                v-for="(sec, si) in groupedResultSections"
                :key="'sec-tab-' + sec.sectionId + '-' + si"
                :label="sec.title"
                :name="'sec-' + sec.sectionId + '-' + si"
              >
                <div class="paper-answer-sec-body">
                  <ExamQuestionReviewCard
                    v-for="(row, ri) in sec.rows"
                    :key="answerRowKey(row.answer, ri)"
                    mode="review"
                    :question-no="row.answer.question_no"
                    :score="Number(row.answer.score) || 0"
                    :is-example="row.answer.is_example"
                    :is-subjective="row.answer.is_subjective"
                    :stem-text="row.answer.stem_text || ''"
                    :screen-text-json="row.answer.screen_text_json || ''"
                    :topic-json="row.meta?.topicJson || ''"
                    :block-index="row.meta?.blockIndex ?? 0"
                    :question-index="row.meta?.questionIndex ?? 0"
                    :block-passage-text="
                      row.meta
                        ? blockReadingPassageFromTopic(row.meta.topicJson, row.meta.blockIndex)
                        : ''
                    "
                    :show-block-passage="row.showBlockPassage"
                    :source-base-url="resolvedSourceBaseUrl"
                    :audio-file="questionAudioById.get(row.answer.question_id) || ''"
                    :options="row.optionsForCard"
                    :show-correct-options="true"
                    :answer-json="row.answer.answer_json || ''"
                    :objective-correct="row.answer.objective_correct"
                    :awarded-score="row.answer.awarded_score"
                    :analysis-text="row.answer.analysis_text || ''"
                  />
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
          <p v-else class="ratio-hint paper-answer-empty">暂无符合条件题目。</p>
        </el-card>
      </div>
    </el-drawer>

    <el-dialog
      v-model="subjectiveDlgVisible"
      class="subjective-grade-dialog"
      title="主观题评分"
      width="min(960px, 98vw)"
      top="4vh"
      destroy-on-close
      append-to-body
      :close-on-click-modal="false"
      @closed="onSubjectiveGradeDlgClosed"
    >
      <p class="subjective-dlg-hint">每题得分不超过该题满分；保存后更新主观分与总分。</p>
      <p v-if="subjectiveDlgMetaLine" class="subjective-dlg-meta">{{ subjectiveDlgMetaLine }}</p>
      <el-skeleton v-if="subjectiveDlgLoading" :rows="5" animated />
      <el-table
        v-else
        :data="subjectiveGradingDialogRows"
        border
        stripe
        size="small"
        class="subjective-dlg-table"
        :max-height="480"
      >
        <el-table-column label="大题" min-width="100" show-overflow-tooltip>
          <template #default="{ row }">{{ row.sectionTitle }}</template>
        </el-table-column>
        <el-table-column prop="question_no" label="题号" width="64" />
        <el-table-column label="题干" min-width="220">
          <template #default="{ row }">
            <div v-if="row.stemIsBlocks" class="subjective-stem-block subjective-stem-wrap">
              <ExamScreenTextBlocks
                :raw="row.stemBlocksRaw"
                :source-base-url="subjectiveDlgSourceBase"
              />
            </div>
            <div
              v-else-if="row.stemHtml"
              class="subjective-stem-wrap exam-rich-html"
              v-html="row.stemHtml"
            />
            <span v-else class="subjective-stem-plain">{{ row.stemPreview || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="作答内容" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">{{ row.answerText || '—' }}</template>
        </el-table-column>
        <el-table-column label="满分" width="72" align="right">
          <template #default="{ row }">{{ row.score }}</template>
        </el-table-column>
        <el-table-column label="得分" width="128" fixed="right" align="center">
          <template #default="{ row }">
            <el-input-number
              v-model="subjectiveDlgScores[row.question_id]"
              :min="0"
              :max="row.score"
              :precision="2"
              :step="0.5"
              size="small"
              controls-position="right"
              style="width: 110px"
            />
          </template>
        </el-table-column>
      </el-table>
      <p v-if="!subjectiveDlgLoading && !subjectiveGradingDialogRows.length" class="subjective-dlg-empty">
        暂无主观题可评
      </p>
      <template #footer>
        <el-button @click="subjectiveDlgVisible = false">取消</el-button>
        <el-button
          type="primary"
          :disabled="!subjectiveGradingDialogRows.length"
          :loading="savingSubjectiveScores"
          @click="saveSubjectiveFromDialog"
        >
          保存主观分
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import type { Component } from 'vue'
import { ElMessage } from 'element-plus'
import { Trophy, Medal, ChatDotRound } from '@element-plus/icons-vue'
import { formatUtcForDisplay, formatUtcText } from '@/utils/datetime'
import {
  getAttemptList,
  getAttemptDetail,
  saveAttemptSubjectiveScores,
  type AttemptListItem,
  type AttemptDetail,
  type AttemptDetailAnswer,
} from '@/api/examAttempt'
import { getExamPaperDetail, type ExamPaperDetail } from '@/api/exam'
import ExamQuestionReviewCard from '@/components/exam/ExamQuestionReviewCard.vue'
import ExamScreenTextBlocks from '@/components/exam/ExamScreenTextBlocks.vue'
import {
  formatAnswerJson,
  looksLikeHtml,
  sanitizeHtmlForDisplay,
  type ExamOptionDisplayRow,
} from '@/utils/examDisplay'
import {
  buildQuestionTopicMetaById,
  blockReadingPassageFromTopic,
  effectivePaperQuestionStem,
  hasQuestionStemSegments,
  questionStemRichRaw,
  type QuestionTopicMeta,
} from '@/utils/examPaperQuestionDisplay'
import { getMockExaminationPapers, type MockExaminationPaperItem } from '@/api/mockAdmin'

type SubjectiveGradeTableRow = AttemptDetailAnswer & {
  sectionTitle: string
  stemPreview: string
  answerText: string
  stemIsBlocks: boolean
  stemBlocksRaw: unknown
  /** 经 sanitize 的 html，与 ExamQuestionStem 的「纯 html 题干」一致 */
  stemHtml: string
}

/** 与详情抽屉无关，供主观题评分解析列表（大题/题干/作答/满分） */
function buildSubjectiveGradingRows(
  d: AttemptDetail | null,
  paper: ExamPaperDetail | null,
): SubjectiveGradeTableRow[] {
  if (!d) return []
  const list = (d.answers ?? []).filter((a) => a.is_subjective && !a.is_example)
  const metaMap = buildQuestionTopicMetaById(paper)
  const orderMap = new Map<number, number>()
  paper?.sections?.forEach((s, i) => {
    orderMap.set(s.id, s.sort_order ?? i)
  })
  const sectionTitleById = new Map<number, string>()
  for (const s of paper?.sections ?? []) {
    sectionTitleById.set(s.id, (s.topic_title || '').trim() || `大题 #${s.id}`)
  }
  for (const a of d.answers ?? []) {
    const sid = Number(a.section_id) || 0
    if (sid && !sectionTitleById.has(sid) && (a.section_title || '').trim()) {
      sectionTitleById.set(sid, (a.section_title || '').trim())
    }
  }
  return list
    .slice()
    .sort((a, b) => {
      const sa = orderMap.get(Number(a.section_id)) ?? 10_000
      const sb = orderMap.get(Number(b.section_id)) ?? 10_000
      if (sa !== sb) return sa - sb
      return (Number(a.question_no) || 0) - (Number(b.question_no) || 0)
    })
    .map((a) => {
      const sid = Number(a.section_id) || 0
      const st = (a.section_title || '').trim()
      const fromPaper = sectionTitleById.get(sid)
      const sectionTitle = (fromPaper || st || (sid ? `大题 #${sid}` : '—')) as string
      const meta = metaMap.get(Number(a.question_id))
      const bi = meta?.blockIndex ?? 0
      const qi = meta?.questionIndex ?? 0
      const tj = meta?.topicJson

      let stemIsBlocks = false
      let stemBlocksRaw: unknown
      let stemHtml = ''
      if (hasQuestionStemSegments(a.stem_text, a.screen_text_json, tj, bi, qi)) {
        const raw = questionStemRichRaw(a.stem_text, a.screen_text_json, tj, bi, qi)
        if (raw != null) {
          stemIsBlocks = true
          stemBlocksRaw = raw
        }
      }
      if (!stemIsBlocks) {
        const h = (a.stem_text || '').trim()
        if (looksLikeHtml(h)) {
          stemHtml = sanitizeHtmlForDisplay(h)
        } else {
          const scj = (a.screen_text_json || '').trim()
          if (scj && looksLikeHtml(scj) && !/^\s*\[/.test(scj)) {
            stemHtml = sanitizeHtmlForDisplay(scj)
          }
        }
      }

      const stemPreview = effectivePaperQuestionStem(
        a.stem_text,
        a.screen_text_json,
        meta?.topicJson,
        bi,
        qi,
      )
      const answerText = formatAnswerJson(a.answer_json || '')
      return { ...a, sectionTitle, stemPreview, answerText, stemIsBlocks, stemBlocksRaw, stemHtml }
    })
}

const loading = ref(false)
const rows = ref<AttemptListItem[]>([])
const total = ref(0)
const route = useRoute()

const query = reactive({
  page: 1,
  size: 10,
  username: '',
  level: '',
  status: 0 as number,
  examination_paper_id: 0 as number,
  exam_batch_id: 0 as number,
  /** 0 不限；1 仅待主观题评阅 */
  subjective_pending: 0 as number,
})

const paperSel = ref<number | undefined>(undefined)
const paperOptions = ref<MockExaminationPaperItem[]>([])

const drawer = ref(false)
const detail = ref<AttemptDetail | null>(null)
/** 打开详情后按需拉取，用于 topic_json / 音频等与试卷对齐 */
const paperDetailForResult = ref<ExamPaperDetail | null>(null)
const detailId = ref(0)

const subjectiveDlgVisible = ref(false)
const subjectiveDlgLoading = ref(false)
const subjectiveDlgAttemptId = ref(0)
const subjectiveGradeDetail = ref<AttemptDetail | null>(null)
const subjectiveGradePaper = ref<ExamPaperDetail | null>(null)
const subjectiveDlgScores = ref<Record<number, number>>({})
const savingSubjectiveScores = ref(false)

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

interface ResultAnswerRow {
  answer: AttemptDetailAnswer
  meta?: QuestionTopicMeta
  showBlockPassage: boolean
  optionsForCard: ExamOptionDisplayRow[]
}

function answerRowKey(a: AttemptDetailAnswer, idx: number) {
  const qid = a.question_id ?? idx
  const qno = a.question_no ?? idx
  return `q-${qid}-${qno}-${idx}`
}

const topicMetaByQuestionId = computed(() => buildQuestionTopicMetaById(paperDetailForResult.value))

const questionAudioById = computed(() => {
  const m = new Map<number, string>()
  for (const sec of paperDetailForResult.value?.sections ?? []) {
    for (const b of sec.blocks ?? []) {
      for (const q of b.questions ?? []) {
        const af = (q.audio_file || '').trim()
        if (af) m.set(q.id, af)
      }
    }
  }
  return m
})

const resolvedSourceBaseUrl = computed(() => {
  const fromAttempt = (detail.value?.paper?.source_base_url || '').trim()
  if (fromAttempt) return fromAttempt
  return (paperDetailForResult.value?.paper?.source_base_url || '').trim()
})

const groupedResultSections = computed(() => {
  const list = filteredAnswers.value
  if (!list.length) return [] as { sectionId: number; title: string; rows: ResultAnswerRow[] }[]
  const paper = paperDetailForResult.value
  const sectionSort = new Map<number, number>()
  const sectionTitlePaper = new Map<number, string>()
  if (paper?.sections?.length) {
    for (const s of paper.sections) {
      sectionSort.set(s.id, s.sort_order ?? 0)
      const tt = (s.topic_title || '').trim()
      if (tt) sectionTitlePaper.set(s.id, tt)
    }
  }
  const bySec = new Map<number, AttemptDetailAnswer[]>()
  for (const a of list) {
    const sid = Number(a.section_id) || 0
    if (!bySec.has(sid)) bySec.set(sid, [])
    bySec.get(sid)!.push(a)
  }
  const metaMap = topicMetaByQuestionId.value
  const ids = [...bySec.keys()].sort((x, y) => {
    const hasX = sectionSort.has(x)
    const hasY = sectionSort.has(y)
    const ox = hasX ? (sectionSort.get(x) ?? x) : x + 1_000_000
    const oy = hasY ? (sectionSort.get(y) ?? y) : y + 1_000_000
    if (ox !== oy) return ox - oy
    return x - y
  })
  return ids.map((sectionId) => {
    const answers = (bySec.get(sectionId) || []).slice().sort((p, q) => {
      return (Number(p.question_no) || 0) - (Number(q.question_no) || 0)
    })
    let prevBlockKey: string | null = null
    const rows: ResultAnswerRow[] = answers.map((answer) => {
      const meta = metaMap.get(Number(answer.question_id))
      const blockKey = meta != null ? String(meta.blockIndex) : `na:${answer.question_id}`
      const showBlock = prevBlockKey !== blockKey
      prevBlockKey = blockKey
      const optionsForCard: ExamOptionDisplayRow[] = (answer.options ?? []).map((o) => ({
        id: o.id,
        flag: o.flag,
        content: o.content,
        is_correct: o.is_correct,
        option_type: o.option_type,
        sort_order: o.sort_order,
      }))
      return { answer, meta, showBlockPassage: showBlock, optionsForCard }
    })
    const title =
      sectionTitlePaper.get(sectionId) ||
      (answers[0]?.section_title || '').trim() ||
      `大题 #${sectionId}`
    return { sectionId, title, rows }
  })
})

const subjectiveGradingDialogRows = computed(() =>
  buildSubjectiveGradingRows(subjectiveGradeDetail.value, subjectiveGradePaper.value),
)

const subjectiveDlgSourceBase = computed(() => {
  const a = (subjectiveGradeDetail.value?.paper?.source_base_url || '').trim()
  if (a) return a
  return (subjectiveGradePaper.value?.paper?.source_base_url || '').trim()
})

const subjectiveDlgMetaLine = computed(() => {
  const d = subjectiveGradeDetail.value
  if (!d) return ''
  const u = d.user
  const name = (u.nickname || '').trim() || u.username || '—'
  const paperN = (d.paper.name || d.paper.title || '').trim() || '—'
  return `学员 ${u.username}（${name}）· 试卷 ${paperN} · 会话 #${d.attempt.id}`
})

const resultAnswersActiveTab = ref('')

watch(
  () => groupedResultSections.value,
  (sections) => {
    if (sections.length <= 1) return
    const names = sections.map((s, i) => `sec-${s.sectionId}-${i}`)
    if (!resultAnswersActiveTab.value || !names.includes(resultAnswersActiveTab.value)) {
      resultAnswersActiveTab.value = names[0] ?? 'sec-0-0'
    }
  },
  { immediate: true },
)

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
      answerFilter.value = 'all'
      nextTick(() => startScoreAnimation())
    } else {
      cancelScoreAnimation()
    }
  },
)

onUnmounted(() => cancelScoreAnimation())

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
      subjective_pending: query.subjective_pending === 1 ? 1 : undefined,
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
  query.subjective_pending = 0
  paperSel.value = undefined
  loadList()
}

async function openDetail(row: AttemptListItem) {
  detailId.value = row.id
  drawer.value = true
  paperDetailForResult.value = null
  try {
    const res = (await getAttemptDetail(row.id)) as { data?: AttemptDetail }
    detail.value = res?.data ?? null
    answerFilter.value = 'all'
    const pid = Number(detail.value?.paper?.exam_paper_id)
    if (detail.value && pid > 0 && !Number.isNaN(pid)) {
      try {
        const pr = (await getExamPaperDetail(pid)) as { data?: ExamPaperDetail }
        paperDetailForResult.value = pr?.data ?? null
      } catch {
        paperDetailForResult.value = null
      }
    }
  } catch {
    detail.value = null
    paperDetailForResult.value = null
  }
}

function onSubjectiveGradeDlgClosed() {
  subjectiveDlgAttemptId.value = 0
  subjectiveGradeDetail.value = null
  subjectiveGradePaper.value = null
  subjectiveDlgScores.value = {}
}

/** 含主观题且已结束会话时显示「主观题评分」入口；已评过分则仅置灰不可点 */
function showSubjectiveGradeButton(row: AttemptListItem) {
  return row.has_subjective === 1 && row.status === 4
}

/** 列表行主观题是否已评（兼容 camelCase 等） */
function rowSubjectiveGradedFlag(row: AttemptListItem) {
  const r = row as unknown as Record<string, unknown>
  const v = r.subjective_graded ?? r.subjectiveGraded
  if (v === true || v === 'true' || v === 1 || v === '1') return 1
  return Number(v) === 1 ? 1 : 0
}

function isSubjectiveGradeButtonDisabled(row: AttemptListItem) {
  return rowSubjectiveGradedFlag(row) === 1
}

function isDetailSubjectiveAlreadyGraded(d: AttemptDetail) {
  for (const a of d.answers ?? []) {
    if (a.is_subjective && !a.is_example && a.awarded_score != null) return true
  }
  return false
}

async function openSubjectiveGrade(row: AttemptListItem) {
  if (row.has_subjective !== 1) {
    ElMessage.info('该场次未包含主观题')
    return
  }
  if (row.status !== 4) {
    ElMessage.warning('仅「已结束」的会话可评分')
    return
  }
  if (rowSubjectiveGradedFlag(row) === 1) {
    ElMessage.warning('主观题已评过分，不可再次修改')
    return
  }
  subjectiveDlgAttemptId.value = row.id
  subjectiveDlgVisible.value = true
  subjectiveDlgLoading.value = true
  subjectiveGradeDetail.value = null
  subjectiveGradePaper.value = null
  subjectiveDlgScores.value = {}
  try {
    const res = (await getAttemptDetail(row.id)) as { data?: AttemptDetail }
    const d = res?.data ?? null
    if (!d) {
      ElMessage.error('未加载到会话详情')
      subjectiveDlgVisible.value = false
      return
    }
    if (isDetailSubjectiveAlreadyGraded(d)) {
      ElMessage.warning('主观题已评过分，不可再次修改')
      subjectiveDlgVisible.value = false
      return
    }
    subjectiveGradeDetail.value = d
    const m: Record<number, number> = {}
    for (const a of d.answers ?? []) {
      if (a.is_subjective && !a.is_example) {
        const qid = Number(a.question_id)
        if (!Number.isNaN(qid)) {
          m[qid] = a.awarded_score != null ? Number(a.awarded_score) : 0
        }
      }
    }
    subjectiveDlgScores.value = m
    const pid = Number(d.paper?.exam_paper_id)
    if (pid > 0 && !Number.isNaN(pid)) {
      try {
        const pr = (await getExamPaperDetail(pid)) as { data?: ExamPaperDetail }
        subjectiveGradePaper.value = pr?.data ?? null
      } catch {
        subjectiveGradePaper.value = null
      }
    }
  } catch {
    ElMessage.error('加载失败')
    subjectiveDlgVisible.value = false
  } finally {
    subjectiveDlgLoading.value = false
  }
}

async function saveSubjectiveFromDialog() {
  const attemptId = subjectiveDlgAttemptId.value
  if (!attemptId || !subjectiveGradeDetail.value) return
  const items = buildSubjectiveGradingRows(
    subjectiveGradeDetail.value,
    subjectiveGradePaper.value,
  ).map((r) => ({
    question_id: r.question_id,
    score: subjectiveDlgScores.value[r.question_id] ?? 0,
  }))
  if (items.length === 0) {
    ElMessage.warning('没有可保存的主观题')
    return
  }
  savingSubjectiveScores.value = true
  try {
    const res = (await saveAttemptSubjectiveScores(attemptId, items)) as {
      data?: { subjective_score: number; total_score: number }
    }
    if (res?.data && subjectiveGradeDetail.value) {
      subjectiveGradeDetail.value.attempt.subjective_score = res.data.subjective_score
      subjectiveGradeDetail.value.attempt.total_score = res.data.total_score
    }
    if (res?.data && detail.value && detailId.value === attemptId) {
      detail.value.attempt.subjective_score = res.data.subjective_score
      detail.value.attempt.total_score = res.data.total_score
    }
    ElMessage.success('保存成功')
    subjectiveDlgVisible.value = false
    loadList()
    nextTick(() => {
      if (drawer.value && detailId.value === attemptId) {
        startScoreAnimation()
      }
    })
  } catch {
    // 业务错误由 request 拦截器统一 ElMessage
  } finally {
    savingSubjectiveScores.value = false
  }
}

function applyRouteQuery() {
  const q = route.query
  const lv = q.level
  if (typeof lv === 'string' && lv) query.level = lv
  const eid = q.examination_paper_id
  if (eid !== undefined && eid !== '') {
    const n = Number(eid)
    if (!Number.isNaN(n) && n > 0) {
      query.examination_paper_id = n
      paperSel.value = n
    }
  }
  const bid = q.exam_batch_id
  if (bid !== undefined && bid !== '') {
    const n = Number(bid)
    if (!Number.isNaN(n) && n > 0) query.exam_batch_id = n
  }
  const st = q.status
  if (st !== undefined && st !== '') {
    const n = Number(st)
    if (!Number.isNaN(n) && n >= 1 && n <= 4) query.status = n
  }
  if (q.subjective_pending === '1' || q.subjective_pending === 1) {
    query.subjective_pending = 1
  }
  const u = q.username
  if (typeof u === 'string' && u) query.username = u
}

onMounted(async () => {
  await loadPapers()
  applyRouteQuery()
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

.result-list-ops {
  display: inline-flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px 8px;
}
.result-list-op__sub-clip {
  display: inline-block;
  line-height: 1;
}
/* link 型按钮 disabled 时仍易偏主色，强制为禁用色 */
.result-list-ops :deep(.result-list-op--subjective.is-disabled) {
  color: var(--el-text-color-disabled) !important;
  cursor: not-allowed;
}
.result-list-ops :deep(.result-list-op--subjective.is-disabled:hover) {
  color: var(--el-text-color-disabled) !important;
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
.answers-card {
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

.paper-answer-by-section {
  margin-top: 2px;
}

.paper-answer-sec-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 大题分 Tab 查看；仅抽屉主体滚动，避免与内层再套一层竖向滚动条 */
.result-sec-tabs {
  --el-tabs-header-height: auto;
}

.result-sec-tabs :deep(.el-tabs__content) {
  padding: 8px 0 0;
}

.result-sec-tabs :deep(.el-tabs__item) {
  max-width: min(200px, 32vw);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.paper-answer-empty {
  margin-top: 4px;
}

.paper-answer-by-section :deep(.exam-opt-flag) {
  color: var(--result-brand);
}

.subjective-dlg-hint {
  margin: 0 0 6px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
}
.subjective-dlg-meta {
  margin: 0 0 12px;
  font-size: 12px;
  color: var(--el-text-color-regular);
  line-height: 1.5;
  word-break: break-word;
}
.subjective-dlg-empty {
  margin: 10px 0 0;
  font-size: 13px;
  color: var(--el-text-color-placeholder);
}

/* 表内富文本题干：与试卷详情 exam-rich-html 同系 */
.subjective-stem-wrap {
  max-width: 480px;
  line-height: 1.55;
  word-break: break-word;
  font-size: 13px;
  color: var(--el-text-color-primary);
}
.subjective-stem-wrap :deep(.exam-screen-text) {
  white-space: normal;
}
.subjective-stem-plain {
  color: var(--el-text-color-regular);
}
</style>
