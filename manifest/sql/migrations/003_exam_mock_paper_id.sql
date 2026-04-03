-- 003: 试卷以 mock_examination_paper.id 为业务主键；子表冗余 mock 卷 id；exam_paper.mock_examination_paper_id NOT NULL + 唯一。
-- 前置条件：每条 exam_paper.mock_examination_paper_id 已正确回填且非 NULL；覆盖导入会物理删卷故可复用同一 mock id。

-- 1) 子表冗余列
ALTER TABLE `exam_section`
    ADD COLUMN `mock_examination_paper_id` bigint NOT NULL DEFAULT 0 COMMENT '冗余 mock_examination_paper.id' AFTER `exam_paper_id`;

ALTER TABLE `exam_question`
    ADD COLUMN `mock_examination_paper_id` bigint NOT NULL DEFAULT 0 COMMENT '冗余 mock_examination_paper.id' AFTER `exam_paper_id`;

ALTER TABLE `exam_attempt`
    ADD COLUMN `mock_examination_paper_id` bigint NOT NULL DEFAULT 0 COMMENT '冗余 mock_examination_paper.id' AFTER `exam_paper_id`;

ALTER TABLE `exam_result`
    ADD COLUMN `mock_examination_paper_id` bigint NOT NULL DEFAULT 0 COMMENT '冗余 mock_examination_paper.id' AFTER `exam_paper_id`;

-- 2) 从 exam_paper 回填冗余列
UPDATE `exam_section` s
    INNER JOIN `exam_paper` p ON p.id = s.exam_paper_id
    SET s.mock_examination_paper_id = p.mock_examination_paper_id
    WHERE p.mock_examination_paper_id IS NOT NULL;

UPDATE `exam_question` q
    INNER JOIN `exam_paper` p ON p.id = q.exam_paper_id
    SET q.mock_examination_paper_id = p.mock_examination_paper_id
    WHERE p.mock_examination_paper_id IS NOT NULL;

UPDATE `exam_attempt` a
    INNER JOIN `exam_paper` p ON p.id = a.exam_paper_id
    SET a.mock_examination_paper_id = p.mock_examination_paper_id
    WHERE p.mock_examination_paper_id IS NOT NULL;

UPDATE `exam_result` r
    INNER JOIN `exam_paper` p ON p.id = r.exam_paper_id
    SET r.mock_examination_paper_id = p.mock_examination_paper_id
    WHERE p.mock_examination_paper_id IS NOT NULL;

-- 3) exam_paper：mock 列 NOT NULL + 唯一（每 mock 卷对应一行 exam_paper）
ALTER TABLE `exam_paper` DROP INDEX `idx_mock_examination_paper`;

ALTER TABLE `exam_paper`
    MODIFY COLUMN `mock_examination_paper_id` bigint NOT NULL COMMENT 'mock 真源 mock_examination_paper.id';

ALTER TABLE `exam_paper`
    ADD UNIQUE KEY `uk_mock_examination_paper_id` (`mock_examination_paper_id`);

-- 4) 子表索引
ALTER TABLE `exam_section` ADD KEY `idx_mock_exam_paper` (`mock_examination_paper_id`);
ALTER TABLE `exam_question` ADD KEY `idx_mock_exam_paper` (`mock_examination_paper_id`);
ALTER TABLE `exam_attempt` ADD KEY `idx_mock_exam_paper` (`mock_examination_paper_id`);
ALTER TABLE `exam_result` ADD KEY `idx_mock_exam_paper` (`mock_examination_paper_id`);
