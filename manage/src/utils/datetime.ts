function pad2(n: number): string {
  return String(n).padStart(2, '0')
}

export function formatUtcText(value?: string | null): string {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())} ${pad2(d.getHours())}:${pad2(d.getMinutes())}:${pad2(d.getSeconds())}`
}

export function formatUtcForDisplay(_row: unknown, _column: unknown, cellValue?: string): string {
  return formatUtcText(cellValue)
}
