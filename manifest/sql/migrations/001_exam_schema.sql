-- HSK 试卷域：试卷、大题、题块、小题、选项（与 internal/model/entity/exam 对齐，供全新库一次性建表）

-- 试卷
CREATE TABLE IF NOT EXISTS `exam_paper` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `level` varchar(32) NOT NULL COMMENT '级别，如 hsk1',
    `paper_id` varchar(64) NOT NULL COMMENT '远程试卷目录ID，如 0d26e5c778ad4ca8',
    `mock_examination_paper_id` bigint NOT NULL COMMENT 'mock 真源 mock_examination_paper.id（业务卷标识）',
    `title` varchar(255) NOT NULL DEFAULT '' COMMENT '试卷标题',
    `prepare_title` varchar(255) NOT NULL DEFAULT '' COMMENT 'prepare.title',
    `prepare_instruction` text COMMENT '考前说明 instruction',
    `prepare_audio_file` varchar(512) NOT NULL DEFAULT '' COMMENT 'prepare.audio_file',
    `source_base_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '资源基址，可拼 index 与媒体 URL',
    `audio_hls_prefix` varchar(512) NOT NULL DEFAULT '' COMMENT 'OSS 桶内 HLS 目录前缀，动态 m3u8 拼接时使用',
    `audio_hls_segment_count` int NOT NULL DEFAULT 0 COMMENT '分片总数，0 表示未配置 HLS',
    `audio_hls_segment_pattern` varchar(64) NOT NULL DEFAULT '' COMMENT '分片文件名 fmt，空则默认 %%05d.ts',
    `audio_hls_key_object` varchar(512) NOT NULL DEFAULT '' COMMENT '密钥对象相对 prefix 的路径，空表示不加密',
    `audio_hls_iv_hex` varchar(32) NOT NULL DEFAULT '' COMMENT 'AES-128 IV 十六进制，写入 #EXT-X-KEY',
    `audio_hls_segment_duration` decimal(10,3) NOT NULL DEFAULT 0.000 COMMENT '#EXTINF 时长秒',
    `index_json` json DEFAULT NULL COMMENT 'index.json 全文快照',
    `duration_seconds` int NOT NULL DEFAULT 0 COMMENT '考试时长秒，0=使用系统默认',
    `creator` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(64) NOT NULL DEFAULT '' COMMENT '更新者',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除：0-否，1-是',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_level_paper` (`level`, `paper_id`),
    KEY `idx_paper_id` (`paper_id`),
    UNIQUE KEY `uk_mock_examination_paper_id` (`mock_examination_paper_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='HSK 试卷';

-- 大题 / 部分（对应 index items）
CREATE TABLE IF NOT EXISTS `exam_section` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `exam_paper_id` bigint NOT NULL COMMENT '试卷ID exam_paper.id',
    `mock_examination_paper_id` bigint NOT NULL COMMENT '冗余 mock_examination_paper.id',
    `sort_order` int NOT NULL DEFAULT 0 COMMENT '在 index.items 中的顺序',
    `topic_title` varchar(255) NOT NULL DEFAULT '' COMMENT 'topic_title',
    `topic_subtitle` varchar(512) NOT NULL DEFAULT '' COMMENT 'topic_subtitle',
    `topic_type` varchar(32) NOT NULL DEFAULT '' COMMENT '题型代码 pt/xp/xt/...',
    `part_code` int NOT NULL DEFAULT 0 COMMENT '大题内 part 序号',
    `segment_code` varchar(32) NOT NULL DEFAULT '' COMMENT 'listen/read',
    `topic_items_file` varchar(255) NOT NULL DEFAULT '' COMMENT 'topic_items 文件名，如 pt.json',
    `topic_json` json DEFAULT NULL COMMENT '该 topic 文件全文快照',
    `creator` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(64) NOT NULL DEFAULT '' COMMENT '更新者',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除',
    PRIMARY KEY (`id`),
    KEY `idx_exam_paper_sort` (`exam_paper_id`, `sort_order`),
    KEY `idx_mock_exam_paper` (`mock_examination_paper_id`),
    KEY `idx_topic_type` (`topic_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='HSK 试卷大题';

-- 题块（扁平：一题一块；套题：多题共享一块）
CREATE TABLE IF NOT EXISTS `exam_question_block` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `section_id` bigint NOT NULL COMMENT '大题ID exam_section.id',
    `block_order` int NOT NULL DEFAULT 0 COMMENT '对应 topic JSON 中 items 下标',
    `group_index` int DEFAULT NULL COMMENT '套题外层 index（若存在）',
    `question_description_json` json DEFAULT NULL COMMENT '块级 question_description_obj 等',
    `creator` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(64) NOT NULL DEFAULT '' COMMENT '更新者',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除',
    PRIMARY KEY (`id`),
    KEY `idx_section_block` (`section_id`, `block_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='HSK 题块';

-- 小题
CREATE TABLE IF NOT EXISTS `exam_question` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `exam_paper_id` bigint NOT NULL COMMENT '试卷ID，冗余便于按卷查询',
    `mock_examination_paper_id` bigint NOT NULL COMMENT '冗余 mock_examination_paper.id',
    `block_id` bigint NOT NULL COMMENT '题块ID exam_question_block.id',
    `sort_in_block` int NOT NULL DEFAULT 0 COMMENT '块内顺序',
    `question_no` int DEFAULT NULL COMMENT '卷面题号（JSON index，如 1-40）',
    `score` decimal(10,2) DEFAULT NULL COMMENT '分值',
    `is_example` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否例题',
    `is_subjective` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否主观题：0否 1是',
    `content_type` varchar(32) NOT NULL DEFAULT '' COMMENT '题干内容类型，如 audio',
    `audio_file` varchar(512) NOT NULL DEFAULT '' COMMENT '音频 content 文件名',
    `audio_hls_prefix` varchar(512) NOT NULL DEFAULT '' COMMENT 'OSS 桶内 HLS 目录前缀（无首尾/）',
    `audio_hls_segment_count` int NOT NULL DEFAULT 0 COMMENT '分片总数，0 表示未配置 HLS',
    `audio_hls_segment_pattern` varchar(64) NOT NULL DEFAULT '' COMMENT '分片文件名 fmt，空则默认 %%05d.ts',
    `audio_hls_key_object` varchar(512) NOT NULL DEFAULT '' COMMENT '密钥对象相对 prefix 的路径，空表示不加密',
    `audio_hls_iv_hex` varchar(32) NOT NULL DEFAULT '' COMMENT 'AES-128 IV 十六进制',
    `audio_hls_segment_duration` decimal(10,3) NOT NULL DEFAULT 0.000 COMMENT '#EXTINF 时长秒',
    `stem_text` text COMMENT 'content_sentence',
    `screen_text_json` json DEFAULT NULL COMMENT 'screen_text 数组',
    `analysis_json` json DEFAULT NULL COMMENT 'analysis 多语言',
    `question_description_json` json DEFAULT NULL COMMENT '小题级 question_description_obj',
    `raw_json` json DEFAULT NULL COMMENT '单题原始 JSON 备份',
    `creator` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(64) NOT NULL DEFAULT '' COMMENT '更新者',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除',
    PRIMARY KEY (`id`),
    KEY `idx_block_sort` (`block_id`, `sort_in_block`),
    KEY `idx_paper_no` (`exam_paper_id`, `question_no`),
    KEY `idx_mock_exam_paper` (`mock_examination_paper_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='HSK 试题';

-- 选项
CREATE TABLE IF NOT EXISTS `exam_option` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `question_id` bigint NOT NULL COMMENT '试题ID exam_question.id',
    `flag` varchar(16) NOT NULL DEFAULT '' COMMENT '选项标识 A/B/C/T/F',
    `sort_order` int NOT NULL DEFAULT 0 COMMENT '对应 answers.index',
    `is_correct` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否正确',
    `option_type` varchar(32) NOT NULL DEFAULT '' COMMENT 'text/image/pinyin 等',
    `content` varchar(2048) NOT NULL DEFAULT '' COMMENT '文本或资源文件名',
    `creator` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(64) NOT NULL DEFAULT '' COMMENT '更新者',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除',
    PRIMARY KEY (`id`),
    KEY `idx_question_sort` (`question_id`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='HSK 试题选项';
