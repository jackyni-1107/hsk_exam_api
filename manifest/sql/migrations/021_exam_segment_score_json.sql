-- 记录按 segment_code 的分数字典（整数分）。
ALTER TABLE `exam_attempt`
    ADD COLUMN `segment_score_json` json DEFAULT NULL COMMENT '按 segment_code 的整数分数字典 JSON' AFTER `total_score`;

ALTER TABLE `exam_result`
    ADD COLUMN `segment_score_json` json DEFAULT NULL COMMENT '按 segment_code 的整数分数字典 JSON' AFTER `total_score`;
