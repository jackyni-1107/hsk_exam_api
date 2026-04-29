-- 019: Extend notification template for system/third-party compatibility.
-- Goals:
-- 1) Add template type and third-party fields.
-- 2) Bind each template to a notification channel config.
-- 3) Provide backfill strategy for historical rows.

ALTER TABLE `sys_notification_template`
    ADD COLUMN `channel_config_id` bigint NOT NULL DEFAULT 0 COMMENT '绑定的通知渠道配置ID（sys_notification_channel_config.id）' AFTER `channel`,
    ADD COLUMN `template_type` tinyint NOT NULL DEFAULT 1 COMMENT '模板类型：1=系统模板 2=第三方模板' AFTER `channel_config_id`,
    ADD COLUMN `third_party_template_id` varchar(128) NOT NULL DEFAULT '' COMMENT '第三方模板ID' AFTER `content`,
    ADD COLUMN `third_party_template_params` varchar(2000) NOT NULL DEFAULT '' COMMENT '第三方模板参数(JSON)' AFTER `third_party_template_id`;

-- Indexes for template query and channel-config binding lookups.
ALTER TABLE `sys_notification_template`
    ADD INDEX `idx_sys_notif_tpl_channel_cfg_del` (`channel_config_id`, `delete_flag`),
    ADD INDEX `idx_sys_notif_tpl_type_status` (`template_type`, `status`),
    ADD INDEX `idx_sys_notif_tpl_channel_type` (`channel`, `template_type`);

-- Backfill strategy (existing data):
-- A. Bind historical templates to active channel config by same channel.
UPDATE `sys_notification_template` t
INNER JOIN `sys_notification_channel_config` c
        ON c.`channel` = t.`channel`
       AND c.`is_active` = 1
       AND c.`delete_flag` = 0
SET t.`channel_config_id` = c.`id`
WHERE t.`delete_flag` = 0
  AND t.`channel_config_id` = 0;

-- B. If no active config exists, fallback to latest available config per channel.
UPDATE `sys_notification_template` t
INNER JOIN (
    SELECT c1.`channel`, MAX(c1.`id`) AS `id`
    FROM `sys_notification_channel_config` c1
    WHERE c1.`delete_flag` = 0
    GROUP BY c1.`channel`
) latest ON latest.`channel` = t.`channel`
SET t.`channel_config_id` = latest.`id`
WHERE t.`delete_flag` = 0
  AND t.`channel_config_id` = 0;

-- C. Backfill template_type:
--    historical records default to system templates unless third-party id already exists.
UPDATE `sys_notification_template`
SET `template_type` = CASE
    WHEN TRIM(IFNULL(`third_party_template_id`, '')) <> '' THEN 2
    ELSE 1
END
WHERE `delete_flag` = 0;

-- D. Validation queries (execute manually after migration):
--    1) Remaining unbound templates should be reviewed and manually bound in admin UI.
--       SELECT id, code, channel FROM sys_notification_template WHERE delete_flag = 0 AND channel_config_id = 0;
--    2) Ensure third-party templates have template IDs.
--       SELECT id, code FROM sys_notification_template WHERE delete_flag = 0 AND template_type = 2 AND TRIM(third_party_template_id) = '';
