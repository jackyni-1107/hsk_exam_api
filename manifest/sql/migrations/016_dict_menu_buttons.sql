-- 016: 字典管理按钮权限。
-- 背景：
-- 1) 字典接口已拆分为独立权限：
--    dict:type:list/create/update/delete
--    dict:data:list/create/update/delete
-- 2) 现网字典页面父菜单通常仍为 path = '/system/dict' / permission = 'dict:list'，
--    且缺少对应按钮节点；不补齐时，页面可见但接口可能返回 403。
--
-- 处理：
-- 1) 在字典页面父菜单下补齐 8 个按钮节点。
-- 2) 将当前已拥有字典页面菜单的角色，自动补齐这 8 个按钮权限。

SET @dict_parent_id = (
  SELECT id
  FROM sys_menu
  WHERE delete_flag = 0
    AND (path = '/system/dict' OR permission = 'dict:list')
  ORDER BY id ASC
  LIMIT 1
);

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '字典类型列表',
  'dict:type:list',
  3,
  10,
  @dict_parent_id,
  '', '', '', '',
  0, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @dict_parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'dict:type:list' AND delete_flag = 0) = 0;

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '字典类型创建',
  'dict:type:create',
  3,
  20,
  @dict_parent_id,
  '', '', '', '',
  0, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @dict_parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'dict:type:create' AND delete_flag = 0) = 0;

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '字典类型更新',
  'dict:type:update',
  3,
  30,
  @dict_parent_id,
  '', '', '', '',
  0, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @dict_parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'dict:type:update' AND delete_flag = 0) = 0;

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '字典类型删除',
  'dict:type:delete',
  3,
  40,
  @dict_parent_id,
  '', '', '', '',
  0, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @dict_parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'dict:type:delete' AND delete_flag = 0) = 0;

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '字典数据列表',
  'dict:data:list',
  3,
  50,
  @dict_parent_id,
  '', '', '', '',
  0, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @dict_parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'dict:data:list' AND delete_flag = 0) = 0;

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '字典数据创建',
  'dict:data:create',
  3,
  60,
  @dict_parent_id,
  '', '', '', '',
  0, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @dict_parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'dict:data:create' AND delete_flag = 0) = 0;

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '字典数据更新',
  'dict:data:update',
  3,
  70,
  @dict_parent_id,
  '', '', '', '',
  0, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @dict_parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'dict:data:update' AND delete_flag = 0) = 0;

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '字典数据删除',
  'dict:data:delete',
  3,
  80,
  @dict_parent_id,
  '', '', '', '',
  0, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @dict_parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'dict:data:delete' AND delete_flag = 0) = 0;

INSERT INTO sys_role_menu (role_id, menu_id, creator, delete_flag)
SELECT DISTINCT
  rm.role_id,
  btn.id,
  'migration',
  0
FROM sys_role_menu rm
JOIN sys_menu parent
  ON parent.id = rm.menu_id
 AND parent.delete_flag = 0
JOIN sys_menu btn
  ON btn.parent_id = parent.id
 AND btn.delete_flag = 0
 AND btn.permission IN (
   'dict:type:list',
   'dict:type:create',
   'dict:type:update',
   'dict:type:delete',
   'dict:data:list',
   'dict:data:create',
   'dict:data:update',
   'dict:data:delete'
 )
WHERE rm.delete_flag = 0
  AND (parent.path = '/system/dict' OR parent.permission = 'dict:list')
  AND NOT EXISTS (
    SELECT 1
    FROM sys_role_menu existing
    WHERE existing.role_id = rm.role_id
      AND existing.menu_id = btn.id
      AND existing.delete_flag = 0
  );
