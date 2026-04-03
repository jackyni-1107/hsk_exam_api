-- 005: 考试批次（时间窗、多选 mock_levels、批次内学员 sys_member）

CREATE TABLE IF NOT EXISTS `exam_batch` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mock_examination_paper_id` bigint NOT NULL COMMENT 'mock 卷 id，与 exam_paper 业务主键一致',
    `title` varchar(255) NOT NULL DEFAULT '' COMMENT '批次名称',
    `exam_start_at` datetime NOT NULL COMMENT '考试允许开始时间',
    `exam_end_at` datetime NOT NULL COMMENT '考试允许结束时间',
    `creator` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(64) NOT NULL DEFAULT '' COMMENT '更新者',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除：0-否，1-是',
    PRIMARY KEY (`id`),
    KEY `idx_exam_batch_mock_paper` (`mock_examination_paper_id`),
    KEY `idx_exam_batch_time` (`exam_start_at`, `exam_end_at`),
    KEY `idx_exam_batch_delete` (`delete_flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='考试批次';

CREATE TABLE IF NOT EXISTS `exam_batch_mock_level` (
    `batch_id` bigint NOT NULL COMMENT 'exam_batch.id',
    `mock_level_id` bigint NOT NULL COMMENT 'mock_levels.id',
    PRIMARY KEY (`batch_id`, `mock_level_id`),
    KEY `idx_ebml_level` (`mock_level_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='批次可选等级（多选）';

CREATE TABLE IF NOT EXISTS `exam_batch_member` (
    `batch_id` bigint NOT NULL COMMENT 'exam_batch.id',
    `member_id` bigint NOT NULL COMMENT 'sys_member.id',
    `creator` varchar(64) NOT NULL DEFAULT '' COMMENT '导入操作者',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '导入时间',
    PRIMARY KEY (`batch_id`, `member_id`),
    KEY `idx_ebm_member` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='批次参考学员';
