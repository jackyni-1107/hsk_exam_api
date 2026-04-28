-- 018: Seed dictionary data for config center groups.
-- The config center reads group options from dict type `sys_config_group`.

INSERT INTO sys_dict_type (
  dict_name, dict_type, status, remark, creator, delete_flag
)
SELECT
  '配置中心分组',
  'sys_config_group',
  0,
  '配置中心分组下拉选项',
  'migration',
  0
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1
  FROM sys_dict_type
  WHERE dict_type = 'sys_config_group'
    AND delete_flag = 0
);

INSERT INTO sys_dict_data (
  dict_type, dict_label, dict_value, sort, status, remark, creator, delete_flag
)
SELECT
  'sys_config_group',
  '默认分组',
  'default',
  0,
  0,
  '系统默认配置分组',
  'migration',
  0
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1
  FROM sys_dict_data
  WHERE dict_type = 'sys_config_group'
    AND dict_value = 'default'
    AND delete_flag = 0
);
