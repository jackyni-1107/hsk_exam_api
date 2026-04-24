-- 015: exam_result.status 新增 5=全部算分完成；历史数据与新版 Upsert 语义对齐

-- 仅客观题（或试卷无主观）且已结束的记录
UPDATE `exam_result` r
SET r.`status` = 5, r.`update_time` = NOW()
WHERE r.`delete_flag` = 0
  AND r.`status` = 4
  AND r.`has_subjective` = 0;

-- 含主观题且已有任意主观题人工分
UPDATE `exam_result` r
INNER JOIN (
  SELECT DISTINCT eaa.`attempt_id`
  FROM `exam_attempt_answer` eaa
  INNER JOIN `exam_question` eq ON eq.`id` = eaa.`exam_question_id`
    AND eq.`is_subjective` = 1 AND eq.`is_example` = 0 AND eq.`delete_flag` = 0
  WHERE eaa.`delete_flag` = 0 AND eaa.`awarded_score` IS NOT NULL
) sub ON sub.`attempt_id` = r.`attempt_id`
SET r.`status` = 5, r.`update_time` = NOW()
WHERE r.`delete_flag` = 0
  AND r.`status` = 4
  AND r.`has_subjective` = 1;
