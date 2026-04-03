-- 运维/排查：管理端 GET /api/admin/exam/paper/{id} 的 path id 为 mock_examination_paper.id。
-- 若接口无数据，先确认 exam 域是否已导入该 mock 卷（存在未删除的 exam_paper 行）。
-- 用法：将下方 195 换成实际 mock 卷 id 后执行。

SELECT
  ep.id AS exam_paper_pk,
  ep.mock_examination_paper_id,
  ep.level,
  ep.paper_id,
  ep.title,
  ep.delete_flag,
  ep.create_time
FROM exam_paper ep
WHERE ep.mock_examination_paper_id = 195;

-- mock 卷是否在 Mock 库存在（表名按实际 mock 库/前缀调整，常见为 mock_examination_paper）
-- SELECT id, name, delete_flag FROM mock_examination_paper WHERE id = 195;
