# 试卷导入（exam_paper）设计说明

## 入口与数据来源

- 管理端仅选择 Mock 卷（`mock_examination_paper.id`），不再由用户填写 `index_url` / 粘贴 `index.json`。
- 服务端读取该 Mock 行的 `resource_url`（通常为资源包 `.zip` 的 URL），将路径后缀 **`.zip` 替换为 `/index.json`**，得到完整 `index.json` 地址并 HTTP 拉取。
- 由该 URL 解析出 `level`、`paper_id`（远程目录 ID）、以及资源基址 `source_base_url`，规则与历史 `parseIndexURL` 一致（路径需满足 `…/{level}/{paper_id}/index.json` 且至少四段路径）。

## 试卷名称（`exam_paper.title`）

- 若导入请求中填写了「试卷名称」，则写入 `title`。
- 若未填写，则使用 Mock 卷的 `name`；若仍为空则回退为 `index.json` 内的 `title` 字段。

## 冲突策略（`conflict_mode`）

对「该 Mock 下是否已存在未逻辑删除的 `exam_paper`」分支处理：

| 值 | 行为 |
|----|------|
| `fail` | 已存在则**不写入**，接口返回 `conflict=true`（与历史行为一致）。 |
| `overwrite` | 对已存在的 `exam_paper`：**仅覆盖指定的 `overwrite_exam_paper_id`**（且该卷必须属于当前 Mock）；不改 `exam_paper.id`，先逻辑删除该卷下题目树，再在同一主键上重建。 |
| `new` | 无论该 Mock 是否已存在导入记录，均**新建一条** `exam_paper`（新 `exam_paper.id`）并导入题目树。 |
| `new_copy`（兼容） | 视为 `new`。 |

## 关于主键与多版本

- 迁移 `022_exam_paper_multi_version_by_mock.sql` 已将 `mock_examination_paper_id` 从唯一约束调整为普通索引，允许同一 Mock 卷下存在多份 `exam_paper`。
- `exam_attempt`、`exam_result` 等表冗余 `exam_paper_id`；物理删卷或改主键需批量迁移关联行，且与「保留历史答题与旧题快照」的目标不一致。
- **「不删除原有数据」**在此指：旧版题目树与选项仍以逻辑删除形式留在库中，便于审计与追溯；`overwrite` 维持同一主键，`new` 会新增主键。

## 缓存

- 导入成功后需失效按 `exam_paper.id` 维度的考前缓存，并调用 `exampaper.InvalidateByMockIDCache(mock_id)`，避免 `ByMockID` 进程内缓存指向过期的卷行快照。

## 前端展示

- 管理端可只读展示由 `resource_url` 推导出的 `index.json` URL 供核对，具体交互与文案以产品为准；**实现细节与上表策略以本文档为准**。

## 管理端试卷列表（`/admin/exam/paper/list`）

- 列表以 **`exam_paper` 为主表**：`id` 为 `exam_paper.id`，并返回 `mock_examination_paper_id` 供导入按 Mock 业务 id 的入口沿用。
- 支持按 `exam_paper.level` 字符串筛选与分页。

## 管理端试卷详情（`GET /admin/exam/paper/{id}`）

- 路径参数 `id` **以 `exam_paper.id` 为准**。控制器不再经 `exampaper.ByMockID` 一次转换；响应体中仍附带 `paper.id`（= `mock_examination_paper_id`）、`paper.exam_paper_id` 两个口径。
- 找不到记录时返回 `code=11114`（`exam_paper_not_found`）。

## 管理端编辑元数据（`POST /admin/exam/paper/edit`）

- 入参 **`exam_paper_id`**（与列表项 `id` 一致），可更新：`title`、`prepare_*`、`source_base_url`（非空时规范为以 `/` 结尾）、`duration_seconds`（`0` 表示走系统默认时长）。
- 不包含题目树与 `index_json` 的修改。

## 管理端听力 HLS（`POST /admin/exam/paper/update`）

- 入参 **`id` 同样为 `exam_paper.id`**，与详情、编辑接口对齐，避免在 Mock 主键与 Exam 主键之间来回转换。

## 考试批次（`exam_batch_paper` / `exam_batch_member`）

- 创建批次时管理端只传 **`exam_paper_ids`**；服务端写入 **`exam_batch_paper`** 时同时写入 **`exam_paper_id`** 与从 `exam_paper` 解析出的 **`mock_examination_paper_id`**（二者必须与 `exam_paper` 行一致）。
- 导入批次成员时按 **`exam_paper_id`** 绑定；写入 **`exam_batch_member`** 时同样冗余 **`mock_examination_paper_id`**，便于按 Mock 卷口径查询或兼容旧逻辑。
- 迁移脚本见 **`migrations/010_exam_batch_paper.sql`**（若库内曾执行过不含 mock 列的旧版 010，需手工 `ALTER` 补齐列并用 `exam_paper` 回填）。

## Mock 详情接口下线

- 管理端 `/admin/mock/examination-paper/{id}` 不再暴露（API 类型与 controller 方法已删除）。
- 列表接口 `/admin/mock/examination-paper/list` 保留，仅用于「导入试卷」弹窗按等级筛选候选 Mock 卷。
- 答题详情等仍需要 Mock 卷名称/元信息的内部逻辑，直接走 `mocksvc.Mock()` 服务层或 DAO，不再依赖对外的管理端 HTTP 接口。
