-- 010: 移除未再使用的 exam_attempt_question_audio（听力播放进度等已不再落库）
-- 若库中无此表，DROP IF EXISTS 可安全执行。

DROP TABLE IF EXISTS `exam_attempt_question_audio`;
