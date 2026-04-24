<template>
  <div class="page">
    <el-card shadow="never">
      <template #header><span>考试批次</span></template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="试卷">
          <el-select
            v-model="query.exam_paper_id"
            clearable
            filterable
            placeholder="全部（exam_paper）"
            style="width: 260px"
          >
            <el-option
              v-for="p in examPaperOptions"
              :key="p.id"
              :label="`${p.id} · ${p.title}`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="考试时间">
          <el-date-picker
            v-model="query.timeRange"
            type="datetimerange"
            value-format="YYYY-MM-DD HH:mm:ss"
            range-separator="至"
            start-placeholder="起"
            end-placeholder="止"
            clearable
            style="width: 400px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
          <el-button type="success" @click="openCreate">新建批次</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="批次ID" width="88" />
        <el-table-column
          label="试卷（exam_paper）"
          min-width="220"
          show-overflow-tooltip
        >
          <template #default="{ row }">{{
            examPaperIdsCell(row.exam_paper_ids)
          }}</template>
        </el-table-column>
        <el-table-column
          prop="title"
          label="批次名称"
          min-width="140"
          show-overflow-tooltip
        />
        <el-table-column label="类型" width="72" align="center">
          <template #default="{ row }">{{ batchKindLabel(row.batch_kind) }}</template>
        </el-table-column>
        <el-table-column label="多次" width="64" align="center">
          <template #default="{ row }">{{ yesNo(row.allow_multiple_attempts) }}</template>
        </el-table-column>
        <el-table-column label="上限" width="72" align="center" show-overflow-tooltip>
          <template #default="{ row }">{{ maxAttemptsCell(row) }}</template>
        </el-table-column>
        <el-table-column label="免成绩" width="72" align="center">
          <template #default="{ row }">{{ yesNo(row.skip_scoring) }}</template>
        </el-table-column>
        <el-table-column label="截止交卷" width="88" align="center">
          <template #default="{ row }">{{ yesNo(row.auto_submit_on_deadline) }}</template>
        </el-table-column>
        <el-table-column
          prop="exam_start_at"
          label="开始时间"
          width="220"
          show-overflow-tooltip
          :formatter="formatUtcForDisplay"
        />
        <el-table-column
          prop="exam_end_at"
          label="结束时间"
          width="220"
          show-overflow-tooltip
          :formatter="formatUtcForDisplay"
        />
        <el-table-column
          prop="member_count"
          label="学员数"
          width="80"
          align="right"
        />
        <el-table-column
          prop="create_time"
          label="创建时间"
          width="220"
          show-overflow-tooltip
          :formatter="formatUtcForDisplay"
        />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              :disabled="isExamBatchEnded(row)"
              :title="isExamBatchEnded(row) ? '考试已结束，不可编辑' : ''"
              @click="openEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="(row.batch_kind ?? 0) !== 1"
              link
              type="primary"
              @click="openMembers(row)"
              >成员</el-button
            >
            <el-button link type="danger" @click="onDelete(row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="pager">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.size"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          background
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="formVisible"
      :title="formMode === 'create' ? '新建批次' : '编辑批次'"
      width="680px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="140px"
      >
        <el-form-item label="试卷" prop="exam_paper_ids">
          <el-select
            v-model="form.exam_paper_ids"
            multiple
            filterable
            placeholder="选择 exam_paper（可多选）"
            style="width: 100%"
          >
            <el-option
              v-for="p in examPaperOptions"
              :key="p.id"
              :label="`${p.id} · ${p.title}`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="批次名称" prop="title">
          <el-input v-model="form.title" clearable placeholder="可选" />
        </el-form-item>
        <el-form-item label="考试开始" prop="exam_start_at">
          <el-date-picker
            v-model="form.exam_start_at"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="选择时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="考试结束" prop="exam_end_at">
          <el-date-picker
            v-model="form.exam_end_at"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="选择时间"
            style="width: 100%"
          />
        </el-form-item>

        <el-divider content-position="left">批次策略</el-divider>
        <el-form-item label="批次类型">
          <el-radio-group v-model="form.batch_kind">
            <el-radio :label="0">正式考试</el-radio>
            <el-radio :label="1">练习 / 模拟</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="允许多次作答">
          <el-switch
            v-model="form.allow_multiple_attempts"
            :active-value="1"
            :inactive-value="0"
            @change="onAllowMultipleChange"
          />
          <span class="field-hint">同一学员同一套卷可多次新建答题会话</span>
        </el-form-item>
        <el-form-item label="每人每卷上限" prop="max_attempts_per_member">
          <el-input-number
            v-model="form.max_attempts_per_member"
            :min="0"
            :precision="0"
            :disabled="form.allow_multiple_attempts !== 1"
            controls-position="right"
            style="width: 200px"
          />
          <span class="field-hint">0 表示不限制（仅「允许多次」开启时有效）</span>
        </el-form-item>
        <el-form-item label="跳过正式成绩">
          <el-switch v-model="form.skip_scoring" :active-value="1" :inactive-value="0" />
          <span class="field-hint">交卷后不写入正式成绩、不参与算分落库</span>
        </el-form-item>
        <el-form-item label="截止自动交卷">
          <el-switch
            v-model="form.auto_submit_on_deadline"
            :active-value="1"
            :inactive-value="0"
          />
          <span class="field-hint">关闭后：个人倒计时到期不自动交卷，超时仍可保存；批次过期也不代交卷</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" :loading="formSaving" @click="submitForm"
          >保存</el-button
        >
      </template>
    </el-dialog>

    <el-drawer
      v-model="memberDrawer"
      title="批次成员"
      size="min(640px, 92vw)"
      destroy-on-close
      @opened="loadMemberList"
    >
      <template v-if="currentBatch">
        <p class="hint">
          批次 #{{ currentBatch.id }} · {{ currentBatch.title || "（无标题）" }}
        </p>
        <el-alert
          v-if="isExamBatchEnded(currentBatch)"
          type="warning"
          :closable="false"
          class="mb"
          show-icon
        >
          该批次考试已结束，无法导入新学员
        </el-alert>
        <el-form class="member-import-form mb" label-width="120px">
          <el-form-item label="导入试卷" required>
            <el-select
              v-model="importExamPaperId"
              placeholder="选择 exam_paper"
              class="import-mock-paper-select"
              filterable
              clearable
              :disabled="isExamBatchEnded(currentBatch)"
            >
              <el-option
                v-for="id in currentBatch.exam_paper_ids || []"
                :key="id"
                :label="examPaperLabel(id)"
                :value="id"
              />
            </el-select>
          </el-form-item>
        </el-form>
        <el-space wrap class="mb">
          <el-input
            v-model="importIdsText"
            type="textarea"
            :rows="2"
            placeholder="输入要导入的会员 ID，逗号或换行分隔"
            style="width: 320px"
            :disabled="isExamBatchEnded(currentBatch)"
          />
          <el-button
            type="primary"
            :loading="importing"
            :disabled="isExamBatchEnded(currentBatch)"
            @click="doImport"
          >
            导入
          </el-button>
        </el-space>
        <el-form :inline="true" class="filter" @submit.prevent="loadMemberList">
          <el-form-item label="账号">
            <el-input
              v-model="memberSearch.username"
              clearable
              placeholder="模糊"
              style="width: 140px"
              :disabled="isExamBatchEnded(currentBatch)"
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :disabled="isExamBatchEnded(currentBatch)"
              @click="searchMembersToPick"
            >
              搜索会员
            </el-button>
          </el-form-item>
        </el-form>
        <el-table
          v-if="pickMembers.length"
          :data="pickMembers"
          size="small"
          max-height="200"
          border
          class="mb"
        >
          <el-table-column prop="id" label="ID" width="72" />
          <el-table-column prop="username" label="账号" />
          <el-table-column prop="nickname" label="昵称" />
          <el-table-column label="" width="100">
            <template #default="{ row }">
              <el-button
                link
                type="primary"
                :disabled="isExamBatchEnded(currentBatch)"
                @click="addPickId(row.id)"
              >
                加入导入
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-table
          v-loading="memberLoading"
          :data="memberRows"
          border
          size="small"
        >
          <el-table-column prop="member_id" label="会员ID" width="88" />
          <el-table-column
            label="试卷 ID"
            min-width="200"
            show-overflow-tooltip
          >
            <template #default="{ row }">{{
              examPaperLabel(row.exam_paper_id)
            }}</template>
          </el-table-column>
          <el-table-column prop="username" label="账号" />
          <el-table-column prop="nickname" label="昵称" />
          <el-table-column
            prop="import_time"
            label="导入时间"
            width="172"
            :formatter="formatUtcForDisplay"
          />
          <el-table-column label="" width="88">
            <template #default="{ row }">
              <el-button link type="danger" @click="removeOne(row)"
                >移除</el-button
              >
            </template>
          </el-table-column>
        </el-table>
        <div class="pager">
          <el-pagination
            v-model:current-page="memberQuery.page"
            v-model:page-size="memberQuery.size"
            :total="memberTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            small
            background
            @size-change="loadMemberList"
            @current-change="loadMemberList"
          />
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import type { FormInstance, FormRules } from "element-plus";
import {
  createExamBatch,
  deleteExamBatch,
  getExamBatchList,
  getExamBatchMemberList,
  getExamPaperList,
  importExamBatchMembers,
  removeExamBatchMembers,
  updateExamBatch,
  type ExamBatchListItem,
  type ExamBatchMemberItem,
  type ExamPaperItem,
} from "@/api/exam";
import { fetchMemberList, type MemberItem } from "@/api/member";
import {
  formatUtcForDisplay,
  getDisplayTimeZone,
  isoOrRfcToWallClockForPicker,
  wallClockStringToRFC3339UTC,
} from "@/utils/datetime";

const batchTimeZone = computed(() => getDisplayTimeZone());

const loading = ref(false);
const rows = ref<ExamBatchListItem[]>([]);
const total = ref(0);
/** 下拉用：拉取足够多的 exam_paper（管理端分页接口） */
const examPaperOptions = ref<ExamPaperItem[]>([]);

const examPaperById = computed(() => {
  const m = new Map<number, ExamPaperItem>();
  for (const p of examPaperOptions.value) {
    m.set(p.id, p);
  }
  return m;
});

/** 单行展示：id · 试卷标题 */
function examPaperLabel(id: number): string {
  const p = examPaperById.value.get(id);
  return p ? `${p.id} · ${p.title}` : String(id);
}

/** 批次列表单元格：多张卷分号分隔 */
function examPaperIdsCell(ids: number[] | undefined): string {
  if (!ids?.length) return "—";
  return ids.map((id) => examPaperLabel(id)).join("；");
}

const query = reactive({
  exam_paper_id: undefined as number | undefined,
  /** [time_from, time_to]，与列表筛选区间求交集 */
  timeRange: null as [string, string] | null,
  page: 1,
  size: 10,
});

async function loadExamPapersForSelect() {
  const res = await getExamPaperList({ page: 1, size: 500 });
  examPaperOptions.value = res.data?.list ?? [];
}

function wallToApiOrUndefined(wall: string | undefined): string | undefined {
  if (!wall?.trim()) return undefined;
  return wallClockStringToRFC3339UTC(wall.trim());
}

async function loadList() {
  loading.value = true;
  try {
    const tr = query.timeRange;
    let time_from: string | undefined;
    let time_to: string | undefined;
    if (tr?.[0] || tr?.[1]) {
      try {
        if (tr[0]) time_from = wallToApiOrUndefined(tr[0]);
        if (tr[1]) time_to = wallToApiOrUndefined(tr[1]);
      } catch (e) {
        ElMessage.error(e instanceof Error ? e.message : "筛选时间转换失败");
        return;
      }
    }
    const res = await getExamBatchList({
      exam_paper_id: query.exam_paper_id || undefined,
      time_from,
      time_to,
      page: query.page,
      size: query.size,
    });
    rows.value = res.data?.list ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function resetQuery() {
  query.exam_paper_id = undefined;
  query.timeRange = null;
  query.page = 1;
  loadList();
}

const formVisible = ref(false);
const formMode = ref<"create" | "edit">("create");
const formSaving = ref(false);
const formRef = ref<FormInstance>();
const editingId = ref(0);
const form = reactive({
  exam_paper_ids: [] as number[],
  title: "",
  exam_start_at: "",
  exam_end_at: "",
  batch_kind: 0,
  allow_multiple_attempts: 0 as 0 | 1,
  max_attempts_per_member: 0,
  skip_scoring: 0 as 0 | 1,
  auto_submit_on_deadline: 1 as 0 | 1,
});

const formRules: FormRules = {
  exam_paper_ids: [
    { type: "array", required: true, message: "请选择试卷", trigger: "change" },
    { type: "array", min: 1, message: "至少选择一张卷", trigger: "change" },
  ],
  exam_start_at: [
    { required: true, message: "请选择开始时间", trigger: "change" },
  ],
  exam_end_at: [
    { required: true, message: "请选择结束时间", trigger: "change" },
  ],
  max_attempts_per_member: [
    {
      validator: (_rule, val: number, cb) => {
        if (val < 0 || !Number.isInteger(val)) {
          cb(new Error("须为不小于 0 的整数"));
        } else {
          cb();
        }
      },
      trigger: "blur",
    },
  ],
};

function batchKindLabel(k: number | undefined): string {
  return k === 1 ? "练习" : "正式";
}

function yesNo(v: number | undefined): string {
  return v === 1 ? "是" : "否";
}

function maxAttemptsCell(row: ExamBatchListItem): string {
  if (row.allow_multiple_attempts !== 1) return "—";
  const m = row.max_attempts_per_member ?? 0;
  return m === 0 ? "不限" : String(m);
}

function onAllowMultipleChange() {
  if (form.allow_multiple_attempts !== 1) {
    form.max_attempts_per_member = 0;
  }
}

function resetFormPolicyDefaults() {
  form.batch_kind = 0;
  form.allow_multiple_attempts = 0;
  form.max_attempts_per_member = 0;
  form.skip_scoring = 0;
  form.auto_submit_on_deadline = 1;
}

function fillFormPolicyFromRow(row: ExamBatchListItem) {
  form.batch_kind = row.batch_kind ?? 0;
  form.allow_multiple_attempts = (row.allow_multiple_attempts === 1 ? 1 : 0) as 0 | 1;
  form.max_attempts_per_member = row.max_attempts_per_member ?? 0;
  form.skip_scoring = (row.skip_scoring === 1 ? 1 : 0) as 0 | 1;
  form.auto_submit_on_deadline = (row.auto_submit_on_deadline === 0 ? 0 : 1) as 0 | 1;
}

function openCreate() {
  formMode.value = "create";
  editingId.value = 0;
  form.exam_paper_ids = [];
  form.title = "";
  form.exam_start_at = "";
  form.exam_end_at = "";
  resetFormPolicyDefaults();
  formVisible.value = true;
}

function isExamBatchEnded(
  row: { exam_end_at: string } | null | undefined,
): boolean {
  if (!row?.exam_end_at) return false;
  const t = new Date(row.exam_end_at).getTime();
  if (Number.isNaN(t)) return false;
  return t <= Date.now();
}

function openEdit(row: ExamBatchListItem) {
  if (isExamBatchEnded(row)) {
    ElMessage.warning("该考试批次已结束，不能编辑");
    return;
  }
  formMode.value = "edit";
  editingId.value = row.id;
  form.exam_paper_ids = [...(row.exam_paper_ids || [])];
  form.title = row.title;
  form.exam_start_at = isoOrRfcToWallClockForPicker(row.exam_start_at);
  form.exam_end_at = isoOrRfcToWallClockForPicker(row.exam_end_at);
  fillFormPolicyFromRow(row);
  formVisible.value = true;
}

async function submitForm() {
  await formRef.value?.validate().catch(() => Promise.reject());
  if (!form.exam_paper_ids.length) return;
  formSaving.value = true;
  try {
    const start = form.exam_start_at.trim();
    const end = form.exam_end_at.trim();
    let exam_start_at: string;
    let exam_end_at: string;
    try {
      exam_start_at = wallClockStringToRFC3339UTC(start);
      exam_end_at = wallClockStringToRFC3339UTC(end);
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : "时间转换失败");
      return;
    }
    const policy = {
      batch_kind: form.batch_kind,
      allow_multiple_attempts: form.allow_multiple_attempts,
      max_attempts_per_member:
        form.allow_multiple_attempts === 1
          ? Math.floor(Number(form.max_attempts_per_member) || 0)
          : 0,
      skip_scoring: form.skip_scoring,
      auto_submit_on_deadline: form.auto_submit_on_deadline,
    };
    if (formMode.value === "create") {
      await createExamBatch({
        title: form.title,
        exam_start_at,
        exam_end_at,
        exam_paper_ids: form.exam_paper_ids,
        ...policy,
      });
      ElMessage.success("已创建");
    } else {
      await updateExamBatch(editingId.value, {
        title: form.title,
        exam_start_at,
        exam_end_at,
        exam_paper_ids: form.exam_paper_ids,
        ...policy,
      });
      ElMessage.success("已保存");
    }
    formVisible.value = false;
    loadList();
  } finally {
    formSaving.value = false;
  }
}

async function onDelete(row: ExamBatchListItem) {
  await ElMessageBox.confirm(`确定删除批次 #${row.id}？`, "确认", {
    type: "warning",
  });
  await deleteExamBatch(row.id);
  ElMessage.success("已删除");
  loadList();
}

const memberDrawer = ref(false);
const currentBatch = ref<ExamBatchListItem | null>(null);
const memberLoading = ref(false);
const memberRows = ref<ExamBatchMemberItem[]>([]);
const memberTotal = ref(0);
const memberQuery = reactive({ page: 1, size: 10 });
const importIdsText = ref("");
const importExamPaperId = ref(0);
const importing = ref(false);
const pickMembers = ref<MemberItem[]>([]);

function openMembers(row: ExamBatchListItem) {
  currentBatch.value = row;
  memberQuery.page = 1;
  importIdsText.value = "";
  importExamPaperId.value = row.exam_paper_ids?.[0] ?? 0;
  pickMembers.value = [];
  memberDrawer.value = true;
}

async function loadMemberList() {
  if (!currentBatch.value) return;
  memberLoading.value = true;
  try {
    const res = await getExamBatchMemberList(currentBatch.value.id, {
      page: memberQuery.page,
      size: memberQuery.size,
    });
    memberRows.value = res.data?.list ?? [];
    memberTotal.value = res.data?.total ?? 0;
  } finally {
    memberLoading.value = false;
  }
}

function parseIdList(text: string): number[] {
  const parts = text
    .split(/[\s,;，；]+/)
    .map((s) => s.trim())
    .filter(Boolean);
  const ids: number[] = [];
  for (const p of parts) {
    const n = Number(p);
    if (Number.isFinite(n) && n > 0) ids.push(n);
  }
  return [...new Set(ids)];
}

async function doImport() {
  if (!currentBatch.value) return;
  if (isExamBatchEnded(currentBatch.value)) {
    ElMessage.warning("该考试批次已结束，不能导入学员");
    return;
  }
  if (!importExamPaperId.value) {
    ElMessage.warning("请选择试卷");
    return;
  }
  const ids = parseIdList(importIdsText.value);
  if (!ids.length) {
    ElMessage.warning("请输入有效的会员 ID");
    return;
  }
  importing.value = true;
  try {
    const res = await importExamBatchMembers(currentBatch.value.id, {
      exam_paper_id: importExamPaperId.value,
      member_ids: ids,
    });
    ElMessage.success(`已导入 ${res.data?.inserted ?? 0} 人`);
    importIdsText.value = "";
    loadMemberList();
    loadList();
  } finally {
    importing.value = false;
  }
}

async function removeOne(row: ExamBatchMemberItem) {
  if (!currentBatch.value) return;
  await ElMessageBox.confirm("从批次中移除该学员？", "确认", {
    type: "warning",
  });
  await removeExamBatchMembers(currentBatch.value.id, {
    exam_paper_id: row.exam_paper_id,
    member_ids: [row.member_id],
  });
  ElMessage.success("已移除");
  loadMemberList();
  loadList();
}

const memberSearch = reactive({ username: "" });

async function searchMembersToPick() {
  const res = await fetchMemberList({
    page: 1,
    size: 20,
    username: memberSearch.username || undefined,
    status: 0,
  });
  pickMembers.value = res.data?.list ?? [];
}

function addPickId(id: number) {
  const cur = parseIdList(importIdsText.value);
  if (!cur.includes(id)) cur.push(id);
  importIdsText.value = cur.join(", ");
}

onMounted(async () => {
  await loadExamPapersForSelect();
  loadList();
});
</script>

<style scoped>
.page {
  padding: 16px;
}
.filter {
  margin-bottom: 12px;
}
.pager {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}
.hint {
  margin: 0 0 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
.mb {
  margin-bottom: 12px;
}
/* 抽屉内避免 inline 表单把 el-select 压成仅箭头宽度 */
.member-import-form {
  width: 100%;
  max-width: 520px;
}
.import-mock-paper-select {
  width: 100%;
}
.import-mock-paper-select :deep(.el-select__wrapper) {
  min-height: var(--el-component-size);
}
.tz-hint {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
}
.tz-hint-inline {
  margin-top: 6px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.field-hint {
  margin-left: 10px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  vertical-align: middle;
}
</style>
