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

export function hasQuestionStemSegments(
  stemText: string | undefined,
  screenTextJson: string | undefined,
  topicJson: string | undefined,
  bi: number,
  qi: number,
) {
  if ((stemText || '').trim()) return false
  const raw =
    screenTextJson && String(screenTextJson).trim()
      ? screenTextJson
      : getQuestionScreenTextRaw(topicJson, bi, qi)
  return parseScreenTextSegments(raw).length > 0
}

export function questionStemRichRaw(
  stemText: string | undefined,
  screenTextJson: string | undefined,
  topicJson: string | undefined,
  bi: number,
  qi: number,
): unknown {
  if ((stemText || '').trim()) return undefined
  return screenTextJson && String(screenTextJson).trim()
    ? screenTextJson
    : getQuestionScreenTextRaw(topicJson, bi, qi)
}
