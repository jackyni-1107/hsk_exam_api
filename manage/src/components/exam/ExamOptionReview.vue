<template>
  <div
    class="exam-opt-review"
    :class="{ 'exam-opt-review--strip': layoutStrip }"
  >
    <button
      v-for="(o, i) in options"
      :key="optionKey(o, i)"
      type="button"
      class="exam-opt-card"
      :class="cardClass(o)"
      disabled
    >
      <span class="exam-opt-card__flag">{{ o.flag || '—' }}</span>
      <template v-if="isImageOption(o) && resolvedUrl(o)">
        <img
          class="exam-opt-card__img"
          :src="resolvedUrl(o)"
          :alt="o.flag || ''"
          loading="lazy"
        />
      </template>
      <span
        v-else-if="htmlText(o)"
        class="exam-opt-card__text exam-opt-card__html exam-rich-html"
        v-html="htmlText(o)"
      />
      <span v-else class="exam-opt-card__text">{{ textContent(o) }}</span>
      <el-tag
        v-if="showCorrect && o.is_correct === 1"
        size="small"
        type="success"
        effect="plain"
        class="exam-opt-card__tag"
        >标答</el-tag
      >
      <el-tag
        v-if="isUserSelected(o) && mode === 'review'"
        size="small"
        type="primary"
        effect="dark"
        class="exam-opt-card__tag"
        >所选</el-tag
      >
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  isImageOptionRow,
  looksLikeHtml,
  optionContentLabel,
  sanitizeHtmlForDisplay,
  type ExamOptionDisplayRow,
} from '@/utils/examDisplay'
import { resolveResourceUrl } from '@/utils/resourceUrl'

const props = withDefaults(
  defineProps<{
    options: ExamOptionDisplayRow[]
    sourceBaseUrl?: string
    /** preview：无「所选」；review：根据 selectedOptionIds 高亮 */
    mode?: 'preview' | 'review'
    showCorrect?: boolean
    /** 用户选的 option id（可多选） */
    selectedOptionIds?: number[]
  }>(),
  {
    sourceBaseUrl: '',
    mode: 'preview',
    showCorrect: true,
    selectedOptionIds: () => [],
  },
)

const layoutStrip = computed(() => {
  const opts = props.options ?? []
  if (opts.length < 2) return false
  return opts.every((o) => isImageOption(o))
})

function isImageOption(o: ExamOptionDisplayRow) {
  return isImageOptionRow(o)
}

function htmlText(o: ExamOptionDisplayRow) {
  const c = (o.content || '').trim()
  if (!c || isImageOptionRow(o)) return ''
  if (!looksLikeHtml(c)) return ''
  return sanitizeHtmlForDisplay(c)
}

function resolvedUrl(o: ExamOptionDisplayRow) {
  return resolveResourceUrl(props.sourceBaseUrl, o.content)
}

function textContent(o: ExamOptionDisplayRow) {
  return optionContentLabel(o)
}

function optionKey(o: ExamOptionDisplayRow, i: number) {
  if (o.id != null && o.id !== 0) return `o-${o.id}`
  return `o-${o.flag ?? i}-${i}`
}

function isUserSelected(o: ExamOptionDisplayRow) {
  if (!o.id) return false
  return (props.selectedOptionIds ?? []).includes(o.id)
}

function cardClass(o: ExamOptionDisplayRow) {
  const cls: Record<string, boolean> = {
    'is-correct': props.showCorrect && o.is_correct === 1,
    'is-user': props.mode === 'review' && isUserSelected(o),
  }
  if (props.mode === 'review' && !props.showCorrect) {
    /* keep */
  }
  if (props.mode === 'review' && o.is_correct === 1 && isUserSelected(o)) {
    cls['is-match'] = true
  }
  if (props.mode === 'review' && isUserSelected(o) && o.is_correct !== 1) {
    cls['is-wrong'] = true
  }
  return cls
}
</script>

<style scoped>
.exam-opt-review {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.exam-opt-review--strip {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
}

@media (max-width: 480px) {
  .exam-opt-review--strip {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

.exam-opt-card {
  flex: 1 1 50px;
  min-width: 0;
  min-height: auto;
  height: auto;
  max-width: 100%;
  margin: 0;
  padding: 10px;
  text-align: left;
  border-radius: 10px;
  border: 1px solid var(--el-border-color);
  background: var(--el-fill-color-blank);
  cursor: default;
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  align-items: center;
  gap: 6px;
  transition:
    border-color 0.15s,
    box-shadow 0.15s;
}

.exam-opt-review--strip .exam-opt-card {
  width: 100%;
  max-width: 100%;
}

.exam-opt-card:disabled {
  opacity: 1;
  color: inherit;
}

.exam-opt-card__flag {
  font-weight: 800;
  font-size: 14px;
  color: var(--el-color-primary);
}

.exam-opt-card__img {
  width: 72px;
  height: 72px;
  max-width: 72px;
  max-height: 72px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid var(--el-border-color-lighter);
  flex-shrink: 0;
}

.exam-opt-review--strip .exam-opt-card__img {
  width: 72px;
  height: 72px;
  max-width: 72px;
  max-height: 72px;
}

.exam-opt-card__text {
  flex: 1;
  min-width: 0;
  font-size: 13px;
  line-height: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.exam-opt-card__html :deep(font),
.exam-opt-card__html :deep(span) {
  vertical-align: baseline;
}

.exam-opt-card__html :deep(img) {
  width: 72px !important;
  height: 72px !important;
  max-width: 72px !important;
  max-height: 72px !important;
  object-fit: cover !important;
  border-radius: 6px;
  vertical-align: middle;
}

.exam-opt-card__tag {
  align-self: center;
  flex-shrink: 0;
}

.exam-opt-card.is-correct {
  border-color: color-mix(in srgb, var(--el-color-success) 45%, var(--el-border-color));
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--el-color-success) 25%, transparent);
}

.exam-opt-card.is-user.is-wrong {
  border-color: color-mix(in srgb, var(--el-color-danger) 50%, var(--el-border-color));
}

.exam-opt-card.is-match {
  border-color: color-mix(in srgb, var(--el-color-success) 60%, var(--el-border-color));
}
</style>
