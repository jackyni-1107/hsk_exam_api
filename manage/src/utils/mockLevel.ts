import type { MockLevelItem } from '@/api/mockAdmin'

/** 侧栏/筛选用：type_name + 空格 + level_name */
export function mockLevelOptionLabel(lv: MockLevelItem): string {
  const a = (lv.type_name || '').trim()
  const b = (lv.level_name || '').trim()
  if (a && b) return `${a} ${b}`
  return a || b || `id:${lv.id}`
}
