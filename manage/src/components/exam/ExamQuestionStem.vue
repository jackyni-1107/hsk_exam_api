<template>
  <div
    class="exam-q-stem"
    :class="{ 'exam-q-stem--html': useHtml }"
  >
    <div v-if="useHtml" class="exam-rich-html" v-html="safeHtml" />
    <template v-else>{{ text }}</template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { looksLikeHtml, sanitizeHtmlForDisplay, stemDisplayText } from '@/utils/examDisplay'

const props = withDefaults(
  defineProps<{
    stemText?: string
    screenTextJson?: string
    emptyFallback?: string
  }>(),
  {
    stemText: '',
    screenTextJson: '',
    emptyFallback: undefined,
  },
)

const text = computed(() =>
  stemDisplayText(props.stemText, props.screenTextJson, props.emptyFallback),
)

const useHtml = computed(() => looksLikeHtml(text.value))

const safeHtml = computed(() => sanitizeHtmlForDisplay(text.value))
</script>

<style scoped>
.exam-q-stem {
  word-break: break-word;
  overflow-wrap: anywhere;
  line-height: 1.55;
}

.exam-q-stem--html :deep(font),
.exam-q-stem--html :deep(span) {
  vertical-align: baseline;
}
</style>
