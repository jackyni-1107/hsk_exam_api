-- 可选：本地/联调填充 exam_question.id=1 的 HLS 元数据（需在业务数据中存在 id=1 的小题）
-- HLS 分片命名 segment_000.ts … segment_180.ts（共 181 个）、AES-128 等说明见原注释。
--
-- 桶内对象（前缀 audio_hls_prefix，无首尾 /）示例：
--   dev/hls/exam_question_1/segment_000.ts … segment_180.ts
--   dev/hls/exam_question_1/encryption.key
--
-- audio_hls_iv_hex：32 个十六进制字符（无 0x），须与转码一致；示例为全 0，上线前请改为真实 IV。
UPDATE `exam_question`
SET
  `audio_hls_prefix` = 'exam/question_1',
  `audio_hls_segment_count` = 10,
  `audio_hls_segment_pattern` = 'segment_%03d.ts',
  `audio_hls_key_object` = 'encryption.key',
  `audio_hls_iv_hex` = '00000000000000000000000000000000',
  `audio_hls_segment_duration` = 10.000,
  `updater` = 'admin',
  `update_time` = NOW()
WHERE `id` = 1 AND `delete_flag` = 0;
