<template>
  <div class="page">
    <el-card shadow="never">
      <template #header>
        <div class="head">
          <span>试卷管理（exam_paper）</span>
          <el-button type="primary" @click="openImport">导入试卷</el-button>
        </div>
      </template>
      <el-form :inline="true" class="filter" @submit.prevent="loadList">
        <el-form-item label="级别">
          <el-input
            v-model="query.level"
            clearable
            placeholder="如 hsk1、new1"
            style="width: 180px"
            @clear="loadList"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadList">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="loading" :data="rows" border stripe>
        <el-table-column prop="id" label="试卷ID" width="88" fixed="left" />
        <el-table-column
          prop="mock_examination_paper_id"
          label="Mock卷ID"
          width="100"
        />
        <el-table-column prop="level" label="级别" width="88" />
        <el-table-column
          prop="paper_id"
          label="远程目录"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column
          prop="title"
          label="标题"
          min-width="160"
          show-overflow-tooltip
        />
        <el-table-column
          label="资源基址"
          min-width="200"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ row.source_base_url }}
          </template>
        </el-table-column>
        <el-table-column
          prop="audio_hls_prefix"
          label="HLS前缀"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column label="分片数" width="80" align="right">
          <template #default="{ row }">
            {{ row.audio_hls_segment_count }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="168">
          <template #default="{ row }">
            {{ formatUtcText(row.create_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)"
              >详情</el-button
            >
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="primary" @click="openImportForRow(row)"
              >导入</el-button
            >
            <el-button link type="primary" @click="openSettings(row)"
              >听力 HLS</el-button
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

    <!-- 导入 -->
    <el-dialog
      v-model="importDlg"
      title="导入试卷"
      width="640px"
      destroy-on-close
      @open="onImportOpen"
    >
      <el-form label-width="140px">
        <el-form-item label="Mock 卷" required>
          <div class="mock-row">
            <el-select
              v-model="importForm.levelFilter"
              placeholder="等级筛选"
              clearable
              style="width: 160px; margin-right: 8px"
              @change="onImportLevelChange"
            >
              <el-option
                v-for="lv in mockLevels"
                :key="lv.id"
                :label="
                  lv.level_name || lv.app_level_name || String(lv.level_id)
                "
                :value="lv.id"
              />
            </el-select>
            <el-select
              v-model="importForm.mock_examination_paper_id"
              placeholder="选择模拟卷"
              filterable
              style="width: 260px"
            >
              <el-option
                v-for="p in mockPapers"
                :key="p.id"
                :label="`${p.id} · ${p.name}`"
                :value="p.id"
              />
            </el-select>
          </div>
          <div class="mock-hint">
            或直接填写 ID：
            <el-input-number
              v-model="importForm.mock_examination_paper_id"
              :min="1"
              :controls="false"
            />
          </div>
        </el-form-item>
        <el-form-item label="试卷名称">
          <el-input
            v-model="importForm.title"
            clearable
            placeholder="可选；不填则使用 Mock 卷名称"
          />
        </el-form-item>
        <el-form-item v-if="derivedImportIndexUrl" label="index 地址">
          <el-input
            :model-value="derivedImportIndexUrl"
            type="textarea"
            :rows="2"
            readonly
          />
        </el-form-item>
        <el-form-item label="听力 HLS 前缀">
          <el-input
            v-model="importForm.audio_hls_prefix"
            placeholder="可选，无首尾 /"
          />
        </el-form-item>
        <el-form-item label="冲突策略">
          <el-radio-group v-model="importForm.conflict_mode">
            <el-radio label="fail">已存在则拒绝</el-radio>
            <el-radio label="overwrite">覆盖</el-radio>
            <el-radio label="new">新试卷</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importDlg = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="submitImport"
          >开始导入</el-button
        >
      </template>
    </el-dialog>

    <!-- 编辑试卷元数据 -->
    <el-dialog
      v-model="editDlg"
      title="编辑试卷"
      width="600px"
      destroy-on-close
    >
      <el-form
        v-loading="editLoading"
        label-width="120px"
        class="edit-paper-form"
      >
        <el-form-item label="试卷ID">
          <span>{{ editForm.exam_paper_id }}</span>
        </el-form-item>
        <el-form-item label="标题" required>
          <el-input v-model="editForm.title" placeholder="试卷标题" />
        </el-form-item>
        <el-form-item label="考前标题">
          <el-input v-model="editForm.prepare_title" />
        </el-form-item>
        <el-form-item label="考前说明">
          <el-input v-model="editForm.prepare_instruction" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="考前音频">
          <el-input v-model="editForm.prepare_audio_file" placeholder="文件名或 URL" />
        </el-form-item>
        <el-form-item label="资源基址" required>
          <el-input
            v-model="editForm.source_base_url"
            type="textarea"
            :rows="2"
            placeholder="以 / 结尾"
          />
        </el-form-item>
        <el-form-item label="考试时长(秒)">
          <el-input-number
            v-model="editForm.duration_seconds"
            :min="0"
            :controls="false"
            style="width: 100%"
          />
          <div class="hint">0 表示使用系统默认时长</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDlg = false">取消</el-button>
        <el-button type="primary" :loading="editSaving" @click="submitEdit"
          >保存</el-button
        >
      </template>
    </el-dialog>

    <!-- 设置（仅已导入 exam 时） -->
    <el-dialog
      v-model="settingsDlg"
      title="听力 HLS 设置"
      width="560px"
      destroy-on-close
    >
      <el-form label-width="168px">
        <el-form-item label="试卷ID">
          <span>{{ settingsPaperId }}</span>
        </el-form-item>
        <div class="hint form-top-hint">
          此处仅配置听力 HLS 相关字段；答题时长在「编辑」中维护。
        </div>
        <el-form-item label="HLS 前缀">
          <el-input
            v-model="settingsForm.audio_hls_prefix"
            placeholder="无首尾 /"
          />
        </el-form-item>
        <el-form-item label="分片总数">
          <el-input-number
            v-model="settingsForm.audio_hls_segment_count"
            :min="0"
            :controls="false"
            style="width: 100%"
          />
          <div class="hint">0 表示未配置整卷听力 HLS</div>
        </el-form-item>
        <el-form-item label="分片文件名 fmt">
          <el-input
            v-model="settingsForm.audio_hls_segment_pattern"
            placeholder="空则默认 %05d.ts"
          />
        </el-form-item>
        <el-form-item label="密钥对象路径">
          <el-input
            v-model="settingsForm.audio_hls_key_object"
            placeholder="相对 prefix，空表示不加密"
          />
        </el-form-item>
        <el-form-item label="IV 十六进制">
          <el-input
            v-model="settingsForm.audio_hls_iv_hex"
            placeholder="AES-128 IV"
          />
        </el-form-item>
        <el-form-item label="#EXTINF 时长(秒)">
          <el-input-number
            v-model="settingsForm.audio_hls_segment_duration"
            :min="0"
            :precision="3"
            :step="0.1"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="settingsDlg = false">取消</el-button>
        <el-button
          type="primary"
          :loading="savingSettings"
          @click="submitSettings"
          >保存</el-button
        >
      </template>
    </el-dialog>

    <!-- 详情：以 exam_paper 为主 -->
    <el-drawer
      v-model="detailDrawer"
      title="试卷详情"
      size="70%"
      destroy-on-close
    >
      <div v-loading="detailLoading" class="detail-wrap">
        <template v-if="examDetail">
          <h4 class="sec-title">考试试卷（exam_paper）</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="试卷ID">{{
              examDetail.paper.exam_paper_id
            }}</el-descriptions-item>
            <el-descriptions-item label="Mock卷ID">{{
              examDetail.paper.id
            }}</el-descriptions-item>
            <el-descriptions-item label="级别">{{
              examDetail.paper.level
            }}</el-descriptions-item>
            <el-descriptions-item label="远程 paper_id">{{
              examDetail.paper.paper_id
            }}</el-descriptions-item>
            <el-descriptions-item label="标题">{{
              examDetail.paper.title
            }}</el-descriptions-item>
            <el-descriptions-item label="考试时长(秒)">{{
              examDetail.paper.duration_seconds
            }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{
              formatUtcText(examDetail.paper.create_time)
            }}</el-descriptions-item>
            <el-descriptions-item label="资源基址" :span="2">{{
              examDetail.paper.source_base_url
            }}</el-descriptions-item>
            <el-descriptions-item label="HLS 前缀" :span="2">{{
              examDetail.paper.audio_hls_prefix
            }}</el-descriptions-item>
            <el-descriptions-item label="HLS 分片数">{{
              examDetail.paper.audio_hls_segment_count
            }}</el-descriptions-item>
            <el-descriptions-item label="分片 fmt">{{
              examDetail.paper.audio_hls_segment_pattern
            }}</el-descriptions-item>
            <el-descriptions-item label="密钥对象" :span="2">{{
              examDetail.paper.audio_hls_key_object
            }}</el-descriptions-item>
            <el-descriptions-item label="IV hex" :span="2">{{
              examDetail.paper.audio_hls_iv_hex
            }}</el-descriptions-item>
            <el-descriptions-item label="#EXTINF(秒)">{{
              examDetail.paper.audio_hls_segment_duration
            }}</el-descriptions-item>
            <el-descriptions-item label="考前标题" :span="2">{{
              examDetail.paper.prepare_title
            }}</el-descriptions-item>
            <el-descriptions-item label="考前说明" :span="2">{{
              examDetail.paper.prepare_instruction
            }}</el-descriptions-item>
            <el-descriptions-item label="考前音频" :span="2">{{
              examDetail.paper.prepare_audio_file
            }}</el-descriptions-item>
          </el-descriptions>
          <h4 class="sec-title">index.json</h4>
          <pre class="json-preview">{{
            truncate(examDetail.paper.index_json, 8000)
          }}</pre>
          <h4 class="sec-title">大题（{{ examDetail.sections.length }}）</h4>
          <el-collapse>
            <el-collapse-item
              v-for="s in examDetail.sections"
              :key="s.id"
              :title="`${s.sort_order}. ${s.topic_title || s.topic_type} (#${s.id})`"
            >
              <el-descriptions :column="2" size="small" border>
                <el-descriptions-item label="topic_type">{{
                  s.topic_type
                }}</el-descriptions-item>
                <el-descriptions-item label="文件">{{
                  s.topic_items_file
                }}</el-descriptions-item>
                <el-descriptions-item label="part_code">{{
                  s.part_code
                }}</el-descriptions-item>
                <el-descriptions-item label="segment_code">{{
                  s.segment_code
                }}</el-descriptions-item>
              </el-descriptions>
              <div class="sub">topic_json 预览</div>
              <pre class="json-preview sm">{{
                truncate(s.topic_json, 4000)
              }}</pre>
              <div class="sub">题块数 {{ s.blocks?.length ?? 0 }}</div>
            </el-collapse-item>
          </el-collapse>
        </template>

        <el-alert
          v-if="detailLoaded && !examDetail"
          type="info"
          :closable="false"
          show-icon
          class="mt"
          title="未加载到考试内容。"
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, computed } from "vue";
import { ElMessage } from "element-plus";
import {
  getExamPaperList,
  getExamPaperDetail,
  importExamPaper,
  updateExamPaper,
  editExamPaperMeta,
  type ExamPaperDetail,
  type ExamPaperItem,
} from "@/api/exam";
import {
  getMockLevelsList,
  getMockExaminationPapers,
  type MockLevelItem,
  type MockExaminationPaperItem,
} from "@/api/mockAdmin";
import { formatUtcText } from "@/utils/datetime";

const loading = ref(false);
const rows = ref<ExamPaperItem[]>([]);
const total = ref(0);
const query = reactive({
  page: 1,
  size: 10,
  /** exam_paper.level，空表示全部 */
  level: "",
});

const editDlg = ref(false);
const editLoading = ref(false);
const editSaving = ref(false);
const editForm = reactive({
  exam_paper_id: 0,
  title: "",
  prepare_title: "",
  prepare_instruction: "",
  prepare_audio_file: "",
  source_base_url: "",
  duration_seconds: 0,
});

const importDlg = ref(false);
/** 从表格「导入」打开时预填 Mock 卷 id（在 onImportOpen 末尾应用） */
const importPreset = ref<{ mockId: number } | null>(null);
const importing = ref(false);
const mockLevels = ref<MockLevelItem[]>([]);
const mockPapers = ref<MockExaminationPaperItem[]>([]);
const importForm = reactive({
  levelFilter: undefined as number | undefined,
  mock_examination_paper_id: undefined as number | undefined,
  title: "",
  audio_hls_prefix: "",
  conflict_mode: "fail" as "fail" | "overwrite" | "new",
});

function derivedIndexUrlFromResource(resource: string): string {
  const s = (resource || "").trim();
  if (!s) return "";
  try {
    const u = new URL(s);
    let path = u.pathname;
    const lp = path.toLowerCase();
    if (lp.endsWith(".zip")) {
      path = path.slice(0, -4) + "/index.json";
    } else if (!lp.endsWith("/index.json")) {
      path = path.replace(/\/$/, "") + "/index.json";
    }
    u.pathname = path;
    u.search = "";
    u.hash = "";
    return u.toString();
  } catch {
    return "";
  }
}

const derivedImportIndexUrl = computed(() => {
  const mid = importForm.mock_examination_paper_id;
  if (!mid) return "";
  const row = mockPapers.value.find((p) => p.id === mid);
  if (row?.resource_url) return derivedIndexUrlFromResource(row.resource_url);
  return "";
});

/** 听力 HLS 表单默认值（新建或未配置时展示/回填） */
const HLS_FORM_DEFAULTS = {
  audio_hls_prefix: "",
  audio_hls_segment_count: 0,
  audio_hls_segment_pattern: "segment_%03d.ts",
  audio_hls_key_object: "static.key",
  audio_hls_iv_hex: "0x0123456789abcdef0123456789abcdef",
  audio_hls_segment_duration: 10,
} as const;

const settingsDlg = ref(false);
const savingSettings = ref(false);
const settingsPaperId = ref(0);
const settingsForm = reactive({ ...HLS_FORM_DEFAULTS });

const detailDrawer = ref(false);
const detailLoading = ref(false);
const detailLoaded = ref(false);
const examDetail = ref<ExamPaperDetail | null>(null);

function truncate(s: string, max: number) {
  if (!s) return "";
  return s.length <= max ? s : s.slice(0, max) + "\n…（已截断）";
}

async function loadList() {
  loading.value = true;
  try {
    const res = (await getExamPaperList({
      page: query.page,
      size: query.size,
      level: query.level.trim() || undefined,
    })) as { data?: { list?: ExamPaperItem[]; total?: number } };
    rows.value = res?.data?.list ?? [];
    total.value = res?.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function resetQuery() {
  query.page = 1;
  query.size = 10;
  query.level = "";
  loadList();
}

async function onImportOpen() {
  importForm.levelFilter = undefined;
  importForm.mock_examination_paper_id = undefined;
  importForm.title = "";
  importForm.audio_hls_prefix = "";
  importForm.conflict_mode = "fail";
  try {
    const res = (await getMockLevelsList()) as {
      data?: { list?: MockLevelItem[] };
    };
    mockLevels.value = res?.data?.list ?? [];
  } catch {
    mockLevels.value = [];
  }
  mockPapers.value = [];
  if (importPreset.value) {
    const { mockId } = importPreset.value;
    importPreset.value = null;
    importForm.mock_examination_paper_id = mockId;
  }
}

async function onImportLevelChange() {
  importForm.mock_examination_paper_id = undefined;
  if (!importForm.levelFilter) {
    mockPapers.value = [];
    return;
  }
  try {
    const res = (await getMockExaminationPapers({
      level_id: importForm.levelFilter,
    })) as {
      data?: { list?: MockExaminationPaperItem[] };
    };
    mockPapers.value = res?.data?.list ?? [];
  } catch {
    mockPapers.value = [];
  }
}

function openImport() {
  importDlg.value = true;
}

async function openEdit(row: ExamPaperItem) {
  editForm.exam_paper_id = row.id;
  editDlg.value = true;
  editLoading.value = true;
  try {
    const res = (await getExamPaperDetail(row.id)) as {
      data?: ExamPaperDetail;
    };
    const p = res?.data?.paper;
    if (!p) {
      ElMessage.warning("无法加载试卷详情");
      editDlg.value = false;
      return;
    }
    editForm.title = p.title ?? "";
    editForm.prepare_title = p.prepare_title ?? "";
    editForm.prepare_instruction = p.prepare_instruction ?? "";
    editForm.prepare_audio_file = p.prepare_audio_file ?? "";
    editForm.source_base_url = p.source_base_url ?? "";
    editForm.duration_seconds = p.duration_seconds ?? 0;
  } catch {
    ElMessage.warning("无法加载试卷详情");
    editDlg.value = false;
  } finally {
    editLoading.value = false;
  }
}

async function submitEdit() {
  if (!editForm.exam_paper_id || !editForm.title.trim()) {
    ElMessage.warning("请填写试卷标题");
    return;
  }
  if (!editForm.source_base_url.trim()) {
    ElMessage.warning("请填写资源基址");
    return;
  }
  editSaving.value = true;
  try {
    await editExamPaperMeta({
      exam_paper_id: editForm.exam_paper_id,
      title: editForm.title.trim(),
      prepare_title: editForm.prepare_title.trim(),
      prepare_instruction: editForm.prepare_instruction.trim(),
      prepare_audio_file: editForm.prepare_audio_file.trim(),
      source_base_url: editForm.source_base_url.trim(),
      duration_seconds: editForm.duration_seconds,
    });
    ElMessage.success("已保存");
    editDlg.value = false;
    await loadList();
  } catch {
    /* */
  } finally {
    editSaving.value = false;
  }
}

function openImportForRow(row: ExamPaperItem) {
  importPreset.value = { mockId: row.mock_examination_paper_id };
  importDlg.value = true;
}

async function submitImport() {
  const mid = importForm.mock_examination_paper_id;
  if (!mid || mid < 1) {
    ElMessage.warning("请选择或填写 Mock 卷 ID");
    return;
  }
  importing.value = true;
  try {
    const payload: Parameters<typeof importExamPaper>[0] = {
      mock_examination_paper_id: mid,
      conflict_mode: importForm.conflict_mode,
      audio_hls_prefix: importForm.audio_hls_prefix || undefined,
      title: importForm.title.trim() || undefined,
    };
    const res = (await importExamPaper(payload)) as {
      data?: {
        conflict?: boolean;
        existing_examination_paper_id?: number;
        section_count?: number;
        question_count?: number;
      };
    };
    const d = res?.data;
    if (d?.conflict) {
      ElMessage.warning(
        `该 Mock 卷已存在导入记录（mock id=${d.existing_examination_paper_id ?? mid}），未写入。可改用覆盖或新试卷。`,
      );
      return;
    }
    ElMessage.success(
      `导入成功：大题 ${d?.section_count ?? 0}，小题 ${d?.question_count ?? 0}`,
    );
    importDlg.value = false;
    await loadList();
  } catch {
    /* axios 已 toast */
  } finally {
    importing.value = false;
  }
}

async function openSettings(row: ExamPaperItem) {
  try {
    const res = (await getExamPaperDetail(row.id)) as {
      data?: ExamPaperDetail;
    };
    const d = res?.data;
    if (!d?.paper) {
      ElMessage.warning("无法加载该试卷详情");
      return;
    }
    settingsPaperId.value = row.id;
    const p = d.paper;
    settingsForm.audio_hls_prefix =
      p.audio_hls_prefix?.trim() || HLS_FORM_DEFAULTS.audio_hls_prefix;
    settingsForm.audio_hls_segment_count =
      p.audio_hls_segment_count ?? HLS_FORM_DEFAULTS.audio_hls_segment_count;
    settingsForm.audio_hls_segment_pattern =
      p.audio_hls_segment_pattern?.trim() ||
      HLS_FORM_DEFAULTS.audio_hls_segment_pattern;
    settingsForm.audio_hls_key_object =
      p.audio_hls_key_object?.trim() || HLS_FORM_DEFAULTS.audio_hls_key_object;
    settingsForm.audio_hls_iv_hex =
      p.audio_hls_iv_hex?.trim() || HLS_FORM_DEFAULTS.audio_hls_iv_hex;
    {
      const dur = p.audio_hls_segment_duration;
      settingsForm.audio_hls_segment_duration =
        dur != null && dur > 0
          ? dur
          : HLS_FORM_DEFAULTS.audio_hls_segment_duration;
    }
    settingsDlg.value = true;
  } catch {
    ElMessage.warning("无法加载该试卷详情");
  }
}

async function submitSettings() {
  savingSettings.value = true;
  try {
    await updateExamPaper({
      id: settingsPaperId.value,
      audio_hls_prefix: settingsForm.audio_hls_prefix,
      audio_hls_segment_count: settingsForm.audio_hls_segment_count,
      audio_hls_segment_pattern: settingsForm.audio_hls_segment_pattern,
      audio_hls_key_object: settingsForm.audio_hls_key_object,
      audio_hls_iv_hex: settingsForm.audio_hls_iv_hex,
      audio_hls_segment_duration: settingsForm.audio_hls_segment_duration,
    });
    ElMessage.success("已保存");
    settingsDlg.value = false;
  } catch {
    /* */
  } finally {
    savingSettings.value = false;
  }
}

async function openDetail(row: ExamPaperItem) {
  examDetail.value = null;
  detailLoaded.value = false;
  detailDrawer.value = true;
  detailLoading.value = true;
  try {
    const examRes = (await getExamPaperDetail(row.id)) as {
      data?: ExamPaperDetail;
    };
    examDetail.value = examRes?.data ?? null;
  } catch {
    examDetail.value = null;
  } finally {
    detailLoading.value = false;
    detailLoaded.value = true;
  }
}

onMounted(async () => {
  await loadList();
});
</script>

<style scoped>
.page {
  padding: 8px 0;
}
.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.filter {
  margin-bottom: 12px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
.mock-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}
.mock-hint {
  margin-top: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.form-top-hint {
  margin: 0 0 12px;
  padding-left: 168px;
  line-height: 1.4;
}
.detail-wrap {
  padding-right: 8px;
}
.sec-title {
  margin: 16px 0 8px;
  font-size: 14px;
}
.json-preview {
  background: var(--el-fill-color-light);
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
  overflow: auto;
  max-height: 320px;
  white-space: pre-wrap;
  word-break: break-all;
}
.json-preview.sm {
  max-height: 200px;
}
.sub {
  margin: 8px 0 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.mt {
  margin-top: 16px;
}
</style>
