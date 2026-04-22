/**
 * 管理端站点展示配置（侧栏系统名、页脚版权、浏览器标题等）。
 * 部署不同品牌时只需改本文件常量。
 */
export const siteConfig = {
  /** 侧栏顶部展示的系统名称 */
  systemName: 'HSK 考试管理平台',
  /**
   * 底部版权文案（可含 ©、年份、公司名；写任意固定字符串即可）
   * 以下为首次加载时的年份，长会话跨年时如需更新可改为页面挂载时拼接
   */
  copyright: `© ${new Date().getFullYear()} HSK 考试管理平台`,
  /** 浏览器标签页标题 */
  documentTitle: 'HSK 考试管理平台',
} as const
