-- 017: 考试会话作弊事件记录表（切屏、录屏等）

CREATE TABLE IF NOT EXISTS `exam_attempt_cheat_event` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `attempt_id` bigint NOT NULL COMMENT 'exam_attempt.id',
    `member_id` bigint NOT NULL COMMENT 'sys_member.id',
    `event_type` varchar(64) NOT NULL DEFAULT '' COMMENT '作弊事件类型，如 switch_screen/screen_record',
    `event_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '事件发生时间（服务端写入，与请求到达时刻一致）',
    `segment_code` varchar(32) NOT NULL DEFAULT '' COMMENT '发生时环节编码 listen/read/write',
    `detail` varchar(1024) NOT NULL DEFAULT '' COMMENT '事件详情',
    `client_ip` varchar(64) NOT NULL DEFAULT '' COMMENT '客户端IP',
    `client_agent` varchar(512) NOT NULL DEFAULT '' COMMENT 'User-Agent',
    `creator` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除',
    PRIMARY KEY (`id`),
    KEY `idx_eace_attempt_create` (`attempt_id`, `create_time`),
    KEY `idx_eace_member_create` (`member_id`, `create_time`),
    KEY `idx_eace_event_type` (`event_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='考试会话作弊事件记录';
