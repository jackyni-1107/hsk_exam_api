<template>
  <div class="exam-screen-text">
    <template v-for="(seg, i) in segments" :key="i">
      <div
        v-if="seg.type === 'text' && seg.content.trim()"
        class="exam-screen-text__text"
        :class="{ 'exam-screen-text__text--html': looksLikeHtml(seg.content) }"
      >
        <div
          v-if="looksLikeHtml(seg.content)"
          class="exam-rich-html"
          v-html="sanitizeHtml(seg.content)"
        />
        <span v-else>{{ seg.content }}</span>
      </div>
      <div
        v-else-if="seg.type === 'image' && seg.content.trim()"
        class="exam-screen-text__img-wrap"
      >
        <img
          v-show="!imgErr[i]"
          class="exam-screen-text__img"
          :src="resolved(seg.content)"
          :alt="seg.content"
          loading="lazy"
          @error="setImgErr(i)"
        />
        <div v-if="imgErr[i]" class="exam-screen-text__fallback">
          <a :href="resolved(seg.content)" target="_blank" rel="noopener noreferrer">{{
            resolved(seg.content)
          }}</a>
        </div>
        <div class="exam-screen-text__path">{{ seg.content }}</div>
      </div>
      <div
        v-else-if="seg.content.trim()"
        class="exam-screen-text__text"
        :class="{ 'exam-screen-text__text--html': looksLikeHtml(seg.content) }"
      >
        <div
          v-if="looksLikeHtml(seg.content)"
          class="exam-rich-html"
          v-html="sanitizeHtml(seg.content)"
        />
        <span v-else>{{ seg.content }}</span>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { looksLikeHtml, sanitizeHtmlForDisplay } from '@/utils/examDisplay'
import {
  mergeAdjacentScreenTextSegments,
  parseScreenTextSegments,
  type ScreenTextSegment,
} from '@/utils/topicJsonDisplay'
import { resolveResourceUrl } from '@/utils/resourceUrl'

const props = withDefaults(
  defineProps<{
    /** topic / DB 中的 screen_text 原始结构 */
    raw: unknown
    /** exam_paper.source_base_url */
    sourceBaseUrl?: string
  }>(),
  {
    sourceBaseUrl: '',
  },
)

const imgErr = reactive<Record<number, boolean>>({})

watch(
  () => props.raw,
  () => {
    for (const k of Object.keys(imgErr)) {
      delete imgErr[Number(k)]
    }
  },
)

const segments = computed((): ScreenTextSegment[] =>
  mergeAdjacentScreenTextSegments(parseScreenTextSegments(props.raw)),
)

function resolved(path: string) {
  return resolveResourceUrl(props.sourceBaseUrl, path)
}

function setImgErr(i: number) {
  imgErr[i] = true
}

function sanitizeHtml(html: string) {
  return sanitizeHtmlForDisplay(html)
}
</script>

<style scoped>
.exam-screen-text {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.exam-screen-text__text {
  white-space: pre-wrap;
  word-break: break-word;
  overflow-wrap: anywhere;
  line-height: 1.55;
}

.exam-screen-text__img-wrap {
  max-width: 100%;
}

.exam-screen-text__img {
  display: block;
  width: 120px;
  height: 120px;
  max-width: 120px;
  max-height: 120px;
  object-fit: cover;
  border-radius: 6px;
  border: 1px solid var(--el-border-color-lighter);
}

.exam-screen-text__fallback {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.exam-screen-text__path {
  margin-top: 4px;
  font-size: 11px;
  color: var(--el-text-color-secondary);
  word-break: break-all;
}

/* 富文本 / 拼音行：不要用 pre-wrap，避免 HTML 源码里标签间换行导致每字断行 */
.exam-screen-text__text--html {
  white-space: normal;
  word-break: normal;
  overflow-wrap: break-word;
}

.exam-screen-text__text--html :deep(img) {
  width: 120px;
  height: 120px;
  max-width: 120px;
  max-height: 120px;
  object-fit: cover;
  border-radius: 6px;
  vertical-align: middle;
}
</style>
