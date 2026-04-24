-- 014: 练习/模拟批次策略 + exam_attempt 可重复会话（放宽唯一索引）

-- 1) exam_batch 策略列（默认与现网行为一致）
ALTER TABLE `exam_batch`
    ADD COLUMN `batch_kind` tinyint NOT NULL DEFAULT 0 COMMENT '0=formal 1=practice' AFTER `exam_end_at`,
    ADD COLUMN `allow_multiple_attempts` tinyint NOT NULL DEFAULT 0 COMMENT '1=同用户同卷可多条 exam_attempt' AFTER `batch_kind`,
    ADD COLUMN `max_attempts_per_member` int NOT NULL DEFAULT 0 COMMENT 'allow_multiple=1 时上限，0=不限制' AFTER `allow_multiple_attempts`,
    ADD COLUMN `skip_scoring` tinyint NOT NULL DEFAULT 0 COMMENT '1=finalize 不落 exam_result、不计分' AFTER `max_attempts_per_member`,
    ADD COLUMN `auto_submit_on_deadline` tinyint NOT NULL DEFAULT 1 COMMENT '0=不因个人 deadline 拦截保存/自动交卷，批次过期任务也不代交卷' AFTER `skip_scoring`;

-- 2) exam_attempt：重建 active_unique_flag，使 allow_multiple 批次 (scope=1) 不参与唯一约束
ALTER TABLE `exam_attempt` DROP INDEX `uk_exam_attempt_active_member_batch_paper`;
ALTER TABLE `exam_attempt` DROP COLUMN `active_unique_flag`;

ALTER TABLE `exam_attempt`
    ADD COLUMN `attempt_uniqueness_scope` tinyint NOT NULL DEFAULT 0 COMMENT '0=参与 member+batch+paper 唯一 1=可重复批次' AFTER `mock_level_id`;

ALTER TABLE `exam_attempt`
    ADD COLUMN `active_unique_flag` tinyint GENERATED ALWAYS AS (
        CASE
            WHEN `delete_flag` = 0 AND `exam_batch_id` > 0 AND `attempt_uniqueness_scope` = 0 THEN 1
            ELSE NULL
        END
    ) STORED;

CREATE UNIQUE INDEX `uk_exam_attempt_active_member_batch_paper`
    ON `exam_attempt` (`member_id`, `exam_batch_id`, `exam_paper_id`, `active_unique_flag`);
