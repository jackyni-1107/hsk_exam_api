/**
 * 试卷题在管理端的展示辅助（与 topic_json 对齐、结果页反查题号坐标）。
 */

import type { ExamPaperDetail } from '@/api/exam'
import {
  blockPassageFromTopicJson,
  effectivePaperQuestionStem,
  getBlockScreenTextRaw,
  getQuestionScreenTextRaw,
  parseScreenTextSegments,
} from '@/utils/topicJsonDisplay'

export type PaperSection = ExamPaperDetail['sections'][number]
export type PaperQuestion = PaperSection['blocks'][number]['questions'][number]

/** 小题在 topic_json 中的坐标（用于 effectivePaperQuestionStem） */
export interface QuestionTopicMeta {
  topicJson: string
  blockIndex: number
  questionIndex: number
  sectionId: number
}

/** 根据试卷详情人 id 反查 question_id → topic 坐标 */
export function buildQuestionTopicMetaById(
  paper: ExamPaperDetail | null | undefined,
): Map<number, QuestionTopicMeta> {
  const m = new Map<number, QuestionTopicMeta>()
  if (!paper?.sections?.length) return m
  for (const sec of paper.sections) {
    const topicJson = sec.topic_json || ''
    const blocks = sec.blocks ?? []
    for (let bi = 0; bi < blocks.length; bi++) {
      const qs = blocks[bi].questions ?? []
      for (let qi = 0; qi < qs.length; qi++) {
        m.set(qs[qi].id, {
          topicJson,
          blockIndex: bi,
          questionIndex: qi,
          sectionId: sec.id,
        })
      }
    }
  }
  return m
}

export { blockPassageFromTopicJson as blockReadingPassageFromTopic, effectivePaperQuestionStem }

export function hasBlockScreenSegments(topicJson: string | undefined, bi: number) {
  return parseScreenTextSegments(getBlockScreenTextRaw(topicJson, bi)).length > 0
}

/** 与试卷预览一致：优先库表 screen_text_json，再 topic_json 内小题的 screen_text。 */
function resolveQuestionScreenTextRaw(
  screenTextJson: string | undefined,
  topicJson: string | undefined,
  bi: number,
  qi: number,
): unknown {
  if (screenTextJson && String(screenTextJson).trim()) return screenTextJson
  return getQuestionScreenTextRaw(topicJson, bi, qi)
}

/**
 * 若答题记录里 stem_text 有值、但题面仍来自 topic/screen 富文本，不能因 stem 非空就退回纯文本；
 * 否则与「试卷管理·详情」用 ExamScreenTextBlocks 的效果不一致。只要题包有 screen_text 片段，即走富文本。
 */
export function hasQuestionStemSegments(
  _stemText: string | undefined,
  screenTextJson: string | undefined,
  topicJson: string | undefined,
  bi: number,
  qi: number,
) {
  const raw = resolveQuestionScreenTextRaw(screenTextJson, topicJson, bi, qi)
  return parseScreenTextSegments(raw).length > 0
}

export function questionStemRichRaw(
  _stemText: string | undefined,
  screenTextJson: string | undefined,
  topicJson: string | undefined,
  bi: number,
  qi: number,
): unknown {
  const raw = resolveQuestionScreenTextRaw(screenTextJson, topicJson, bi, qi)
  if (parseScreenTextSegments(raw).length > 0) return raw
  return undefined
}
