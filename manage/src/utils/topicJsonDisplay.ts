/**
 * 试卷 section.topic_json 展示辅助（套题如阅读理解：共用材料在 item.screen_text，小题题干可能在 questions[].screen_text）
 */

export interface TopicJsonParseResult {
  items?: TopicJsonItem[]
}

export interface TopicJsonItem {
  screen_text?: unknown
  questions?: TopicJsonQuestion[]
  [key: string]: unknown
}

export interface TopicJsonQuestion {
  screen_text?: unknown
  [key: string]: unknown
}

/** 解析 section 级 topic_json */
export function parseTopicJson(topicJson: string | undefined | null): TopicJsonParseResult | null {
  if (!topicJson || !String(topicJson).trim()) return null
  try {
    const o = JSON.parse(topicJson) as TopicJsonParseResult
    if (o && Array.isArray(o.items)) return o
    return null
  } catch {
    return null
  }
}

/**
 * 将 topic 中的 screen_text 转为纯文本。
 * 支持：数组 [{ type, content }]，或 JSON 字符串形式的数组，或普通字符串。
 */
export function screenTextsToPlain(screenText: unknown): string {
  if (screenText == null) return ''
  if (typeof screenText === 'string') {
    const t = screenText.trim()
    if (!t) return ''
    if (t.startsWith('[') || t.startsWith('{')) {
      try {
        return screenTextsToPlain(JSON.parse(t) as unknown)
      } catch {
        return t
      }
    }
    return t
  }
  if (Array.isArray(screenText)) {
    const parts: string[] = []
    for (const seg of screenText) {
      if (seg && typeof seg === 'object' && seg !== null && 'content' in seg) {
        const c = (seg as { content?: unknown }).content
        if (typeof c === 'string' && c.trim()) parts.push(c.trim())
      }
    }
    return parts.join('\n')
  }
  return ''
}

/** 某套题（与 block 同序的 item）共用阅读材料 */
export function blockPassageFromTopicJson(
  topicJson: string | undefined | null,
  blockIndex: number,
): string {
  const p = parseTopicJson(topicJson)
  const item = p?.items?.[blockIndex] as TopicJsonItem | undefined
  return screenTextsToPlain(item?.screen_text)
}

/** 小题在 topic_json 中的题干（items[bi].questions[qi].screen_text） */
export function questionStemFromTopicJson(
  topicJson: string | undefined | null,
  blockIndex: number,
  questionIndex: number,
): string {
  const p = parseTopicJson(topicJson)
  const q = p?.items?.[blockIndex]?.questions?.[questionIndex] as TopicJsonQuestion | undefined
  return screenTextsToPlain(q?.screen_text)
}

/**
 * 管理端展示用题干：优先 DB stem_text，其次 DB screen_text_json（数组 JSON），最后 topic_json 中对应小题的 screen_text。
 */
export function effectivePaperQuestionStem(
  stemText: string | undefined,
  screenTextJson: string | undefined,
  topicJson: string | undefined | null,
  blockIndex: number,
  questionIndex: number,
): string {
  if ((stemText || '').trim()) return (stemText || '').trim()
  const fromDbScreen = screenTextsToPlain(screenTextJson)
  if (fromDbScreen) return fromDbScreen
  return questionStemFromTopicJson(topicJson, blockIndex, questionIndex)
}

export interface ScreenTextSegment {
  type: string
  content: string
}

/**
 * 将 screen_text 解析为片段（text / image 等），用于带图展示。
 */
export function parseScreenTextSegments(screenText: unknown): ScreenTextSegment[] {
  if (screenText == null) return []
  let arr: unknown = screenText
  if (typeof screenText === 'string') {
    const t = screenText.trim()
    if (!t) return []
    if (t.startsWith('[')) {
      try {
        arr = JSON.parse(t) as unknown
      } catch {
        return [{ type: 'text', content: t }]
      }
    } else {
      return [{ type: 'text', content: t }]
    }
  }
  if (!Array.isArray(arr)) return []
  const out: ScreenTextSegment[] = []
  for (const seg of arr) {
    if (!seg || typeof seg !== 'object') continue
    const o = seg as { type?: string; content?: string }
    const typ = String(o.type || 'text').toLowerCase()
    const content = typeof o.content === 'string' ? o.content : ''
    if (!content.trim() && typ !== 'image') continue
    out.push({ type: typ, content })
  }
  return out
}

/**
 * 合并相邻的 text 片段为一条，避免管理端每个片段单独包一层块级 div、在 flex 列里竖成「一列字」。
 * （题包常见把每个字/词拆成数组里多条 type:text。）
 */
export function mergeAdjacentScreenTextSegments(segments: ScreenTextSegment[]): ScreenTextSegment[] {
  const out: ScreenTextSegment[] = []
  for (const seg of segments) {
    const typ = String(seg.type || 'text').toLowerCase()
    if (typ !== 'text') {
      out.push({ ...seg })
      continue
    }
    const prev = out[out.length - 1]
    if (prev && String(prev.type || 'text').toLowerCase() === 'text') {
      prev.content = prev.content + seg.content
    } else {
      out.push({ type: seg.type, content: seg.content })
    }
  }
  return out
}

/** topic_json 中套题共用 screen_text 原始值 */
export function getBlockScreenTextRaw(
  topicJson: string | undefined | null,
  blockIndex: number,
): unknown {
  const p = parseTopicJson(topicJson)
  const item = p?.items?.[blockIndex] as TopicJsonItem | undefined
  return item?.screen_text
}

/** topic_json 中小题 screen_text 原始值 */
export function getQuestionScreenTextRaw(
  topicJson: string | undefined | null,
  blockIndex: number,
  questionIndex: number,
): unknown {
  const p = parseTopicJson(topicJson)
  const q = p?.items?.[blockIndex]?.questions?.[questionIndex] as TopicJsonQuestion | undefined
  return q?.screen_text
}
