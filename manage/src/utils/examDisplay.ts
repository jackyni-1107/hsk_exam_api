/** 管理端试卷/考试结果共用的展示逻辑（题干、选项、考生答案 JSON） */

export interface ExamOptionDisplayRow {
  id?: number
  flag?: string
  content?: string
  is_correct?: number
  option_type?: string
  sort_order?: number
}

const DEFAULT_EMPTY_STEM =
  '（无文字题干；听力/图片等题型可能仅有资源文件名，见下方选项或试卷资源）'

/**
 * 题干展示：优先 stem_text；空串时可附带 screen_text_json 的简要提示。
 */
export function stemDisplayText(
  stemText: string | undefined,
  screenTextJson?: string,
  emptyFallback: string = DEFAULT_EMPTY_STEM,
): string {
  const s = (stemText || '').trim()
  if (s) return s
  const raw = (screenTextJson || '').trim()
  if (raw) {
    return `（无 stem_text；screen_text_json 长度 ${raw.length} 字符）`
  }
  return emptyFallback
}

/** 是否按 HTML 渲染（题干/选项富文本） */
export function looksLikeHtml(s: string | undefined): boolean {
  const t = (s || '').trim()
  if (t.length < 3) return false
  return /<\/?[a-z][\s\S]*>/i.test(t)
}

/** 管理端展示用：去掉 script，降低 v-html 风险 */
export function sanitizeHtmlForDisplay(html: string): string {
  return html.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
}

/** 选项是否为图片资源（与 ExamOptionReview 一致） */
export function isImageOptionRow(o: ExamOptionDisplayRow): boolean {
  const c = (o.content || '').trim()
  const t = (o.option_type || '').toLowerCase()
  return t === 'image' || /\.(jpe?g|png|gif|webp|svg)(\?|$)/i.test(c)
}

/** 选项内容展示：图片/音频等加前缀提示 */
export function optionContentLabel(o: ExamOptionDisplayRow): string {
  const c = (o.content || '').trim()
  if (!c) return '—'
  const t = (o.option_type || '').toLowerCase()
  if (t === 'image' || /\.(jpe?g|png|gif|webp|svg)(\?|$)/i.test(c)) {
    return `[图片] ${c}`
  }
  if (t === 'audio' || /\.(mp3|wav|m4a|ogg)(\?|$)/i.test(c)) {
    return `[音频] ${c}`
  }
  return c
}

/** 解析答题 JSON（客观题 option_id / 多选 / 主观 text） */
export function parseAnswerPayload(raw: string | undefined | null): {
  optionId?: number
  optionIds?: number[]
  text?: string
} {
  if (!raw) return {}
  try {
    const o = JSON.parse(raw) as Record<string, unknown>
    const out: { optionId?: number; optionIds?: number[]; text?: string } = {}
    if (typeof o.option_id === 'number') out.optionId = o.option_id
    /** 客户端常见字段，与 option_id 等价 */
    if (out.optionId == null) {
      const oid = o.o_id
      if (typeof oid === 'number' && !Number.isNaN(oid)) out.optionId = oid
      else if (typeof oid === 'string' && /^\d+$/.test(oid)) out.optionId = Number(oid)
    }
    const sel = o.selected_option_ids
    if (Array.isArray(sel)) {
      const ids = sel.filter((x): x is number => typeof x === 'number')
      if (ids.length) out.optionIds = ids
    }
    if (typeof o.text === 'string') out.text = o.text
    return out
  } catch {
    return {}
  }
}

/** 考生答案 JSON 可读化 */
export function formatAnswerJson(raw: string): string {
  if (!raw) return '—'
  const p = parseAnswerPayload(raw)
  if (p.optionId != null) return `选项 ID: ${p.optionId}`
  if (p.optionIds && p.optionIds.length > 0) {
    return `选项 ID: ${p.optionIds.map((x) => String(x)).join('、')}`
  }
  if (typeof p.text === 'string')
    return p.text.length > 80 ? `${p.text.slice(0, 80)}…` : p.text
  try {
    const o = JSON.parse(raw) as Record<string, unknown>
    return JSON.stringify(o)
  } catch {
    return raw.length > 100 ? `${raw.slice(0, 100)}…` : raw
  }
}

/** 解析类文本：试卷 analysis_json 可能为 JSON 字符串 */
export function analysisDisplayText(raw: string | undefined): string {
  if (!raw) return ''
  const t = raw.trim()
  if (!t) return ''
  try {
    const o = JSON.parse(t) as unknown
    if (typeof o === 'string') return o
    return JSON.stringify(o)
  } catch {
    return t
  }
}
