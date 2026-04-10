# exam logic 代码优化记录

> 本文档记录 `internal/logic/exam/` 包的代码整理与优化内容，便于后续维护与 Code Review。

---

## 1. 统一 DAO 引入方式

| 项目 | 说明 |
|------|------|
| 涉及文件 | `internal/logic/exam/import.go` |
| 问题 | 该文件使用 `examdao "exam/internal/dao/exam"` 直接访问底层 DAO 子包，而包内其余 13 个文件均通过 `"exam/internal/dao"` 统一入口访问。两套入口增加阅读和维护成本。 |
| 方案 | 将所有 `examdao.ExamXxx` 替换为 `dao.ExamXxx`，删除 `examdao` 导入，与全包风格一致。 |

改动前：

```go
import (
    examdao "exam/internal/dao/exam"
)
// ...
examdao.ExamPaper.Ctx(ctx)...
```

改动后：

```go
import (
    "exam/internal/dao"
)
// ...
dao.ExamPaper.Ctx(ctx)...
```

---

## 2. 提取重复 attempt 查询为公共函数

| 项目 | 说明 |
|------|------|
| 涉及文件 | `client_attempt.go`、`iexam_bridge.go`、`random_fill_persist.go`、`hls_audio.go` |
| 问题 | 「按 ID + 用户 + delete_flag 加载 attempt 并判 Id==0」的查询模式在 7 处重复出现（`GetAttempt`、`SaveAnswers`、`SubmitAttempt`、`maybeAutoSubmitIfOverdue`、`assertAttemptInProgress`、两处 `RandomFillAnswersForTest`）。 |
| 方案 | 新增两个包级复用函数，各调用点改为调用它们。 |

新增函数（定义在 `client_attempt.go`）：

```go
func loadAttemptByUser(ctx context.Context, attemptID, userID int64) (examentity.ExamAttempt, error)
func assertAttemptInProgressByUser(ctx context.Context, attemptID, userID int64) (*examentity.ExamAttempt, error)
```

- `loadAttemptByUser`：加载并校验存在性，返回 attempt 实体或 `CodeExamAttemptNotFound`。
- `assertAttemptInProgressByUser`：在上者基础上额外校验 `status == InProgress` 且未超时。

---

## 3. 解决 RandomFillAnswersForTest 命名歧义

| 项目 | 说明 |
|------|------|
| 涉及文件 | `random_fill_persist.go`、`iexam_bridge.go` |
| 问题 | 包级函数 `RandomFillAnswersForTest` 会调用 `SaveAnswers` 实际写库；`sExam.RandomFillAnswersForTest` 仅返回 `[]RandomAnswerDraftItem` 不写库。两者同名但语义截然不同，极易误用。 |
| 方案 | 包级函数重命名为 `RandomFillAndSaveAnswers`，方法保留原名并补充注释。 |

改动前：

```go
// random_fill_persist.go
func RandomFillAnswersForTest(ctx, userID, paperID, attemptID) (filled int, err error) { ... SaveAnswers(...) }

// iexam_bridge.go
func (s *sExam) RandomFillAnswersForTest(ctx, ...) ([]bo.RandomAnswerDraftItem, error) { ... /* 不写库 */ }
```

改动后：

```go
// random_fill_persist.go
func RandomFillAndSaveAnswers(ctx, userID, paperID, attemptID) (filled int, err error) { ... SaveAnswers(...) }

// iexam_bridge.go — 仅返回草稿，不写库
func (s *sExam) RandomFillAnswersForTest(ctx, ...) ([]bo.RandomAnswerDraftItem, error) { ... }
```

---

## 4. Redis KEYS 替换为 SCAN

| 项目 | 说明 |
|------|------|
| 涉及文件 | `internal/logic/exam/for_exam.go` |
| 问题 | `invalidatePaperSectionExamCachesByPaper` 使用 `KEYS pattern` 命令批量删除 section 缓存。`KEYS` 在数据量大时会阻塞 Redis 主线程，属于生产环境已知的高风险操作（Redis 官方文档明确标注为 O(N) 阻塞）。 |
| 方案 | 替换为 `SCAN cursor MATCH pattern COUNT 100` 迭代模式，分批收集 key 后 `DEL`。 |

改动前：

```go
v, err := g.Redis().Do(ctx, "KEYS", pattern)
keys := v.Strings()
g.Redis().Del(ctx, keys...)
```

改动后：

```go
var cursor int64
for {
    v, err := g.Redis().Do(ctx, "SCAN", cursor, "MATCH", pattern, "COUNT", 100)
    // 解析 cursor 和 keys
    nextCursor := gvar.New(arr[0]).Int64()
    keys := gvar.New(arr[1]).Strings()
    if len(keys) > 0 {
        g.Redis().Del(ctx, keys...)
    }
    cursor = nextCursor
    if cursor == 0 { break }
}
```

---

## 5. 删除脆弱的 isNoRowsErr 字符串判断

| 项目 | 说明 |
|------|------|
| 涉及文件 | `internal/logic/exam/import.go` |
| 问题 | `isNoRowsErr` 通过 `strings.Contains(err.Error(), "no rows in result set")` 判断空结果。这种写法在数据库驱动版本升级或错误信息国际化时会失效，属于脆弱代码。 |
| 方案 | GoFrame 的 `Scan` 在无匹配行时不返回错误而是返回零值结构体，调用方已通过 `exist.Id > 0` 做了正确判断。`isNoRowsErr` 属于冗余防御代码，直接删除函数及其 3 处调用分支。 |

改动前：

```go
if err := dao.ExamPaper.Ctx(ctx)...Scan(&exist); err != nil {
    if !isNoRowsErr(err) {
        return nil, err
    }
}
```

改动后：

```go
if err := dao.ExamPaper.Ctx(ctx)...Scan(&exist); err != nil {
    return nil, err
}
```

---

## 6. loadQuestionScoreMetaTx 消除 N+1 查询

| 项目 | 说明 |
|------|------|
| 涉及文件 | `internal/logic/exam/client_attempt.go` |
| 问题 | `loadQuestionScoreMetaTx` 在循环中对每道题单独执行 `SELECT ... WHERE question_id = ? AND is_correct = 1`，形成 N+1 查询模式。此外错误被 `_` 忽略，导致数据库查询失败时静默返回空切片。 |
| 方案 | 改为一次性 `WhereIn("question_id", allQIDs).Where("is_correct", 1)` 批量查询，按 `question_id` 在内存中分组，同时正确处理 `Scan` 错误。 |

改动前：

```go
for _, q := range qs {
    var opts []examentity.ExamOption
    _ = tx.Model(dao.ExamOption.Table())...
        Where("question_id", q.Id).Where("is_correct", 1)...Scan(&opts)
    // ...
}
```

改动后：

```go
var correctOpts []examentity.ExamOption
if err := tx.Model(dao.ExamOption.Table())...
    WhereIn("question_id", qIDs).Where("is_correct", 1)...
    Scan(&correctOpts); err != nil {
    return nil, err
}
correctByQ := make(map[int64][]int64)
for _, o := range correctOpts {
    correctByQ[o.QuestionId] = append(correctByQ[o.QuestionId], o.Id)
}
```

---

## 7. 硬编码 updater 字符串常量化

| 项目 | 说明 |
|------|------|
| 涉及文件 | `s_exam.go`（常量定义）、`client_attempt.go`、`attempt_admin.go` |
| 问题 | `"client"`、`"admin"`、`"task"` 作为 creator/updater 分散硬编码于多处，存在拼写不一致风险且不利于全局搜索。 |
| 方案 | 在 `s_exam.go` 中定义包级常量，全包统一引用。 |

```go
// s_exam.go
const (
    updaterClient = "client"
    updaterAdmin  = "admin"
    updaterTask   = "task"
)
```

替换范围：`client_attempt.go` 中 ~10 处、`attempt_admin.go` 中 ~4 处。

---

## 未改动部分（保持现状）

| 项目 | 原因 |
|------|------|
| `iexam_bridge.go` 的 bridge 委托模式 | GoFrame `service` 接口注册的标准做法，利于解耦与测试替换 |
| `for_exam.go` 与 `paper.go` 的相似树构建逻辑 | 两者场景不同（考前脱敏 vs 管理端完整），字段裁剪差异大，合并反而增加复杂度 |
| `batch_admin.go` 中 `ExamBatchDelete` 的子表保留策略 | 属于产品/审计决策，非代码问题 |
| `hls_audio.go` 整体结构 | 职责明确，代码质量可接受 |
