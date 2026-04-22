<template>
  <article
    class="exam-q-card"
    :class="{
      'exam-q-card--wrong': isWrongObjective,
      'exam-q-card--ok': isOkObjective,
    }"
  >
    <div class="exam-q-card__head">
      <span class="exam-q-badge">{{ questionNo }}</span>
      <div class="exam-q-card__meta">
        <span v-if="isExample === 1" class="exam-q-pill">例题</span>
        <span v-else-if="isSubjective === 1" class="exam-q-pill exam-q-pill--sub">主观</span>
        <span v-else class="exam-q-pill exam-q-pill--obj">客观</span>
        <span class="exam-q-score">{{ scoreDisplay }}</span>
        <template v-if="mode === 'review' && !isExample && !isSubjective">
          <el-icon v-if="objectiveCorrect === true" class="exam-q-ic exam-q-ic--ok"><Check /></el-icon>
          <el-icon v-else-if="objectiveCorrect === false" class="exam-q-ic exam-q-ic--bad"><Close /></el-icon>
          <el-icon v-else class="exam-q-ic exam-q-ic--muted"><Minus /></el-icon>
        </template>
      </div>
    </div>

    <div v-if="audioFile" class="exam-q-audio">
      音频：
      <a :href="audioHref" target="_blank" rel="noopener noreferrer">{{ audioFile }}</a>
    </div>

    <div v-if="showBlockPassageBlock" class="exam-q-passage">
      <div class="exam-q-label">阅读材料（套题共用）</div>
      <ExamScreenTextBlocks
        v-if="blockHasRich"
        :raw="blockScreenRaw"
        :source-base-url="sourceBaseUrl"
      />
      <div v-else class="exam-q-passage-plain">{{ blockPassageText }}</div>
    </div>

    <template v-if="showStemSection">
      <div class="exam-q-label">题干</div>
      <ExamScreenTextBlocks
        v-if="stemRich"
        :raw="stemRichRaw"
        :source-base-url="sourceBaseUrl"
      />
      <ExamQuestionStem v-else :stem-text="stemPlain" />
    </template>

    <template v-if="showYourAnswerBlock">
      <div class="exam-q-label">你的答案</div>
      <ExamUserAnswer :answer-json="answerJson || ''" />
    </template>

    <template v-if="options?.length">
      <div class="exam-q-label">选项</div>
      <ExamOptionReview
        :options="options"
        :source-base-url="sourceBaseUrl"
        :mode="mode"
        :show-correct="showCorrectOptions"
        :selected-option-ids="resolvedSelectedIds"
      />
    </template>

    <div v-if="analysisText" class="exam-q-analysis">
      <div class="exam-q-label">解析</div>
      <div
        v-if="looksLikeHtml(analysisText)"
        class="exam-q-analysis-body exam-q-analysis-body--html exam-rich-html"
        v-html="sanitizeHtmlForDisplay(analysisText)"
      />
      <div v-else class="exam-q-analysis-body">{{ analysisText }}</div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Check, Close, Minus } from '@element-plus/icons-vue'
import ExamQuestionStem from '@/components/exam/ExamQuestionStem.vue'
import ExamScreenTextBlocks from '@/components/exam/ExamScreenTextBlocks.vue'
import ExamUserAnswer from '@/components/exam/ExamUserAnswer.vue'
import ExamOptionReview from '@/components/exam/ExamOptionReview.vue'
import {
  looksLikeHtml,
  parseAnswerPayload,
  sanitizeHtmlForDisplay,
  type ExamOptionDisplayRow,
} from '@/utils/examDisplay'
import { resolveResourceUrl } from '@/utils/resourceUrl'
import {
  effectivePaperQuestionStem,
  hasBlockScreenSegments,
  hasQuestionStemSegments,
  questionStemRichRaw,
} from '@/utils/examPaperQuestionDisplay'
import { getBlockScreenTextRaw } from '@/utils/topicJsonDisplay'

const props = withDefaults(
  defineProps<{
    mode?: 'preview' | 'review'
    questionNo: number
    score: number
    isExample?: number
    isSubjective?: number
    stemText?: string
    screenTextJson?: string
    topicJson?: string
    blockIndex?: number
    questionIndex?: number
    /** 套题：当前题块下的共用阅读段 */
    blockPassageText?: string
    /** 同一题块内仅首题传 true，避免阅读材料重复 */
    showBlockPassage?: boolean
    sourceBaseUrl?: string
    audioFile?: string
    /** 选项 */
    options?: ExamOptionDisplayRow[]
    showCorrectOptions?: boolean
    /** 批阅 */
    answerJson?: string
    objectiveCorrect?: boolean | null
    analysisText?: string
    /** 管理端结果：主观题已评得分，与 score（满分）同时展示 */
    awardedScore?: number | null
  }>(),
  {
    mode: 'preview',
    isExample: 0,
    isSubjective: 0,
    stemText: '',
    screenTextJson: '',
    topicJson: '',
    blockIndex: 0,
    questionIndex: 0,
    blockPassageText: '',
    showBlockPassage: true,
    sourceBaseUrl: '',
    audioFile: '',
    options: () => [],
    showCorrectOptions: true,
    answerJson: '',
    objectiveCorrect: null,
    analysisText: '',
    awardedScore: undefined,
  },
)

const scoreDisplay = computed(() => {
  const max = Number(props.score) || 0
  if (
    props.mode === 'review' &&
    props.isSubjective === 1 &&
    props.isExample !== 1
  ) {
    const aw = props.awardedScore
    if (aw != null && !Number.isNaN(Number(aw))) {
      return `${Number(aw).toFixed(2)} / ${max.toFixed(2)} 分`
    }
    return `未评 · 满分 ${max.toFixed(2)}`
  }
  return `${max} 分`
})

const stemRich = computed(() => {
  return hasQuestionStemSegments(
    props.stemText,
    props.screenTextJson,
    props.topicJson,
    props.blockIndex ?? 0,
    props.questionIndex ?? 0,
  )
})

const stemRichRaw = computed(() => {
  return questionStemRichRaw(
    props.stemText,
    props.screenTextJson,
    props.topicJson,
    props.blockIndex ?? 0,
    props.questionIndex ?? 0,
  )
})

const stemPlain = computed(() =>
  effectivePaperQuestionStem(
    props.stemText,
    props.screenTextJson,
    props.topicJson,
    props.blockIndex ?? 0,
    props.questionIndex ?? 0,
  ),
)

const blockScreenRaw = computed(() =>
  getBlockScreenTextRaw(props.topicJson, props.blockIndex ?? 0),
)

const blockHasRich = computed(() =>
  hasBlockScreenSegments(props.topicJson, props.blockIndex ?? 0),
)

const showBlockPassageBlock = computed(() => {
  if (props.showBlockPassage === false) return false
  return !!(props.blockPassageText?.trim() || blockHasRich.value)
})

const audioHref = computed(() =>
  resolveResourceUrl(props.sourceBaseUrl, props.audioFile),
)

const isWrongObjective = computed(
  () =>
    props.mode === 'review' &&
    !props.isExample &&
    !props.isSubjective &&
    props.objectiveCorrect === false,
)

const isOkObjective = computed(
  () =>
    props.mode === 'review' &&
    !props.isExample &&
    !props.isSubjective &&
    props.objectiveCorrect === true,
)

const resolvedSelectedIds = computed(() => {
  const raw = parseAnswerPayload(props.answerJson)
  const ids: number[] = []
  if (raw.optionId != null) ids.push(raw.optionId)
  if (raw.optionIds?.length) ids.push(...raw.optionIds)
  return ids
})

const showStemSection = computed(() => {
  if (stemRich.value) return true
  return !!(stemPlain.value || '').trim()
})

/** 仅主观题展示文字作答；客观题以选项区「所选」标签为准 */
const showYourAnswerBlock = computed(
  () => props.mode === 'review' && props.isSubjective === 1,
)
</script>

<style scoped>
.exam-q-card {
  padding: 14px 16px;
  border-radius: 10px;
  border: 1px solid var(--el-border-color-lighter);
  background: var(--el-bg-color);
}

.exam-q-card--wrong {
  border-color: color-mix(in srgb, var(--el-color-danger) 35%, var(--el-border-color-lighter));
  background: color-mix(in srgb, var(--el-color-danger) 6%, var(--el-bg-color));
}

.exam-q-card--ok {
  border-color: color-mix(in srgb, var(--el-color-success) 28%, var(--el-border-color-lighter));
}

.exam-q-card__head {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.exam-q-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  border-radius: 50%;
  font-weight: 800;
  font-size: 14px;
  color: #fff;
  background: var(--el-color-primary);
  flex-shrink: 0;
}

.exam-q-card__meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.exam-q-pill {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 999px;
  background: var(--el-fill-color-dark);
  color: var(--el-text-color-secondary);
}

.exam-q-pill--sub {
  background: color-mix(in srgb, var(--el-color-warning) 22%, transparent);
}

.exam-q-pill--obj {
  background: color-mix(in srgb, var(--el-color-primary) 16%, transparent);
}

.exam-q-score {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.exam-q-ic {
  font-size: 20px;
  margin-left: auto;
}
.exam-q-ic--ok {
  color: var(--el-color-success);
}
.exam-q-ic--bad {
  color: var(--el-color-danger);
}
.exam-q-ic--muted {
  color: var(--el-text-color-placeholder);
}

.exam-q-label {
  display: block;
  font-size: 11px;
  font-weight: 600;
  color: var(--el-text-color-secondary);
  margin: 12px 0 6px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.exam-q-label:first-of-type {
  margin-top: 0;
}

.exam-q-audio {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}

.exam-q-passage {
  margin-bottom: 8px;
  padding: 10px;
  border-radius: 8px;
  background: var(--el-fill-color-light);
  border: 1px dashed var(--el-border-color);
}

.exam-q-passage-plain {
  white-space: pre-wrap;
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
}

.exam-q-analysis-body {
  font-size: 13px;
  line-height: 1.55;
  word-break: break-word;
}

.exam-q-analysis-body--html :deep(font),
.exam-q-analysis-body--html :deep(span) {
  vertical-align: baseline;
}
</style>
