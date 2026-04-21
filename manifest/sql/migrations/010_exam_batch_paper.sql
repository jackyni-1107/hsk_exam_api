-- 010: 批次关联卷改为 exam_paper（exam_batch_paper），同时冗余 mock_examination_paper_id；成员主键含 exam_paper_id，并冗余 mock
-- 须在 008_exam_batch_multi_mock_paper.sql 之后执行。
-- 说明：exam_paper.mock_examination_paper_id 有唯一约束，可与旧 mock id 一一对应。

-- 1) 新建 exam_batch_paper（exam_paper + mock 冗余）
CREATE TABLE IF NOT EXISTS `exam_batch_paper` (
    `batch_id` bigint NOT NULL COMMENT 'exam_batch.id',
    `exam_paper_id` bigint NOT NULL COMMENT 'exam_paper.id',
    `mock_examination_paper_id` bigint NOT NULL COMMENT 'mock_examination_paper.id，与 exam_paper 同步',
    PRIMARY KEY (`batch_id`, `exam_paper_id`),
    KEY `idx_ebp_paper` (`exam_paper_id`),
    KEY `idx_ebp_mock` (`mock_examination_paper_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='批次可选试卷（exam_paper 多选）';

INSERT IGNORE INTO `exam_batch_paper` (`batch_id`, `exam_paper_id`, `mock_examination_paper_id`)
SELECT ebp.`batch_id`, MIN(ep.`id`) AS `exam_paper_id`, ebp.`mock_examination_paper_id`
FROM `exam_batch_mock_paper` ebp
INNER JOIN `exam_paper` ep ON ep.`mock_examination_paper_id` = ebp.`mock_examination_paper_id`
    AND ep.`delete_flag` = 0
GROUP BY ebp.`batch_id`, ebp.`mock_examination_paper_id`;

DROP TABLE IF EXISTS `exam_batch_mock_paper`;

-- 2) exam_batch_member：补充 exam_paper_id；保留 mock_examination_paper_id；主键改为含 exam_paper_id
ALTER TABLE `exam_batch_member`
    ADD COLUMN `exam_paper_id` bigint NOT NULL DEFAULT 0 COMMENT 'exam_paper.id' AFTER `member_id`;

UPDATE `exam_batch_member` ebm
INNER JOIN `exam_paper` ep ON ep.`mock_examination_paper_id` = ebm.`mock_examination_paper_id`
    AND ep.`delete_flag` = 0
SET ebm.`exam_paper_id` = ep.`id`;

DELETE FROM `exam_batch_member` WHERE `exam_paper_id` = 0;

ALTER TABLE `exam_batch_member` DROP PRIMARY KEY;
ALTER TABLE `exam_batch_member` DROP INDEX `idx_ebm_batch_paper`;

ALTER TABLE `exam_batch_member`
    ADD PRIMARY KEY (`batch_id`, `member_id`, `exam_paper_id`),
    ADD KEY `idx_ebm_batch_paper` (`batch_id`, `exam_paper_id`),
    ADD KEY `idx_ebm_mock` (`mock_examination_paper_id`);
