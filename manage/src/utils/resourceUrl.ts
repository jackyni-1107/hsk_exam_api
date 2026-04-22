/**
 * 试卷资源基址（exam_paper.source_base_url）与题目包内相对路径拼接。
 */

export function normalizeSourceBaseUrl(base: string | undefined | null): string {
  const b = (base || '').trim()
  if (!b) return ''
  return b.replace(/\/+$/, '') + '/'
}

/** 将相对路径拼到 source_base_url 上；已是 http(s)/data/blob 则原样返回 */
export function resolveResourceUrl(
  sourceBaseUrl: string | undefined | null,
  path: string | undefined | null,
): string {
  const p = (path || '').trim()
  if (!p) return ''
  if (/^(https?:|data:|blob:|\/\/)/i.test(p)) return p
  const base = normalizeSourceBaseUrl(sourceBaseUrl)
  if (!base) return p
  const rel = p.startsWith('/') ? p.slice(1) : p
  return base + rel
}
