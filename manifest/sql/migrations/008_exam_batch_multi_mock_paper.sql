-- 008: 批次多 Mock 卷（exam_batch_mock_paper）；成员主键含 mock_examination_paper_id；exam_batch 去掉单列 mock 卷
-- 须在 005_exam_batch.sql、007_exam_batch_attempt_member_level.sql 之后执行。
-- 历史数据：同一批次同一学员因多 mock_level 多行、但卷相同时，保留 mock_level_id 最小的一行。

CREATE TABLE IF NOT EXISTS `exam_batch_mock_paper` (
    `batch_id` bigint NOT NULL COMMENT 'exam_batch.id',
    `mock_examination_paper_id` bigint NOT NULL COMMENT 'mock_examination_paper.id',
    PRIMARY KEY (`batch_id`, `mock_examination_paper_id`),
    KEY `idx_ebmp_mock_paper` (`mock_examination_paper_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='批次可选 Mock 卷（多选）';

INSERT IGNORE INTO `exam_batch_mock_paper` (`batch_id`, `mock_examination_paper_id`)
SELECT `id`, `mock_examination_paper_id` FROM `exam_batch`;

ALTER TABLE `exam_batch_member`
    ADD COLUMN `mock_examination_paper_id` bigint NOT NULL DEFAULT 0 COMMENT 'mock_examination_paper.id' AFTER `member_id`;

UPDATE `exam_batch_member` ebm
    INNER JOIN `exam_batch` eb ON eb.`id` = ebm.`batch_id`
    SET ebm.`mock_examination_paper_id` = eb.`mock_examination_paper_id`;

DELETE e1 FROM `exam_batch_member` e1
    INNER JOIN `exam_batch_member` e2
        ON e1.`batch_id` = e2.`batch_id` AND e1.`member_id` = e2.`member_id`
        AND e1.`mock_examination_paper_id` = e2.`mock_examination_paper_id`
        AND e1.`mock_level_id` > e2.`mock_level_id`;

ALTER TABLE `exam_batch_member` DROP PRIMARY KEY;
ALTER TABLE `exam_batch_member` DROP INDEX `idx_ebm_batch_level`;

ALTER TABLE `exam_batch_member` DROP COLUMN `mock_level_id`;

ALTER TABLE `exam_batch_member`
    ADD PRIMARY KEY (`batch_id`, `member_id`, `mock_examination_paper_id`),
    ADD KEY `idx_ebm_batch_paper` (`batch_id`, `mock_examination_paper_id`);

ALTER TABLE `exam_batch` DROP INDEX `idx_exam_batch_mock_paper`;
ALTER TABLE `exam_batch` DROP COLUMN `mock_examination_paper_id`;
