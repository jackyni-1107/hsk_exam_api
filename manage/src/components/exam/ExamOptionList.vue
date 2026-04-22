<template>
  <ul class="exam-opt-list">
    <li
      v-for="(o, i) in options"
      :key="optionKey(o, i)"
      class="exam-opt-row"
      :class="{ 'exam-opt-row--correct': showCorrect && o.is_correct === 1 }"
    >
      <span class="exam-opt-flag">{{ o.flag || '—' }}</span>
      <span class="exam-opt-content">
        <template v-if="isImageOption(o) && resolvedUrl(o)">
          <img
            class="exam-opt-img"
            :src="resolvedUrl(o)"
            :alt="optionCaption(o)"
            loading="lazy"
          />
          <span class="exam-opt-caption">{{ optionContentLabel(o) }}</span>
        </template>
        <template v-else>{{ optionContentLabel(o) }}</template>
      </span>
      <el-tag
        v-if="showCorrect && o.is_correct === 1"
        size="small"
        type="success"
        effect="plain"
        class="exam-opt-tag"
        >标答</el-tag
      >
    </li>
  </ul>
</template>

<script setup lang="ts">
import { optionContentLabel, type ExamOptionDisplayRow } from '@/utils/examDisplay'
import { resolveResourceUrl } from '@/utils/resourceUrl'

const props = withDefaults(
  defineProps<{
    options: ExamOptionDisplayRow[]
    /** 是否高亮正确答案并显示「标答」标签 */
    showCorrect?: boolean
    /** exam_paper.source_base_url，用于拼接选项中的相对图片路径 */
    sourceBaseUrl?: string
  }>(),
  {
    showCorrect: true,
    sourceBaseUrl: '',
  },
)

function optionKey(o: ExamOptionDisplayRow, i: number) {
  if (o.id != null && o.id !== 0) return o.id
  return `${o.flag ?? ''}-${i}`
}

function isImageOption(o: ExamOptionDisplayRow) {
  const c = (o.content || '').trim()
  const t = (o.option_type || '').toLowerCase()
  return t === 'image' || /\.(jpe?g|png|gif|webp|svg)(\?|$)/i.test(c)
}

function resolvedUrl(o: ExamOptionDisplayRow) {
  return resolveResourceUrl(props.sourceBaseUrl, o.content)
}

function optionCaption(o: ExamOptionDisplayRow) {
  return optionContentLabel(o)
}
</script>

<style scoped>
.exam-opt-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.exam-opt-row {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 6px 8px;
  margin-bottom: 8px;
  padding: 6px 8px;
  border-radius: 8px;
  background: color-mix(in srgb, var(--el-fill-color) 92%, transparent);
}

.exam-opt-row--correct {
  font-weight: 500;
}

.exam-opt-content {
  flex: 1 1 140px;
  min-width: 0;
  word-break: break-word;
  overflow-wrap: anywhere;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.exam-opt-img {
  max-width: min(100%, 360px);
  height: auto;
  border-radius: 6px;
  border: 1px solid var(--el-border-color-lighter);
}

.exam-opt-caption {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.exam-opt-tag {
  flex-shrink: 0;
}

.exam-opt-flag {
  font-weight: 700;
  flex-shrink: 0;
  color: var(--exam-brand, #2c5282);
}
</style>
