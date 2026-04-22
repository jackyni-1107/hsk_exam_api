<template>
  <section class="exam-section-panel">
    <header v-if="title || subtitle" class="exam-section-panel__head">
      <h3 v-if="title" class="exam-section-panel__title">{{ title }}</h3>
      <p v-if="subtitle" class="exam-section-panel__sub">{{ subtitle }}</p>
    </header>
    <div class="exam-section-panel__body">
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
  subtitle?: string
}>()
</script>

<style scoped>
.exam-section-panel {
  margin-bottom: 20px;
  border-radius: 12px;
  border: 1px solid var(--el-border-color-lighter);
  background: var(--el-fill-color-blank);
  overflow: hidden;
}

.exam-section-panel__head {
  padding: 14px 18px;
  background: linear-gradient(180deg, var(--el-fill-color-light) 0%, transparent 100%);
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.exam-section-panel__title {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

.exam-section-panel__sub {
  margin: 6px 0 0;
  font-size: 13px;
  color: var(--el-text-color-secondary);
  line-height: 1.45;
}

.exam-section-panel__body {
  padding: 16px 18px 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}
</style>
