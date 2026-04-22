-- =============================================================================
-- 删除「引用 exam_paper.id，但该 id 在 exam_paper 表中不存在」的关联数据（孤儿数据）
-- 条件：LEFT JOIN exam_paper p ON p.id = ... WHERE p.id IS NULL
-- 说明：不区分 exam_paper.delete_flag；只要主键行不存在即删。若需同时清理「卷已逻辑删除」
--       的子树，请使用单独策略（勿与本脚本混用）。
--
-- 务必：先全库备份；在从库或维护窗口执行；建议先运行 exam_paper_integrity_audit.sql 核对。
-- =============================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

START TRANSACTION;

-- ---------------------------------------------------------------------------
-- 1. 答题侧（依赖 exam_attempt.exam_paper_id → exam_paper.id）
-- ---------------------------------------------------------------------------

DELETE a
FROM exam_attempt_answer a
         INNER JOIN exam_attempt e ON e.id = a.attempt_id
         LEFT JOIN exam_paper p ON p.id = e.exam_paper_id
WHERE p.id IS NULL;

DELETE r
FROM exam_result r
         LEFT JOIN exam_paper p ON p.id = r.exam_paper_id
WHERE p.id IS NULL;

-- 会话本身 exam_paper_id 指向不存在的卷
DELETE e
FROM exam_attempt e
         LEFT JOIN exam_paper p ON p.id = e.exam_paper_id
WHERE p.id IS NULL;

-- 仍可能残留：answer 引用不存在的 attempt（无 paper 时已上面删 attempt，此处兜底）
DELETE a
FROM exam_attempt_answer a
         LEFT JOIN exam_attempt e ON e.id = a.attempt_id
WHERE e.id IS NULL;

DELETE r
FROM exam_result r
         LEFT JOIN exam_attempt e ON e.id = r.attempt_id
WHERE e.id IS NULL;

-- ---------------------------------------------------------------------------
-- 2. 批次（exam_batch_paper / exam_batch_member 含 exam_paper_id）
-- ---------------------------------------------------------------------------

DELETE bp
FROM exam_batch_paper bp
         LEFT JOIN exam_paper p ON p.id = bp.exam_paper_id
WHERE p.id IS NULL;

DELETE bm
FROM exam_batch_member bm
         LEFT JOIN exam_paper p ON p.id = bm.exam_paper_id
WHERE p.id IS NULL;

-- ---------------------------------------------------------------------------
-- 3. 题目树（自下而上）：option → question → block → section
-- ---------------------------------------------------------------------------

-- 3.1 选项：试题已不存在
DELETE o
FROM exam_option o
         LEFT JOIN exam_question q ON q.id = o.question_id
WHERE q.id IS NULL;

-- 3.2 选项：试题存在但块/节/卷链断裂
DELETE o
FROM exam_option o
         INNER JOIN exam_question q ON q.id = o.question_id
         LEFT JOIN exam_question_block b ON b.id = q.block_id
WHERE b.id IS NULL;

DELETE o
FROM exam_option o
         INNER JOIN exam_question q ON q.id = o.question_id
         INNER JOIN exam_question_block b ON b.id = q.block_id
         LEFT JOIN exam_section s ON s.id = b.section_id
WHERE s.id IS NULL;

DELETE o
FROM exam_option o
         INNER JOIN exam_question q ON q.id = o.question_id
         INNER JOIN exam_question_block b ON b.id = q.block_id
         INNER JOIN exam_section s ON s.id = b.section_id
         LEFT JOIN exam_paper p ON p.id = s.exam_paper_id
WHERE p.id IS NULL;

-- 3.3 小题：块不存在；或节/卷不存在；或 exam_paper_id 冗余字段指向不存在卷
DELETE q
FROM exam_question q
         LEFT JOIN exam_question_block b ON b.id = q.block_id
WHERE b.id IS NULL;

DELETE q
FROM exam_question q
         INNER JOIN exam_question_block b ON b.id = q.block_id
         LEFT JOIN exam_section s ON s.id = b.section_id
WHERE s.id IS NULL;

DELETE q
FROM exam_question q
         INNER JOIN exam_question_block b ON b.id = q.block_id
         INNER JOIN exam_section s ON s.id = b.section_id
         LEFT JOIN exam_paper p ON p.id = s.exam_paper_id
WHERE p.id IS NULL;

DELETE q
FROM exam_question q
         LEFT JOIN exam_paper p ON p.id = q.exam_paper_id
WHERE p.id IS NULL;

-- 3.4 题块：大题不存在；或大题指向不存在卷
DELETE b
FROM exam_question_block b
         LEFT JOIN exam_section s ON s.id = b.section_id
WHERE s.id IS NULL;

DELETE b
FROM exam_question_block b
         INNER JOIN exam_section s ON s.id = b.section_id
         LEFT JOIN exam_paper p ON p.id = s.exam_paper_id
WHERE p.id IS NULL;

-- 3.5 大题：exam_paper_id 无对应主表行
DELETE s
FROM exam_section s
         LEFT JOIN exam_paper p ON p.id = s.exam_paper_id
WHERE p.id IS NULL;

COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

-- =============================================================================
-- 执行后建议：核对 exam_paper_integrity_audit.sql 第一节 SELECT 是否均为空行
-- =============================================================================
