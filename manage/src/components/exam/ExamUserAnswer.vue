<template>
  <div class="exam-user-answer">{{ formatted }}</div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { formatAnswerJson } from '@/utils/examDisplay'

const props = defineProps<{
  answerJson: string
}>()

const formatted = computed(() => formatAnswerJson(props.answerJson ?? ''))
</script>

<style scoped>
.exam-user-answer {
  font-family: ui-monospace, 'Cascadia Code', 'Consolas', monospace;
  font-size: 12px;
  word-break: break-word;
  overflow-wrap: anywhere;
  line-height: 1.55;
}
</style>
