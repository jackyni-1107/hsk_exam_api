-- 020: sys_member 改为可逆 SM2 密码后，扩容存储字段避免密文被截断
-- 说明：
-- 1) bcrypt 常见长度约 60，而 SM2 密文（hex）通常 > 128。
-- 2) 旧字段长度不足会导致写入时被截断，进而登录恒失败。

ALTER TABLE `sys_member`
    MODIFY COLUMN `password` varchar(512) NOT NULL DEFAULT '' COMMENT '密码（sys_member: 可逆密文）';

ALTER TABLE `sys_password_history`
    MODIFY COLUMN `password_hash` varchar(512) NOT NULL DEFAULT '' COMMENT '历史密码（admin:bcrypt / member:可逆密文）';
