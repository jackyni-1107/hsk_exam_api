<template>
  <div v-if="!sections.length">
    <slot name="empty">
      <el-empty :description="emptyDescription" :image-size="72" />
    </slot>
  </div>
  <div v-else :class="tabsClass">
    <section
      v-for="group in groupedSections"
      :key="group.code"
      class="exam-segment-block"
    >
      <h3 v-if="showSegmentTitle(group.code)" class="exam-segment-title">
        {{ group.label }}（{{ group.questionCount }} 题）
      </h3>
      <section
        v-for="item in group.items"
        :key="item.section.id"
        class="exam-section-block"
      >
        <h4 class="exam-section-title">{{ item.section.title }}</h4>
        <slot name="section-before" :section="item.section" :index="item.index" />
        <div :class="bodyClass">
          <div
            v-for="(card, ci) in item.section.cards"
            :id="questionAnchorId(item, ci)"
            :key="card.key || `card-${item.section.id}-${ci}`"
            class="exam-question-anchor"
          >
            <ExamQuestionReviewCard v-bind="card" />
          </div>
        </div>
        <slot name="section-after" :section="item.section" :index="item.index" />
      </section>
    </section>
    <div v-if="navItemCount > 0" class="exam-global-nav-float">
      <button
        type="button"
        class="exam-global-nav-toggle"
        @click="navPanelOpen = !navPanelOpen"
      >
        {{ navPanelOpen ? "收起题号" : "展开题号" }}
      </button>
      <div v-show="navPanelOpen" class="exam-global-nav">
        <div class="exam-global-nav__title">题号导航</div>
        <div class="exam-global-nav__group-list">
          <section
            v-for="seg in navTree"
            :key="seg.segmentLabel"
            class="exam-global-nav__group"
          >
            <h5 class="exam-global-nav__group-title">{{ seg.segmentLabel }}</h5>
            <section
              v-for="sec in seg.sections"
              :key="`${seg.segmentLabel}-${sec.sectionTitle}`"
              class="exam-global-nav__subgroup"
            >
              <div class="exam-global-nav__subgroup-title">{{ sec.sectionTitle }}</div>
              <div class="exam-global-nav__list">
                <button
                  v-for="item in sec.items"
                  :key="item.anchorId"
                  type="button"
                  class="exam-global-nav__btn"
                  :class="{ 'is-wrong': item.isWrong }"
                  @click="jumpToQuestion(item.anchorId)"
                >
                  <span>{{ item.label }}</span>
                </button>
              </div>
            </section>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import ExamQuestionReviewCard from "@/components/exam/ExamQuestionReviewCard.vue";

type QuestionCardProps = InstanceType<typeof ExamQuestionReviewCard>["$props"] & {
  key?: string;
};

export type ExamQuestionSectionTabItem = {
  id: number | string;
  title: string;
  cards: QuestionCardProps[];
  segmentCode?: string;
};

const props = withDefaults(
  defineProps<{
    sections: ExamQuestionSectionTabItem[];
    tabsClass?: string;
    bodyClass?: string;
    emptyDescription?: string;
  }>(),
  {
    tabsClass: "",
    bodyClass: "paper-answer-sec-body",
    emptyDescription: "暂无题目数据",
  },
);

const sectionsWithName = computed(() =>
  props.sections.map((section, index) => ({
    section,
    index,
    segmentCode: (section.segmentCode || "").trim() || "default",
  })),
);

const groupedSections = computed(() => {
  const groupMap = new Map<
    string,
    {
      code: string;
      label: string;
      items: typeof sectionsWithName.value;
      questionCount: number;
      navItems: { anchorId: string; label: string; isWrong: boolean }[];
    }
  >();
  for (const row of sectionsWithName.value) {
    if (!groupMap.has(row.segmentCode)) {
      groupMap.set(row.segmentCode, {
        code: row.segmentCode,
        label: row.segmentCode === "default" ? "未分组" : row.segmentCode,
        items: [],
        questionCount: 0,
        navItems: [],
      });
    }
    const group = groupMap.get(row.segmentCode)!;
    group.items.push(row);
    row.section.cards.forEach((card, ci) => {
      group.questionCount += 1;
        const wrongFlag =
          (card as { objectiveCorrect?: boolean | null }).objectiveCorrect === false;
      const qNo =
        typeof (card as { questionNo?: number }).questionNo === "number"
          ? (card as { questionNo?: number }).questionNo
          : group.questionCount;
      group.navItems.push({
        anchorId: questionAnchorId(row, ci),
        label: String(qNo),
          isWrong: wrongFlag,
      });
    });
  }
  return [...groupMap.values()];
});

const navTree = computed(() => {
  const segMap = new Map<
    string,
    {
      segmentLabel: string;
      sections: {
        sectionTitle: string;
        items: { anchorId: string; label: string; isWrong: boolean }[];
      }[];
    }
  >();
  for (const g of groupedSections.value) {
    if (!segMap.has(g.label)) {
      segMap.set(g.label, { segmentLabel: g.label, sections: [] });
    }
    const seg = segMap.get(g.label)!;
    for (const row of g.items) {
      const secItems = row.section.cards.map((card, ci) => {
        const qNo =
          typeof (card as { questionNo?: number }).questionNo === "number"
            ? (card as { questionNo?: number }).questionNo
            : ci + 1;
        const isWrong =
          (card as { objectiveCorrect?: boolean | null }).objectiveCorrect === false;
        return {
          anchorId: questionAnchorId(row, ci),
          label: String(qNo),
          isWrong,
        };
      });
      seg.sections.push({
        sectionTitle: row.section.title,
        items: secItems,
      });
    }
  }
  return [...segMap.values()];
});
const navItemCount = computed(() =>
  navTree.value.reduce(
    (sum, seg) =>
      sum + seg.sections.reduce((acc, sec) => acc + sec.items.length, 0),
    0,
  ),
);
const navPanelOpen = ref(true);

function showSegmentTitle(code: string) {
  if (groupedSections.value.length > 1) return true;
  return code !== "default";
}

function questionAnchorId(
  item: { section: ExamQuestionSectionTabItem; index: number },
  cardIndex: number,
) {
  return `q-anchor-${item.section.id}-${item.index}-${cardIndex}`;
}

function jumpToQuestion(anchorId: string) {
  const el = document.getElementById(anchorId);
  if (!el) return;
  el.scrollIntoView({ behavior: "auto", block: "start" });
}
</script>

<style scoped>
.exam-segment-block + .exam-segment-block {
  margin-top: 18px;
}

.exam-segment-title {
  margin: 0 0 10px;
  font-size: 15px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

.exam-global-nav {
  z-index: 31;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 10px;
  padding: 10px;
  max-width: min(360px, 82vw);
  max-height: min(50vh, 420px);
  overflow: auto;
  box-shadow: 0 8px 22px rgba(0, 0, 0, 0.12);
}

.exam-global-nav-float {
  position: fixed;
  right: 20px;
  bottom: 20px;
  z-index: 30;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.exam-global-nav-toggle {
  border: 1px solid var(--el-border-color);
  background: var(--el-color-primary-light-9, #ecf5ff);
  color: var(--el-color-primary);
  font-size: 12px;
  line-height: 1;
  padding: 8px 10px;
  border-radius: 999px;
  cursor: pointer;
}

.exam-global-nav__title {
  font-size: 13px;
  font-weight: 700;
  margin-bottom: 8px;
  color: var(--el-text-color-primary);
}

.exam-global-nav__group-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.exam-global-nav__group-title {
  margin: 0 0 6px;
  font-size: 12px;
  font-weight: 700;
  color: var(--el-text-color-primary);
}

.exam-global-nav__subgroup + .exam-global-nav__subgroup {
  margin-top: 6px;
}

.exam-global-nav__subgroup-title {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}

.exam-global-nav__list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.exam-global-nav__btn {
  border: 1px solid var(--el-border-color);
  background: var(--el-fill-color-blank);
  color: var(--el-text-color-primary);
  font-size: 12px;
  line-height: 1;
  padding: 5px 8px;
  border-radius: 999px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.exam-global-nav__btn:hover {
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
}

.exam-global-nav__btn.is-wrong {
  border-color: color-mix(in srgb, var(--el-color-danger) 60%, var(--el-border-color));
  background: color-mix(in srgb, var(--el-color-danger) 10%, var(--el-fill-color-blank));
  color: var(--el-color-danger);
}

.exam-global-nav__seg {
  color: var(--el-text-color-secondary);
}

.exam-question-anchor {
  scroll-margin-top: 12px;
}

.exam-section-block + .exam-section-block {
  margin-top: 14px;
}

.exam-section-title {
  margin: 0 0 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}
</style>
