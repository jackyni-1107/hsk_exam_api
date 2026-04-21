-- =============================================================================
-- exam_paper 域数据完整性：审计与清洗指引（MySQL 8+）
-- 使用前务必备份；清洗类语句默认注释，确认后再执行。
-- delete_flag：0 未删除，1 已删除（与 internal/consts 一致）
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 一、只读审计：孤立数据（父不存在或指向不存在的 exam_paper）
-- ---------------------------------------------------------------------------

-- 1.1 大题指向不存在的试卷（应无行）
SELECT s.id AS section_id, s.exam_paper_id
FROM exam_section s
LEFT JOIN exam_paper p ON p.id = s.exam_paper_id AND p.delete_flag = 0
WHERE s.delete_flag = 0 AND p.id IS NULL;

-- 1.2 题块指向不存在的大题
SELECT b.id AS block_id, b.section_id
FROM exam_question_block b
LEFT JOIN exam_section s ON s.id = b.section_id AND s.delete_flag = 0
WHERE b.delete_flag = 0 AND s.id IS NULL;

-- 1.3 小题指向不存在的题块
SELECT q.id AS question_id, q.block_id, q.exam_paper_id
FROM exam_question q
LEFT JOIN exam_question_block b ON b.id = q.block_id AND b.delete_flag = 0
WHERE q.delete_flag = 0 AND b.id IS NULL;

-- 1.4 选项指向不存在的试题
SELECT o.id AS option_id, o.question_id
FROM exam_option o
LEFT JOIN exam_question q ON q.id = o.question_id AND q.delete_flag = 0
WHERE o.delete_flag = 0 AND q.id IS NULL;

-- 1.5 答题会话 / 结果指向不存在的试卷（业务上应无行）
SELECT a.id AS attempt_id, a.exam_paper_id
FROM exam_attempt a
LEFT JOIN exam_paper p ON p.id = a.exam_paper_id AND p.delete_flag = 0
WHERE a.delete_flag = 0 AND p.id IS NULL;

SELECT r.attempt_id, r.exam_paper_id
FROM exam_result r
LEFT JOIN exam_paper p ON p.id = r.exam_paper_id AND p.delete_flag = 0
WHERE r.delete_flag = 0 AND p.id IS NULL;

-- 1.6 批次关联指向不存在的试卷
SELECT ebp.batch_id, ebp.exam_paper_id
FROM exam_batch_paper ebp
LEFT JOIN exam_paper p ON p.id = ebp.exam_paper_id AND p.delete_flag = 0
WHERE p.id IS NULL;

SELECT ebm.batch_id, ebm.member_id, ebm.exam_paper_id
FROM exam_batch_member ebm
LEFT JOIN exam_paper p ON p.id = ebm.exam_paper_id AND p.delete_flag = 0
WHERE p.id IS NULL;

-- ---------------------------------------------------------------------------
-- 二、冗余字段一致性：mock_examination_paper_id 应与 exam_paper 一致
-- （仅检查 delete_flag=0 的活跃数据）
-- ---------------------------------------------------------------------------

-- 2.1 section 与 paper 的 mock 不一致
SELECT s.id, s.exam_paper_id, s.mock_examination_paper_id AS sec_mock, p.mock_examination_paper_id AS paper_mock
FROM exam_section s
JOIN exam_paper p ON p.id = s.exam_paper_id AND p.delete_flag = 0
WHERE s.delete_flag = 0 AND s.mock_examination_paper_id <> p.mock_examination_paper_id;

-- 2.2 question 与 paper 的 mock 不一致
SELECT q.id, q.exam_paper_id, q.mock_examination_paper_id AS q_mock, p.mock_examination_paper_id AS paper_mock
FROM exam_question q
JOIN exam_paper p ON p.id = q.exam_paper_id AND p.delete_flag = 0
WHERE q.delete_flag = 0 AND q.mock_examination_paper_id <> p.mock_examination_paper_id;

-- 2.3 question.exam_paper_id 与 block→section 推导不一致（块所属卷 ≠ 题上冗余卷）
SELECT q.id AS question_id, q.exam_paper_id AS q_paper, p2.id AS derived_paper
FROM exam_question q
JOIN exam_question_block b ON b.id = q.block_id AND b.delete_flag = 0
JOIN exam_section s ON s.id = b.section_id AND s.delete_flag = 0
JOIN exam_paper p2 ON p2.id = s.exam_paper_id AND p2.delete_flag = 0
WHERE q.delete_flag = 0 AND q.exam_paper_id <> p2.id;

-- ---------------------------------------------------------------------------
-- 三、软删除层级不一致（子未删、父已删）—— 仅报表，是否“清洗”看产品策略
-- ---------------------------------------------------------------------------

SELECT 'section_child_paper_deleted' AS kind, s.id
FROM exam_section s
JOIN exam_paper p ON p.id = s.exam_paper_id
WHERE s.delete_flag = 0 AND p.delete_flag = 1;

SELECT 'block_child_section_deleted' AS kind, b.id
FROM exam_question_block b
JOIN exam_section s ON s.id = b.section_id
WHERE b.delete_flag = 0 AND s.delete_flag = 1;

SELECT 'question_child_block_deleted' AS kind, q.id
FROM exam_question q
JOIN exam_question_block b ON b.id = q.block_id
WHERE q.delete_flag = 0 AND b.delete_flag = 1;

SELECT 'option_child_question_deleted' AS kind, o.id
FROM exam_option o
JOIN exam_question q ON q.id = o.question_id
WHERE o.delete_flag = 0 AND q.delete_flag = 1;

-- ---------------------------------------------------------------------------
-- 四、清洗思路（务必在事务中、先备份）
-- ---------------------------------------------------------------------------
-- A) 冗余字段纠偏（低风险）：将 section / question 的 mock_examination_paper_id 回填为所属 exam_paper 的值
/*
START TRANSACTION;
UPDATE exam_section s
JOIN exam_paper p ON p.id = s.exam_paper_id AND p.delete_flag = 0
SET s.mock_examination_paper_id = p.mock_examination_paper_id, s.updater = 'integrity_fix'
WHERE s.delete_flag = 0 AND s.mock_examination_paper_id <> p.mock_examination_paper_id;

UPDATE exam_question q
JOIN exam_paper p ON p.id = q.exam_paper_id AND p.delete_flag = 0
SET q.mock_examination_paper_id = p.mock_examination_paper_id, q.updater = 'integrity_fix'
WHERE q.delete_flag = 0 AND q.mock_examination_paper_id <> p.mock_examination_paper_id;

UPDATE exam_question q
JOIN exam_question_block b ON b.id = q.block_id AND b.delete_flag = 0
JOIN exam_section s ON s.id = b.section_id AND s.delete_flag = 0
JOIN exam_paper p ON p.id = s.exam_paper_id AND p.delete_flag = 0
SET q.exam_paper_id = p.id, q.mock_examination_paper_id = p.mock_examination_paper_id, q.updater = 'integrity_fix'
WHERE q.delete_flag = 0 AND q.exam_paper_id <> p.id;
COMMIT;
*/

-- B) 孤立子树：若无业务引用（无 attempt/batch），可按应用内 deletePaperTree 顺序物理删；
--    若有会话/批次引用，必须先迁数据或拒绝删卷。生产环境建议用管理端「物理删除」或专用脚本单卷处理。
--    顺序：exam_option → exam_question → exam_question_block → exam_section →（最后）exam_paper

-- C) 批次表脏行：指向已删或不存在的 exam_paper 时，应先业务确认再删 exam_batch_paper / 修正 exam_batch_member

-- ---------------------------------------------------------------------------
-- 五、应用层建议
-- ---------------------------------------------------------------------------
-- - 日常以「导入/覆盖」事务为准，避免手工改 exam_paper_id。
-- - 物理删除试卷前，业务已约束：无批次引用、无未删答题会话（见 internal/logic/paper/paper_purge.go）。
-- - Redis 缓存：大改后执行 InvalidatePaperForExamCache 或重启/清相关 key（见 internal/logic/paper/redis_paper.go）。
