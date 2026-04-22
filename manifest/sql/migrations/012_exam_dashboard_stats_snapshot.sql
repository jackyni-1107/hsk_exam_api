-- 012: Dashboard 考试监控全量预聚合快照 + 定时任务种子

CREATE TABLE IF NOT EXISTS exam_dashboard_stats_snapshot (
  scope_key   VARCHAR(64) NOT NULL PRIMARY KEY COMMENT 'global=全量',
  payload     JSON         NOT NULL COMMENT 'bo.AttemptAdminStatsSnapshotPayload JSON',
  update_time DATETIME(3)  NOT NULL,
  KEY idx_exam_dash_stats_ut (update_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='考试监控统计快照（由定时任务写入）';

INSERT INTO sys_task (
  name, code, type, cron_expr, delay_seconds, handler, params,
  retry_times, retry_interval, concurrency, alert_on_fail, alert_receivers,
  status, remark, creator, delete_flag, create_time, update_time
)
SELECT
  '考试监控统计刷新',
  'exam_dashboard_stats_refresh',
  1,
  '0 */2 * * * *',
  0,
  'ExamDashboardStatsHandler',
  '{}',
  0, 0, 0, 0, '',
  0,
  '预聚合 Dashboard 考试指标（全量 scope_key=global），可改 cron/停用',
  'migration',
  0,
  NOW(3), NOW(3)
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM sys_task t WHERE t.code = 'exam_dashboard_stats_refresh' AND t.delete_flag = 0
);
