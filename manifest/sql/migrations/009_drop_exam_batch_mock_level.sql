-- 009: 移除已废弃的 exam_batch_mock_level（批次可选卷后续见 010 exam_batch_paper）
-- 须在 008_exam_batch_multi_mock_paper.sql 之后执行。

DROP TABLE IF EXISTS `exam_batch_mock_level`;
