-- 007: 批次成员带等级主键；exam_attempt / exam_result 冗余 exam_batch_id、mock_level_id

-- 1) exam_batch_member：增加 mock_level_id，回填后改主键为 (batch_id, member_id, mock_level_id)
ALTER TABLE `exam_batch_member`
    ADD COLUMN `mock_level_id` bigint NOT NULL DEFAULT 0 COMMENT 'mock_levels.id' AFTER `member_id`;

-- 按批次在 exam_batch_mock_level 中的最小等级回填（无配置行则保持 0，须运维修正）
UPDATE `exam_batch_member` ebm
    INNER JOIN (
        SELECT `batch_id`, MIN(`mock_level_id`) AS `mid`
        FROM `exam_batch_mock_level`
        GROUP BY `batch_id`
    ) t ON t.`batch_id` = ebm.`batch_id`
    SET ebm.`mock_level_id` = t.`mid`
    WHERE ebm.`mock_level_id` = 0;

ALTER TABLE `exam_batch_member` DROP PRIMARY KEY;
ALTER TABLE `exam_batch_member`
    ADD PRIMARY KEY (`batch_id`, `member_id`, `mock_level_id`),
    ADD KEY `idx_ebm_batch_level` (`batch_id`, `mock_level_id`);

-- 2) exam_attempt
ALTER TABLE `exam_attempt`
    ADD COLUMN `exam_batch_id` bigint NOT NULL DEFAULT 0 COMMENT 'exam_batch.id，0=历史非批次' AFTER `mock_examination_paper_id`,
    ADD COLUMN `mock_level_id` bigint NOT NULL DEFAULT 0 COMMENT 'mock_levels.id，0=历史非批次' AFTER `exam_batch_id`,
    ADD KEY `idx_exam_attempt_batch` (`exam_batch_id`),
    ADD KEY `idx_exam_attempt_member_batch_level` (`member_id`, `exam_batch_id`, `mock_level_id`);

-- 3) exam_result
ALTER TABLE `exam_result`
    ADD COLUMN `exam_batch_id` bigint NOT NULL DEFAULT 0 COMMENT '冗余 exam_batch.id' AFTER `mock_examination_paper_id`,
    ADD COLUMN `mock_level_id` bigint NOT NULL DEFAULT 0 COMMENT '冗余 mock_levels.id' AFTER `exam_batch_id`,
    ADD KEY `idx_exam_result_batch` (`exam_batch_id`);
