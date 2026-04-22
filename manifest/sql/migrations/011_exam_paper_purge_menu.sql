-- 011: 试卷物理删除按钮权限 exam:paper:purge，仅授予 sys_role.code = 'super_admin' 的角色。
-- 父菜单：取 permission = exam:list 的第一条（试卷列表）；若库中无 exam:list，请先补全菜单或手工指定 parent_id 后执行 INSERT。

SET @parent_id = (SELECT id FROM sys_menu WHERE permission = 'exam:list' AND delete_flag = 0 ORDER BY id ASC LIMIT 1);

INSERT INTO sys_menu (
  name, permission, type, sort, parent_id, path, icon, component, component_name,
  status, visible, keep_alive, always_show, creator, delete_flag
)
SELECT
  '试卷物理删除',
  'exam:paper:purge',
  2,
  100,
  @parent_id,
  '', '', '', '',
  1, 1, 0, 0,
  'migration',
  0
FROM DUAL
WHERE @parent_id IS NOT NULL
  AND (SELECT COUNT(*) FROM sys_menu WHERE permission = 'exam:paper:purge' AND delete_flag = 0) = 0;

SET @menu_id = (SELECT id FROM sys_menu WHERE permission = 'exam:paper:purge' AND delete_flag = 0 ORDER BY id DESC LIMIT 1);

INSERT INTO sys_role_menu (role_id, menu_id, creator, delete_flag)
SELECT r.id, @menu_id, 'migration', 0
FROM sys_role r
WHERE r.code = 'super_admin'
  AND r.delete_flag = 0
  AND r.status = 0
  AND @menu_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM sys_role_menu rm
    WHERE rm.role_id = r.id AND rm.menu_id = @menu_id AND rm.delete_flag = 0
  );
