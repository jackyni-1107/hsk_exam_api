/**
 * 批次等接口要求 RFC3339（含时区偏移或 Z）。日期选择器产出 `YYYY-MM-DD HH:mm:ss` 墙钟时间时，
 * 按「展示时区」解析后再提交 UTC 的 ISO 字符串（…Z）。
 *
 * 环境变量：`VITE_DISPLAY_TIMEZONE`（IANA，如 `Asia/Shanghai`）。未设置则用浏览器本机时区。
 */

function readEnvDisplayTimeZone(): string | undefined {
  try {
    const v = (import.meta as unknown as { env?: Record<string, string | undefined> })
      .env?.VITE_DISPLAY_TIMEZONE;
    const s = typeof v === "string" ? v.trim() : "";
    return s || undefined;
  } catch {
    return undefined;
  }
}

/** 业务上「管理员所选墙钟时间」对应的 IANA 时区 */
export function getDisplayTimeZone(): string {
  return (
    readEnvDisplayTimeZone() ||
    Intl.DateTimeFormat().resolvedOptions().timeZone ||
    "UTC"
  );
}

function wallPartsInTimeZone(tMs: number, timeZone: string) {
  const str = new Intl.DateTimeFormat("sv-SE", {
    timeZone,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
    hourCycle: "h23",
  }).format(new Date(tMs));
  const [datePart, timePart = "00:00:00"] = str.split(" ");
  const [y, mo, d] = datePart.split("-").map(Number);
  const [h, mi, s] = timePart.split(":").map(Number);
  return { y, mo, d, h, mi, s };
}

function wallEqual(
  a: ReturnType<typeof wallPartsInTimeZone>,
  b: ReturnType<typeof wallPartsInTimeZone>,
) {
  return (
    a.y === b.y &&
    a.mo === b.mo &&
    a.d === b.d &&
    a.h === b.h &&
    a.mi === b.mi &&
    a.s === b.s
  );
}

/**
 * 将 `YYYY-MM-DD HH:mm:ss`（或 `T` 分隔）视为 `timeZone` 时区下的墙钟时刻，返回对应 UTC 的 ISO 字符串（RFC3339 Z）。
 */
export function wallClockStringToRFC3339UTC(
  wall: string,
  timeZone = getDisplayTimeZone(),
): string {
  const m = /^(\d{4})-(\d{2})-(\d{2})[ T](\d{2}):(\d{2})(?::(\d{2}))?$/.exec(
    wall.trim(),
  );
  if (!m) {
    throw new Error(`无效时间格式：${wall}`);
  }
  const y = +m[1];
  const mo = +m[2];
  const d = +m[3];
  const h = +m[4];
  const mi = +m[5];
  const s = m[6] != null ? +m[6] : 0;
  const target = { y, mo, d, h, mi, s };
  const center = Date.UTC(y, mo - 1, d, h, mi, s);
  for (let delta = -40 * 3600000; delta <= 40 * 3600000; delta += 60000) {
    const t = center + delta;
    if (wallEqual(wallPartsInTimeZone(t, timeZone), target)) {
      return new Date(t).toISOString();
    }
  }
  throw new Error(
    `无法在时区「${timeZone}」内解析该本地时间（可能落在夏令时跳空内），请调整时刻或修改 VITE_DISPLAY_TIMEZONE`,
  );
}

/** 将接口返回的 ISO/RFC3339 转为与日期选择器一致的墙钟串（展示时区） */
export function isoOrRfcToWallClockForPicker(
  iso?: string | null,
  timeZone = getDisplayTimeZone(),
): string {
  if (!iso) return "";
  const t = new Date(iso).getTime();
  if (Number.isNaN(t)) return String(iso);
  const str = new Intl.DateTimeFormat("sv-SE", {
    timeZone,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
    hourCycle: "h23",
  }).format(new Date(t));
  return str.replace("T", " ");
}

function formatTimeZoneSuffix(isoOrWall: string, timeZone: string): string {
  const d = new Date(isoOrWall);
  if (Number.isNaN(d.getTime())) return "";
  const parts = new Intl.DateTimeFormat("en-GB", {
    timeZone,
    timeZoneName: "shortOffset",
  }).formatToParts(d);
  return parts.find((p) => p.type === "timeZoneName")?.value?.trim() || timeZone;
}

/** 表格/详情：墙钟 + 时区后缀 */
export function formatUtcText(value?: string | null): string {
  if (!value) return "";
  const t = new Date(value).getTime();
  if (Number.isNaN(t)) return value;
  const tz = getDisplayTimeZone();
  return new Intl.DateTimeFormat("sv-SE", {
    timeZone: tz,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
    hourCycle: "h23",
  }).format(new Date(t));
  // const suf = formatTimeZoneSuffix(value, tz);
  // return wall;
}

export function formatUtcForDisplay(
  _row: unknown,
  _column: unknown,
  cellValue?: string,
): string {
  return formatUtcText(cellValue);
}
