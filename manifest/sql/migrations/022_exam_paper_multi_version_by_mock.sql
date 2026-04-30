-- 021: 支持同一 mock_examination_paper_id 对应多份 exam_paper（用于 conflict_mode=new 新建新卷）。
-- 将 exam_paper.mock_examination_paper_id 从唯一约束改为普通索引：
-- - fail: 仍按“存在任一未删除记录即冲突”
-- - overwrite: 由调用方指定 overwrite_exam_paper_id 仅覆盖一份
-- - new: 允许插入新行，生成新的 exam_paper.id

ALTER TABLE `exam_paper`
    DROP INDEX `uk_mock_examination_paper_id`;

ALTER TABLE `exam_paper`
DROP INDEX `uk_level_paper`;

ALTER TABLE `exam_paper`
    ADD KEY `idx_mock_examination_paper_id` (`mock_examination_paper_id`);
